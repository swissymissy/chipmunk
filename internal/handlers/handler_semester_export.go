package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/swissymissy/chipmunk/internal/database"
	"github.com/xuri/excelize/v2"
)

type exportReq struct {
	CourseID string `json:"course_id"`
}

// let professor export semester record into excel file
func (cfg *ApiConfig) HandlerExportSemesterRecords(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req exportReq
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding request: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to export")
		return
	}

	// get course infor
	course, err := cfg.DB.GetCourseByID(r.Context(), req.CourseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "course not found")
			return
		}
		log.Printf("error fetching course information: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to fetch course. Something went wrong")
		return
	}

	// get range filter for dates
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	// flexible type to store either type of returned struct from query
	var rows [][]interface{}

	if from != "" && to != "" {
		records, err := cfg.DB.GetAttendanceSummaryByCourseInDateRange(r.Context(), database.GetAttendanceSummaryByCourseInDateRangeParams{
			CourseID:      req.CourseID,
			SessionDate:   from,
			SessionDate_2: to,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("attempt to fetch non-existed or deleted course: %s\n", err)
				ResponseWithError(w, http.StatusNotFound, "course not found")
				return
			}
			log.Printf("error fetching semester summary for course: %s\n", err)
			ResponseWithError(w, http.StatusInternalServerError, "failed to export file. Something went wrong")
			return
		}
		// convert to flexible type and store it outside of scope
		for _, r := range records {
			rows = append(rows, toExcelRow(r.StudentID, r.FirstName, r.LastName, r.Specialty, r.TotalPresent, r.TotalSessions, r.Average))
		}
	} else {
		records, err := cfg.DB.GetAttendanceSummaryByCourse(r.Context(), req.CourseID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("attempt to fetch non-existed or deleted course: %s\n", err)
				ResponseWithError(w, http.StatusNotFound, "course not found")
				return
			}
			log.Printf("error fetching semester summary for course: %s\n", err)
			ResponseWithError(w, http.StatusInternalServerError, "failed to export file. Something went wrong")
			return
		}
		for _, r := range records {
			rows = append(rows, toExcelRow(r.StudentID, r.FirstName, r.LastName, r.Specialty, r.TotalPresent, r.TotalSessions, r.Average))
		}
	}

	// create new file. NewFile will automatically creates sheet1
	f := excelize.NewFile()
	defer f.Close()

	// set sheet name - CourseName_SectionDate_StartTime
	// sanitize course name and time to clean spaces and special characters
	cleanCourseName := strings.ReplaceAll(course.CourseName, " ", "_")
	cleanStartTime := strings.ReplaceAll(course.StartTime, ":", "")

	sheetName := fmt.Sprintf("%s %s %s", cleanCourseName, course.SectionDate, cleanStartTime)
	if len(sheetName) > 31 {
		sheetName = sheetName[:31]
	}
	f.SetSheetName("Sheet1", sheetName)

	// design: merge cells to make a title
	err = f.MergeCell(sheetName, "A1", "H1")
	if err != nil {
		log.Printf("unable to merge cell from A1 to H1 to make title: %s\n", err)
	}
	f.SetCellValue(sheetName, "A1", sheetName)
	styles := NewExcelStyles(f)
	f.SetCellStyle(sheetName, "A1", "H1", styles.Title)
	f.SetCellStyle(sheetName, "A2", "H2", styles.Header)

	// set up columns headers
	f.SetSheetRow(sheetName, "A2", &[]interface{}{
		"Student ID",
		"First Name",
		"Last Name",
		"Major", // specialty
		"Total Present",
		"Total Sessions",
		"Average (%)",
		"Status",
	})

	// build excel file using rows
	for i, r := range rows {
		rowNum := i + 3
		cell := fmt.Sprintf("A%d", rowNum)
		f.SetSheetRow(sheetName, cell, &r)

		// fill colors based on status
		status := r[7].(string)
		applyStatusStyle(f, sheetName, rowNum, status, styles)
	}

	fileName := fmt.Sprintf("%s_%s_%s.xlsx", cleanCourseName, course.SectionDate, cleanStartTime)
	// set up download headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") //this long string is the MIME type for xlsx file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// write file directly to response
	_, err = f.WriteTo(w)
	if err != nil {
		log.Printf("error writing excel file to response: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to download file")
		return
	}
}

// helper function that write values to excel row without restriction of struct type
// bc GetAttendanceSummaryByCourse and GetAttendanceSummaryByCourseInDateRange returns 2 different struct types
// technically both structs type have exact same fields
// but Go does not allow the structs to be used interchangeably
// so I wrote this function to remove that struct restriction by returning an []interface{}
// which is more flexible for branching
func toExcelRow(studentID, firstName, lastName string, specialty sql.NullString, totalPresent, totalSess int64, avg float64) []interface{} {
	var status string
	switch {
	case avg >= 85:
		status = "Qualified"
	case avg >= 70:
		status = "At Risk"
	default:
		status = "Not Qualified"
	}

	return []interface{}{
		studentID,
		firstName,
		lastName,
		specialty.String,
		totalPresent,
		totalSess,
		avg,
		status,
	}
}

// helper to apply color to status cell
func applyStatusStyle(f *excelize.File, sheet string, row int, status string, styles ExcelStyle) {
	cell := fmt.Sprintf("H%d", row)
	switch status {
	case "Qualified":
		f.SetCellStyle(sheet, cell, cell, styles.Qualified)
	case "At Risk":
		f.SetCellStyle(sheet, cell, cell, styles.AtRisk)
	case "Not Qualified":
		f.SetCellStyle(sheet, cell, cell, styles.NotQualified)
	}
}
