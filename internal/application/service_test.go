package application

import (
	"errors"
	"testing"

	"pos/internal/domain"
	"pos/internal/infrastructure/persistence"
	"pos/internal/ports"
)

func createTestService() *Service {
	repository := persistence.NewInMemoryRepository()

	strategy := domain.WeightedAverageStrategy{}
	calculator := domain.NewResultDomainService(strategy)

	factory := NewEvaluationFactory(repository, repository)

	return NewService(
		repository,
		repository,
		repository,
		calculator,
		factory,
	)
}

func TestServiceCreateStudent(t *testing.T) {
	service := createTestService()

	student, err := service.CreateStudent("s1", "Иван Иванов")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if student.ID() != "s1" {
		t.Errorf("expected ID s1, got %s", student.ID())
	}

	if student.Name() != "Иван Иванов" {
		t.Errorf("unexpected name: %s", student.Name())
	}
}

func TestServiceCreateDuplicateStudent(t *testing.T) {
	service := createTestService()

	_, err := service.CreateStudent("s1", "Иван Иванов")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.CreateStudent("s1", "Пётр Петров")

	if !errors.Is(err, ErrStudentAlreadyExists) {
		t.Errorf("expected ErrStudentAlreadyExists, got %v", err)
	}
}

func TestServiceCreateCriterion(t *testing.T) {
	service := createTestService()

	criterion, err := service.CreateCriterion(
		"c1",
		"Практика",
		10,
		2,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if criterion.ID() != "c1" {
		t.Errorf("expected ID c1, got %s", criterion.ID())
	}
}

func TestServiceCreateDuplicateCriterion(t *testing.T) {
	service := createTestService()

	_, err := service.CreateCriterion("c1", "Практика", 10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.CreateCriterion("c1", "Теория", 20, 1)

	if !errors.Is(err, ErrCriterionAlreadyExists) {
		t.Errorf("expected ErrCriterionAlreadyExists, got %v", err)
	}
}

func TestServiceCreateEvaluation(t *testing.T) {
	service := createTestService()

	_, err := service.CreateStudent("s1", "Иван Иванов")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.CreateCriterion("c1", "Практика", 10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	evaluation, err := service.CreateEvaluation(
		"e1",
		"s1",
		"c1",
		8,
		"Хорошая работа",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if evaluation.StudentID() != "s1" {
		t.Errorf("expected student ID s1, got %s", evaluation.StudentID())
	}

	if evaluation.Score().Value() != 8 {
		t.Errorf("expected score 8, got %f", evaluation.Score().Value())
	}
}

func TestServiceCreateEvaluationForUnknownStudent(t *testing.T) {
	service := createTestService()

	_, err := service.CreateCriterion("c1", "Практика", 10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.CreateEvaluation("e1", "unknown", "c1", 8, "")

	if !errors.Is(err, ports.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestServiceCreateEvaluationWithInvalidScore(t *testing.T) {
	service := createTestService()

	_, err := service.CreateStudent("s1", "Иван Иванов")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.CreateCriterion("c1", "Практика", 10, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.CreateEvaluation("e1", "s1", "c1", 15, "")

	if !errors.Is(err, domain.ErrScoreExceedsMax) {
		t.Errorf("expected ErrScoreExceedsMax, got %v", err)
	}
}

func TestServiceUpdateEvaluationScore(t *testing.T) {
	service := createTestService()

	_, _ = service.CreateStudent("s1", "Иван Иванов")
	_, _ = service.CreateCriterion("c1", "Практика", 10, 2)
	_, _ = service.CreateEvaluation("e1", "s1", "c1", 5, "")

	err := service.UpdateEvaluationScore("e1", 9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	evaluations, err := service.GetStudentEvaluations("s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if evaluations[0].Score().Value() != 9 {
		t.Errorf("expected score 9, got %f", evaluations[0].Score().Value())
	}
}

func TestServiceUpdateEvaluationComment(t *testing.T) {
	service := createTestService()

	_, _ = service.CreateStudent("s1", "Иван Иванов")
	_, _ = service.CreateCriterion("c1", "Практика", 10, 2)
	_, _ = service.CreateEvaluation("e1", "s1", "c1", 5, "Старый комментарий")

	err := service.UpdateEvaluationComment("e1", "Новый комментарий")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	evaluations, err := service.GetStudentEvaluations("s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if evaluations[0].Comment() != "Новый комментарий" {
		t.Errorf("unexpected comment: %s", evaluations[0].Comment())
	}
}

func TestServiceDeleteEvaluation(t *testing.T) {
	service := createTestService()

	_, _ = service.CreateStudent("s1", "Иван Иванов")
	_, _ = service.CreateCriterion("c1", "Практика", 10, 2)
	_, _ = service.CreateEvaluation("e1", "s1", "c1", 5, "")

	err := service.DeleteEvaluation("e1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.GetStudentEvaluations("s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = service.CalculateStudentResult("s1")

	if !errors.Is(err, domain.ErrNoEvaluations) {
		t.Errorf("expected ErrNoEvaluations, got %v", err)
	}
}

func TestServiceGetStudentReport(t *testing.T) {
	service := createTestService()

	_, _ = service.CreateStudent("s1", "Иван Иванов")
	_, _ = service.CreateCriterion("c1", "Практика", 10, 1)
	_, _ = service.CreateEvaluation("e1", "s1", "c1", 8, "")

	facade := NewFacade(service)

	report, err := facade.GetStudentReport("s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.Student.Name() != "Иван Иванов" {
		t.Errorf("unexpected student name: %s", report.Student.Name())
	}

	if len(report.Evaluations) != 1 {
		t.Errorf("expected 1 evaluation, got %d", len(report.Evaluations))
	}

	if report.Result.Score() != 8 {
		t.Errorf("expected result 8, got %f", report.Result.Score())
	}
}
