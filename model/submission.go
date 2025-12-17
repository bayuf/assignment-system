package model

import (
	"time"
)

type Submission struct {
	Model
	CourseName      string
	AssignmentID    int
	StudentID       int
	SubmittedAt     time.Time
	StudentName     string
	AssignmentTitle string
	AssignmentDec   string
	Deadline        time.Time
	FileURL         string
	Status          string
	Grade           *float64
}
