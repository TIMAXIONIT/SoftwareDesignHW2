package persistence

import (
	"errors"
	"sort"

	"pos/internal/domain"
	"pos/internal/ports"
)

type InMemoryRepository struct {
	students    map[string]*domain.Student
	criteria    map[string]*domain.Criterion
	evaluations map[string]*domain.Evaluation
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		students:    make(map[string]*domain.Student),
		criteria:    make(map[string]*domain.Criterion),
		evaluations: make(map[string]*domain.Evaluation),
	}
}

func (r *InMemoryRepository) SaveStudent(student *domain.Student) error {
	if student == nil {
		return errors.New("student cannot be nil")
	}

	if _, exists := r.students[student.ID()]; exists {
		return ports.ErrAlreadyExists
	}

	r.students[student.ID()] = student
	return nil
}

func (r *InMemoryRepository) FindStudentByID(
	id string,
) (*domain.Student, error) {
	student, exists := r.students[id]
	if !exists {
		return nil, ports.ErrNotFound
	}

	return student, nil
}

func (r *InMemoryRepository) FindAllStudents() ([]*domain.Student, error) {
	result := make([]*domain.Student, 0, len(r.students))

	for _, student := range r.students {
		result = append(result, student)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result, nil
}

func (r *InMemoryRepository) SaveCriterion(
	criterion *domain.Criterion,
) error {
	if criterion == nil {
		return errors.New("criterion cannot be nil")
	}

	if _, exists := r.criteria[criterion.ID()]; exists {
		return ports.ErrAlreadyExists
	}

	r.criteria[criterion.ID()] = criterion
	return nil
}

func (r *InMemoryRepository) FindCriterionByID(
	id string,
) (*domain.Criterion, error) {
	criterion, exists := r.criteria[id]
	if !exists {
		return nil, ports.ErrNotFound
	}

	return criterion, nil
}

func (r *InMemoryRepository) FindAllCriteria() ([]*domain.Criterion, error) {
	result := make([]*domain.Criterion, 0, len(r.criteria))

	for _, criterion := range r.criteria {
		result = append(result, criterion)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result, nil
}

func (r *InMemoryRepository) SaveEvaluation(
	evaluation *domain.Evaluation,
) error {
	if evaluation == nil {
		return errors.New("evaluation cannot be nil")
	}

	if _, exists := r.evaluations[evaluation.ID()]; exists {
		return ports.ErrAlreadyExists
	}

	r.evaluations[evaluation.ID()] = evaluation
	return nil
}

func (r *InMemoryRepository) FindByID(
	id string,
) (*domain.Evaluation, error) {
	evaluation, exists := r.evaluations[id]
	if !exists {
		return nil, ports.ErrNotFound
	}

	return evaluation, nil
}

func (r *InMemoryRepository) FindByStudentID(
	studentID string,
) ([]*domain.Evaluation, error) {
	result := make([]*domain.Evaluation, 0)

	for _, evaluation := range r.evaluations {
		if evaluation.StudentID() == studentID {
			result = append(result, evaluation)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result, nil
}

func (r *InMemoryRepository) Delete(id string) error {
	if _, exists := r.evaluations[id]; !exists {
		return ports.ErrNotFound
	}

	delete(r.evaluations, id)
	return nil
}

func (r *InMemoryRepository) UpdateEvaluation(evaluation *domain.Evaluation) error {
	if evaluation == nil {
		return errors.New("evaluation cannot be nil")
	}

	if _, exists := r.evaluations[evaluation.ID()]; !exists {
		return ports.ErrNotFound
	}

	r.evaluations[evaluation.ID()] = evaluation
	return nil
}

var (
	_ ports.StudentRepository    = (*InMemoryRepository)(nil)
	_ ports.CriterionRepository  = (*InMemoryRepository)(nil)
	_ ports.EvaluationRepository = (*InMemoryRepository)(nil)
)
