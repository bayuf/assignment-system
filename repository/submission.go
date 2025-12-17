package repository

import (
	"context"
	"database/sql"
	"fmt"
	"session-17/database"
	"session-17/model"
)

type SubmissionRepo interface {
	CountByStudentAndAssignment(studentID, assignmentID int) (int64, error)
	Create(submission *model.Submission) error
	GetAllWithStudentAndAssignment() ([]model.Submission, error)
	FindByStudentAndAssignment(studentID, assignmentID int) (*model.Submission, error)
	UpdateGrade(sub *model.Submission) error

	FindAllByLectureId(id int) ([]model.Submission, error)
	GradeDetailByAssignment(id int) ([]model.Submission, error)
}

type submissionRepo struct {
	db database.PgxIface
}

func NewSubmissionRepo(db database.PgxIface) SubmissionRepo {
	return &submissionRepo{db}
}

func (r *submissionRepo) GradeDetailByAssignment(assignmentId int) ([]model.Submission, error) {

	query := `
		SELECT a.id, a.created_at, a.updated_at, 
a.title, u.name, a.deadline,s.submitted_at, s.file_Url, s.grade
		FROM assignments a
		LEFT JOIN submissions s ON a.id=s.assignment_id 
		LEFT JOIN users u ON u.id=s.student_id
		WHERE a.id=$1 AND a.deleted_at IS NULL
		ORDER BY a.deadline ASC
	`
	rows, err := r.db.Query(context.Background(), query, assignmentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissionsDetail []model.Submission
	var studentName sql.NullString
	var submittedAt sql.NullTime
	var fileUrl sql.NullString
	for rows.Next() {
		var a model.Submission
		err := rows.Scan(
			&a.AssignmentID, &a.Model.CreatedAt, &a.Model.UpdatedAt,
			&a.AssignmentTitle, &studentName, &a.Deadline, &submittedAt, &fileUrl, &a.Grade,
		)
		if err != nil {
			return nil, err
		}

		if studentName.Valid {
			a.StudentName = studentName.String
		}
		if submittedAt.Valid {
			a.SubmittedAt = submittedAt.Time
		}
		if fileUrl.Valid {
			a.FileURL = fileUrl.String
		}

		submissionsDetail = append(submissionsDetail, a)
	}

	return submissionsDetail, nil
}

func (r *submissionRepo) FindAllByLectureId(id int) ([]model.Submission, error) {
	query := `
		SELECT a.id, a.created_at, a.updated_at, a.deleted_at, c.name, a.title, a.description, a.deadline
		FROM assignments a
		LEFT JOIN courses c ON a.course_id=c.id
		WHERE a.lecturer_id=$1 AND a.deleted_at IS NULL
		ORDER BY a.deadline ASC
	`
	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []model.Submission
	for rows.Next() {
		var a model.Submission
		err := rows.Scan(&a.AssignmentID, &a.Model.CreatedAt, &a.Model.UpdatedAt,
			&a.Model.DeletedAt, &a.CourseName, &a.AssignmentTitle, &a.AssignmentDec, &a.Deadline)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, a)
	}
	return submissions, nil
}

func (r *submissionRepo) CountByStudentAndAssignment(studentID, assignmentID int) (int64, error) {
	var count int64
	err := r.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM submissions WHERE student_id=$1 AND assignment_id=$2", studentID, assignmentID).Scan(&count)
	return count, err
}

func (r *submissionRepo) Create(sub *model.Submission) error {
	_, err := r.db.Exec(context.Background(), "INSERT INTO submissions (assignment_id, student_id, submitted_at, file_url, status) VALUES ($1, $2, $3, $4, $5)",
		sub.AssignmentID, sub.StudentID, sub.SubmittedAt, sub.FileURL, sub.Status)
	return err
}

func (r *submissionRepo) GetAllWithStudentAndAssignment() ([]model.Submission, error) {
	query := `
		SELECT s.id, s.assignment_id, s.student_id, u.name as student_name,
		       a.title as assignment_title, s.file_url, s.status, s.grade
		FROM submissions s
		JOIN users u ON s.student_id = u.id
		JOIN assignments a ON s.assignment_id = a.id
		ORDER BY s.submitted_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []model.Submission
	for rows.Next() {
		var s model.Submission
		err := rows.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.StudentName, &s.AssignmentTitle, &s.FileURL, &s.Status, &s.Grade)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	fmt.Printf("data %+v", submissions)
	return submissions, nil
}

func (r *submissionRepo) FindByStudentAndAssignment(studentID, assignmentID int) (*model.Submission, error) {
	query := `SELECT id, assignment_id, student_id, submitted_at, file_url, status, grade 
			  FROM submissions 
			  WHERE student_id = $1 AND assignment_id = $2 LIMIT 1`

	row := r.db.QueryRow(context.Background(), query, studentID, assignmentID)

	var sub model.Submission
	err := row.Scan(
		&sub.ID,
		&sub.AssignmentID,
		&sub.StudentID,
		&sub.SubmittedAt,
		&sub.FileURL,
		&sub.Status,
		&sub.Grade,
	)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *submissionRepo) UpdateGrade(sub *model.Submission) error {
	query := `UPDATE submissions SET grade = $1 WHERE student_id = $2 AND assignment_id = $3`
	_, err := r.db.Exec(context.Background(), query, sub.Grade, sub.StudentID, sub.AssignmentID)
	return err
}
