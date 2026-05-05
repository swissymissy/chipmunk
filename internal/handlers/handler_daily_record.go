package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/database"
	"github.com/xuri/excelize/v2"
)


// export excel file for daily record
func (cfg *ApiConfig) HandlerExportDailyRecord(w http.ResponseWriter, r *http.Request) {
	// get date from url
	date := r.PathValue("date")

	// fetching data
	records, err := cfg.DB.GetAttendanceByDate(r.Context(), date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "date not found")
			return
		}
		log.Printf("error fetching records by session date: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to export file. Something went wrong")
		return
	}

	// create map to group course name and start time into a group
	courseGroup := make(map[string][]database.GetAttendanceByDateRow)
	var courseOrder []string

	for _, r := range records {
		// create key for course name and start time
		key := fmt.Sprintf("%s_%s", r.CourseName, r.StartTime)
		if _, ok := courseGroup[key]; !ok {
			// add new key to order ( order of the appearing of courses)
			courseOrder = append(courseOrder, key)
		}
		courseGroup[key] = append(courseGroup[key], r)
	}

	// create new excel file
	f := excelize.NewFile()
	defer f.Close()
	styles := NewExcelStyles(f)
	firstSheet := true

	// create sheet for each course
	for _, course := range courseOrder {
		courseRecords := courseGroup[course]
		sheetName := course
		if len(sheetName) > 31 {
			sheetName = sheetName[:31]
		}

		if firstSheet {
			f.SetSheetName("Sheet1", sheetName)
			firstSheet = false
		} else {
			f.NewSheet(sheetName)
		}

		// title
		f.MergeCell(sheetName, "A1", "G1")
		f.SetCellValue(sheetName, "A1", sheetName)
		f.SetCellStyle(sheetName, "A1", "G1", styles.Title)

		// headers
		f.SetSheetRow(sheetName, "A2", &[]interface{}{
			"Student ID",
			"First Name",
			"Last Name",
			"Status",
			"Check-in Time",
			"Course",
			"Start time",
		})
		f.SetCellStyle(sheetName, "A2", "G2", styles.Header)

		// write rows
		for i, r := range courseRecords {
			rowNum := i + 3
			cell := fmt.Sprintf("A%d", rowNum)
			f.SetSheetRow(sheetName, cell, &[]interface{}{
				r.StudentID,
				r.FirstName,
				r.LastName,
				r.Status,
				r.CheckInAt.String,
				r.CourseName,
				r.StartTime,
			})
		}
	}

	fileName := fmt.Sprintf("%s_report.xlsx", date)

	// download headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// write file to response
	_, err = f.WriteTo(w)
	if err != nil {
		log.Printf("error writing excel file to response: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to download file")
		return
	}
}
