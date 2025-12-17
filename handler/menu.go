package handler

import (
	"html/template"
	"net/http"
)

type MenuHandler struct {
	Templates *template.Template
}

type HomeViewData struct {
	Role string
}

func NewMenuHandler(templates *template.Template) MenuHandler {
	return MenuHandler{
		Templates: templates,
	}
}

func (h *MenuHandler) HomeView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := HomeViewData{Role: role}

	if err := h.Templates.ExecuteTemplate(w, "home", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *MenuHandler) AssignmentView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := HomeViewData{Role: role}

	if err := h.Templates.ExecuteTemplate(w, "assignment", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *MenuHandler) SubmitView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := HomeViewData{Role: role}

	if err := h.Templates.ExecuteTemplate(w, "submission_form", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *MenuHandler) GradeView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	role := ""

	roleCookie, err := r.Cookie("role")
	if err != nil {
		return
	} else {
		role = roleCookie.Value
	}

	data := HomeViewData{Role: role}

	if err := h.Templates.ExecuteTemplate(w, "grade", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *MenuHandler) PageUnauthorized(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Templates.ExecuteTemplate(w, "page401", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
