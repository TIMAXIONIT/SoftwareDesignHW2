package ports

import "pos/internal/domain"

type CriterionRepository interface {
	SaveCriterion(criterion *domain.Criterion) error
	FindCriterionByID(id string) (*domain.Criterion, error)
	FindAllCriteria() ([]*domain.Criterion, error)
}
