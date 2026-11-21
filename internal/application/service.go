package application

import (
	"errors"

	"pos/internal/domain"
	"pos/internal/ports"
)

var (
	ErrStudentAlreadyExists   = errors.New("student already exists")
	ErrCriterionAlreadyExists = errors.New("criterion already exists")
)

type Service struct {
	students    ports.StudentRepository
	criteria    ports.CriterionRepository
	evaluations ports.EvaluationRepository
	calculator  *domain.ResultDomainService
	factory     *EvaluationFactory
}

func NewService(
	students ports.StudentRepository,
	criteria ports.CriterionRepository,
	evaluations ports.EvaluationRepository,
	calculator *domain.ResultDomainService,
	factory *EvaluationFactory,
) *Service {
	return &Service{
		students:    students,
		criteria:    criteria,
		evaluations: evaluations,
		calculator:  calculator,
		factory:     factory,
	}
}

func (s *Service) CreateStudent(id, name string) (*domain.Student, error) {
	student, err := domain.NewStudent(id, name)
	if err != nil {
		return nil, err
	}

	if _, err := s.students.FindStudentByID(student.ID()); err == nil {
		return nil, ErrStudentAlreadyExists
	} else if !errors.Is(err, ports.ErrNotFound) {
		return nil, err
	}

	if err := s.students.SaveStudent(student); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *Service) CreateCriterion(
	id, name string,
	maxScore, weight float64,
) (*domain.Criterion, error) {
	criterion, err := domain.NewCriterion(id, name, maxScore, weight)
	if err != nil {
		return nil, err
	}

	if _, err := s.criteria.FindCriterionByID(criterion.ID()); err == nil {
		return nil, ErrCriterionAlreadyExists
	} else if !errors.Is(err, ports.ErrNotFound) {
		return nil, err
	}

	if err := s.criteria.SaveCriterion(criterion); err != nil {
		return nil, err
	}

	return criterion, nil
}

func (s *Service) CreateEvaluation(
	id, studentID, criterionID string,
	scoreValue float64,
	comment string,
) (*domain.Evaluation, error) {
	evaluation, err := s.factory.Create(
		id,
		studentID,
		criterionID,
		scoreValue,
		comment,
	)
	if err != nil {
		return nil, err
	}

	if err := s.evaluations.SaveEvaluation(evaluation); err != nil {
		return nil, err
	}

	return evaluation, nil
}

func (s *Service) GetStudentEvaluations(studentID string) ([]*domain.Evaluation, error) {
	return s.evaluations.FindByStudentID(studentID)
}

func (s *Service) CalculateStudentResult(studentID string) (*domain.Result, error) {
	evaluations, err := s.evaluations.FindByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	allCriteria, err := s.criteria.FindAllCriteria()
	if err != nil {
		return nil, err
	}

	criteriaByID := make(map[string]*domain.Criterion, len(allCriteria))
	for _, criterion := range allCriteria {
		criteriaByID[criterion.ID()] = criterion
	}

	return s.calculator.CalculateResult(studentID, evaluations, criteriaByID)
}

func (s *Service) UpdateEvaluationScore(
	evaluationID string,
	scoreValue float64,
) error {
	evaluation, err := s.evaluations.FindByID(evaluationID)
	if err != nil {
		return err
	}

	criterion, err := s.criteria.FindCriterionByID(evaluation.CriterionID())
	if err != nil {
		return err
	}

	score, err := domain.NewScore(scoreValue)
	if err != nil {
		return err
	}

	if err := evaluation.UpdateScore(score, criterion); err != nil {
		return err
	}

	return s.evaluations.UpdateEvaluation(evaluation)
}

func (s *Service) UpdateEvaluationComment(
	evaluationID string,
	comment string,
) error {
	evaluation, err := s.evaluations.FindByID(evaluationID)
	if err != nil {
		return err
	}

	evaluation.UpdateComment(comment)

	return s.evaluations.UpdateEvaluation(evaluation)
}

func (s *Service) DeleteEvaluation(evaluationID string) error {
	return s.evaluations.Delete(evaluationID)
}

func (s *Service) GetAllStudents() ([]*domain.Student, error) {
	return s.students.FindAllStudents()
}

func (s *Service) GetAllCriteria() ([]*domain.Criterion, error) {
	return s.criteria.FindAllCriteria()
}

func (s *Service) GetStudent(studentID string) (*domain.Student, error) {
	return s.students.FindStudentByID(studentID)
}

func (s *Service) GetAllStudentReports() ([]*StudentReport, error) {
	students, err := s.students.FindAllStudents()
	if err != nil {
		return nil, err
	}

	reports := make([]*StudentReport, 0, len(students))

	for _, student := range students {
		evaluations, err := s.evaluations.FindByStudentID(student.ID())
		if err != nil {
			return nil, err
		}

		result, err := s.CalculateStudentResult(student.ID())
		if err != nil {
			if errors.Is(err, domain.ErrNoEvaluations) {
				continue
			}
			return nil, err
		}

		reports = append(reports, &StudentReport{
			Student:     student,
			Evaluations: evaluations,
			Result:      result,
		})
	}

	return reports, nil
}
