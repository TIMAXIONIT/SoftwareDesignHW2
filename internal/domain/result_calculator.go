package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoEvaluations          = errors.New("no evaluations available")
	ErrUnknownCriterion       = errors.New("evaluation references unknown criterion")
	ErrDuplicateEvaluation    = errors.New("duplicate evaluation for criterion")
	ErrDifferentStudents      = errors.New("evaluations belong to different students")
	ErrInvalidCalculationData = errors.New("invalid calculation data")
)

type ResultCalculationStrategy interface {
	Calculate(
		evaluations []*Evaluation,
		criteria map[string]*Criterion,
	) (float64, error)
}

type WeightedAverageStrategy struct{}

func (s WeightedAverageStrategy) Calculate(
	evaluations []*Evaluation,
	criteria map[string]*Criterion,
) (float64, error) {
	if len(evaluations) == 0 {
		return 0, ErrNoEvaluations
	}

	var weightedSum float64
	var totalWeight float64

	for _, evaluation := range evaluations {
		criterion := criteria[evaluation.CriterionID()]

		if criterion == nil {
			return 0, ErrUnknownCriterion
		}

		normalizedScore := evaluation.Score().Value() / criterion.MaxScore()

		weightedSum += normalizedScore * criterion.Weight()
		totalWeight += criterion.Weight()
	}

	if totalWeight <= 0 {
		return 0, ErrInvalidCalculationData
	}

	return 10 * weightedSum / totalWeight, nil
}

type ResultDomainService struct {
	strategy ResultCalculationStrategy
}

func NewResultDomainService(
	strategy ResultCalculationStrategy,
) *ResultDomainService {
	return &ResultDomainService{
		strategy: strategy,
	}
}

func (s *ResultDomainService) CalculateResult(
	studentID string,
	evaluations []*Evaluation,
	criteria map[string]*Criterion,
) (*Result, error) {
	if s.strategy == nil {
		return nil, ErrInvalidCalculationData
	}

	if len(evaluations) == 0 {
		return nil, ErrNoEvaluations
	}

	seenCriteria := make(map[string]bool)

	for _, evaluation := range evaluations {
		if evaluation == nil {
			return nil, ErrInvalidCalculationData
		}

		if evaluation.StudentID() != studentID {
			return nil, ErrDifferentStudents
		}

		criterionID := evaluation.CriterionID()

		if seenCriteria[criterionID] {
			return nil, ErrDuplicateEvaluation
		}

		seenCriteria[criterionID] = true

		criterion := criteria[criterionID]
		if criterion == nil {
			return nil, fmt.Errorf(
				"%w: %s",
				ErrUnknownCriterion,
				criterionID,
			)
		}
	}

	score, err := s.strategy.Calculate(evaluations, criteria)
	if err != nil {
		return nil, err
	}

	return NewResult(studentID, score)
}
