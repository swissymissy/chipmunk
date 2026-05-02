package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	// get course report by courseID
	// the returned result will have courseID and section date, time
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

	// create new file. NewFile will automatically creates sheet1
	f := excelize.NewFile()
	defer f.Close()
	
	// set sheet name - CourseName_SectionDate_StartTime
	sheetName := fmt.Sprintf("%s_%s_%s", course.CourseName, course.SectionDate, course.StartTime)
	// design: merge cells to make a title
	err = f.MergeCell(sheetName, "A1", "H1")
	if err != nil {
		log.Printf("unable to merge cell from A1 to H1 to make title: %s\n", err)
	}

	// set up columns headers
	f.SetSheetRow(sheetName, "A2", &[]interface{}{
		"Student ID",
		"First Name",
		"Last Name",
		"Major",	// specialty
		"Total Present",
		"Total Sessions",
		"Average (%)",
		"Status",
	})
	
	// fill up other cells with data , row by row
	for i, r := range records{
		var status string 
		switch {
		case r.Average >= 85:
			status = "Qualified"
		case r.Average >= 70:
			status = "At Risk"
		default:
			status = "Not Qualified"
		}

		cell := fmt.Sprintf("A%d", i+3)
		f.SetSheetRow(sheetName, cell, &[]interface{}{
			r.StudentID,
			r.FirstName,
			r.LastName,
			r.Specialty,
			r.TotalPresent,
			r.TotalSessions,
			r.Average,
			status,
		})
	}

	// sanitize course name and time to clean spaces and special characters
	cleanCourseName := strings.ReplaceAll(course.CourseName, " ", "_")
	cleanStartTime := strings.ReplaceAll(course.StartTime, ":", "")
	fileName := fmt.Sprintf("%s_%s_%s.xlsx", cleanCourseName, course.SectionDate, cleanStartTime)

	// set up download headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") //this long string is the MIME type for xlsx file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attatchment; filename=%s", fileName))

	// write file directly to response
	_, err = f.WriteTo(w)
	if err != nil {
		log.Printf("error writing excel file to response: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to download file")
		return
	}
}	

