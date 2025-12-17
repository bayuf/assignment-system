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

type AssignmentHandler struct {
	AssignmentService service.AssignmentService
	Templates         *template.Template
}

func NewAssignmentHandler(templates *template.Template, assignmenetService service.AssignmentService) AssignmentHandler {
	return AssignmentHandler{
		AssignmentService: assignmenetService,
		Templates:         templates,
	}
}

type AssignmentViewData struct {
	Role        string
	UserId      string
	Assignments []model.Assignment
	Assignment  *model.Assignment
}

func (AssignmentHandler *AssignmentHandler) List(w http.ResponseWriter, r *http.Request) {
	assignments, err := AssignmentHandler.AssignmentService.GetAllAssignments()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := AssignmentViewData{
		Role:        role,
		Assignments: assignments,
	}

	if err := AssignmentHandler.Templates.ExecuteTemplate(w, "assignment", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (AssignmentHandler *AssignmentHandler) SubmitView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	assignmentIDStr := r.URL.Query().Get("assignment_id")

	assignmentID, _ := strconv.Atoi(assignmentIDStr)

	assignment, err := AssignmentHandler.AssignmentService.GetAssignmentByID(assignmentID)
	if err != nil {
		return
	}

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := AssignmentViewData{
		Role:       role,
		Assignment: assignment,
	}

	if err := AssignmentHandler.Templates.ExecuteTemplate(w, "submission_form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (AssignmentHandler *AssignmentHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "error file size", http.StatusBadRequest)
			return
		}
	}

	// get assignment id
	assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	// get student id
	c, _ := r.Cookie("session")
	idStr := strings.TrimPrefix(c.Value, "lumos-")
	studentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// get file
	file, fileHeander, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "error file", http.StatusBadRequest)
		return
	}

	status, err := AssignmentHandler.AssignmentService.SubmitAssignment(studentID, assignmentID, file, fileHeander)
	if err != nil {
		http.Error(w, "error submit", http.StatusBadRequest)
		return
	}

	fmt.Println(status)
	http.Redirect(w, r, "/user/success-submit", http.StatusSeeOther)
}

func (AssignmentHandler *AssignmentHandler) SuccessSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := AssignmentHandler.Templates.ExecuteTemplate(w, "success_submit", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
