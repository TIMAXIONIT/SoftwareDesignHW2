package application

import (
	"pos/internal/domain"
	"pos/internal/ports"
)

type EvaluationFactory struct {
	students ports.StudentRepository
	criteria ports.CriterionRepository
}

func NewEvaluationFactory(
	students ports.StudentRepository,
	criteria ports.CriterionRepository,
) *EvaluationFactory {
	return &EvaluationFactory{
		students: students,
		criteria: criteria,
	}
}

func (f *EvaluationFactory) Create(
	id, studentID, criterionID string,
	scoreValue float64,
	comment string,
) (*domain.Evaluation, error) {
	student, err := f.students.FindStudentByID(studentID)
	if err != nil {
		return nil, err
	}

	criterion, err := f.criteria.FindCriterionByID(criterionID)
	if err != nil {
		return nil, err
	}

	score, err := domain.NewScore(scoreValue)
	if err != nil {
		return nil, err
	}

	return domain.NewEvaluation(
		id,
		student,
		criterion,
		score,
		comment,
	)
}
