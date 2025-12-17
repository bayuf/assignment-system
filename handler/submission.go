package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"session-17/model"
	"session-17/service"
	"strconv"
	"strings"
)

type SubmissionHandler struct {
	SubmissionService service.SubmissionService
	Templates         *template.Template
}

func NewSubmissionHandler(templates *template.Template, submissionService service.SubmissionService) SubmissionHandler {
	return SubmissionHandler{
		SubmissionService: submissionService,
		Templates:         templates,
	}
}

type SubmissionViewData struct {
	Role        string
	Submissions []model.Submission
}

func (submissionHandler SubmissionHandler) ListByLectureId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	c, _ := r.Cookie("session")
	idStr := strings.TrimPrefix(c.Value, "lumos-")
	lecturerId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lecture ID", http.StatusBadRequest)
		return
	}

	submissions, err := submissionHandler.SubmissionService.FindAllByLectureId(lecturerId)
	if err != nil {
		fmt.Println(err)
		return
	}

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := SubmissionViewData{
		Role:        role,
		Submissions: submissions,
	}

	if err := submissionHandler.Templates.ExecuteTemplate(w, "grade", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (submissionHandler SubmissionHandler) GetGradeDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	assignmentId := r.URL.Query().Get("assignment_id")
	assignmentIdInt, err := strconv.Atoi(assignmentId)
	if err != nil {
		return
	}

	submissions, err := submissionHandler.SubmissionService.GradeDetailByAssignment(assignmentIdInt)
	if err != nil {
		fmt.Println(err)
		return
	}

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := SubmissionViewData{
		Role:        role,
		Submissions: submissions,
	}

	if err := submissionHandler.Templates.ExecuteTemplate(w, "grade_detail", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
