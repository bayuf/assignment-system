// service/lecturer_service.go
package service

import (
	"fmt"
	"session-17/model"
	"session-17/repository"
)

type SubmissionService interface {
	GetAllSubmissions() ([]model.Submission, error)
	GradeSubmission(studentID, assignmentID int, grade float64) error

	FindAllByLectureId(id int) ([]model.Submission, error)
	GradeDetailByAssignment(assignmentId int) ([]model.Submission, error)
}

type submissionService struct {
	Repo repository.Repository
}

func NewSubmissionService(subRepo repository.Repository) SubmissionService {
	return &submissionService{
		Repo: subRepo,
	}
}

func (submissionService *submissionService) GradeDetailByAssignment(assignmentId int) ([]model.Submission, error) {

	submissions, err := submissionService.Repo.SubmissionRepo.GradeDetailByAssignment(assignmentId)
	if err != nil {
		return []model.Submission{}, err
	}

	// check name is not null
	for _, a := range submissions {
		if a.StudentName == "" {
			return []model.Submission{}, nil
		}
	}
	return submissions, nil
}

func (submissionService *submissionService) FindAllByLectureId(id int) ([]model.Submission, error) {
	return submissionService.Repo.SubmissionRepo.FindAllByLectureId(id)
}

func (submissionService *submissionService) GetAllSubmissions() ([]model.Submission, error) {
	return submissionService.Repo.SubmissionRepo.GetAllWithStudentAndAssignment()
}

func (s *submissionService) GradeSubmission(studentID, assignmentID int, grade float64) error {
	// Cek apakah submission-nya ada
	sub, err := s.Repo.SubmissionRepo.FindByStudentAndAssignment(studentID, assignmentID)
	if err != nil {
		return fmt.Errorf("submission not found: %w", err)
	}

	// Update nilai
	sub.Grade = &grade
	return s.Repo.SubmissionRepo.UpdateGrade(sub)
}
