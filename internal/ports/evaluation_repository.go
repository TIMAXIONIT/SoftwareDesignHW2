package ports

import "pos/internal/domain"

type EvaluationRepository interface {
	SaveEvaluation(evaluation *domain.Evaluation) error
	FindByID(id string) (*domain.Evaluation, error)
	FindByStudentID(studentID string) ([]*domain.Evaluation, error)
	Delete(id string) error
	UpdateEvaluation(evaluation *domain.Evaluation) error
}
