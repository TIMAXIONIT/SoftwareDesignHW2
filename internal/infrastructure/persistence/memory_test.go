package persistence

import (
	"testing"

	"pos/internal/domain"
	"pos/internal/ports"
)

func TestSaveAndFindStudent(t *testing.T) {
	repo := NewInMemoryRepository()

	student, err := domain.NewStudent("student-1", "Ivan Ivanov")
	if err != nil {
		t.Fatalf("unexpected error creating student: %v", err)
	}

	if err := repo.SaveStudent(student); err != nil {
		t.Fatalf("unexpected error saving student: %v", err)
	}

	found, err := repo.FindStudentByID("student-1")
	if err != nil {
		t.Fatalf("unexpected error finding student: %v", err)
	}

	if found.ID() != student.ID() || found.Name() != student.Name() {
		t.Fatalf("found student does not match the saved one")
	}
}

func TestSaveDuplicateStudent(t *testing.T) {
	repo := NewInMemoryRepository()

	first, _ := domain.NewStudent("student-1", "Ivan Ivanov")
	second, _ := domain.NewStudent("student-1", "Petr Petrov")

	if err := repo.SaveStudent(first); err != nil {
		t.Fatalf("unexpected error saving first student: %v", err)
	}

	err := repo.SaveStudent(second)
	if err != ports.ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestFindUnknownStudent(t *testing.T) {
	repo := NewInMemoryRepository()

	_, err := repo.FindStudentByID("unknown")
	if err != ports.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFindAllStudents(t *testing.T) {
	repo := NewInMemoryRepository()

	first, _ := domain.NewStudent("student-1", "Ivan Ivanov")
	second, _ := domain.NewStudent("student-2", "Petr Petrov")

	if err := repo.SaveStudent(first); err != nil {
		t.Fatal(err)
	}
	if err := repo.SaveStudent(second); err != nil {
		t.Fatal(err)
	}

	students, err := repo.FindAllStudents()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(students) != 2 {
		t.Fatalf("expected 2 students, got %d", len(students))
	}
}

func TestSaveAndFindCriterion(t *testing.T) {
	repo := NewInMemoryRepository()

	criterion, err := domain.NewCriterion("criterion-1", "Programming", 10, 2)
	if err != nil {
		t.Fatalf("unexpected error creating criterion: %v", err)
	}

	if err := repo.SaveCriterion(criterion); err != nil {
		t.Fatalf("unexpected error saving criterion: %v", err)
	}

	found, err := repo.FindCriterionByID("criterion-1")
	if err != nil {
		t.Fatalf("unexpected error finding criterion: %v", err)
	}

	if found.ID() != criterion.ID() || found.Name() != criterion.Name() {
		t.Fatalf("found criterion does not match the saved one")
	}
}

func TestSaveDuplicateCriterion(t *testing.T) {
	repo := NewInMemoryRepository()

	first, _ := domain.NewCriterion("criterion-1", "Programming", 10, 2)
	second, _ := domain.NewCriterion("criterion-1", "Mathematics", 5, 1)

	if err := repo.SaveCriterion(first); err != nil {
		t.Fatal(err)
	}

	err := repo.SaveCriterion(second)
	if err != ports.ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestSaveFindAndDeleteEvaluation(t *testing.T) {
	repo := NewInMemoryRepository()

	student, err := domain.NewStudent("student-1", "Ivan Ivanov")
	if err != nil {
		t.Fatal(err)
	}

	criterion, err := domain.NewCriterion("criterion-1", "Programming", 10, 2)
	if err != nil {
		t.Fatal(err)
	}

	score, err := domain.NewScore(8)
	if err != nil {
		t.Fatal(err)
	}

	evaluation, err := domain.NewEvaluation(
		"evaluation-1",
		student,
		criterion,
		score,
		"Good work",
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.SaveEvaluation(evaluation); err != nil {
		t.Fatalf("unexpected error saving evaluation: %v", err)
	}

	found, err := repo.FindByID("evaluation-1")
	if err != nil {
		t.Fatalf("unexpected error finding evaluation: %v", err)
	}

	if found.ID() != evaluation.ID() {
		t.Fatalf("found evaluation has an unexpected ID")
	}

	evaluations, err := repo.FindByStudentID("student-1")
	if err != nil {
		t.Fatalf("unexpected error finding evaluations by student: %v", err)
	}

	if len(evaluations) != 1 {
		t.Fatalf("expected 1 evaluation, got %d", len(evaluations))
	}

	if err := repo.Delete("evaluation-1"); err != nil {
		t.Fatalf("unexpected error deleting evaluation: %v", err)
	}

	_, err = repo.FindByID("evaluation-1")
	if err != ports.ErrNotFound {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestUpdateEvaluation(t *testing.T) {
	repo := NewInMemoryRepository()

	student, _ := domain.NewStudent("student-1", "Ivan Ivanov")
	criterion, _ := domain.NewCriterion("criterion-1", "Programming", 10, 2)
	score, _ := domain.NewScore(8)

	evaluation, err := domain.NewEvaluation(
		"evaluation-1",
		student,
		criterion,
		score,
		"Initial comment",
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.SaveEvaluation(evaluation); err != nil {
		t.Fatal(err)
	}

	if err := repo.UpdateEvaluation(evaluation); err != nil {
		t.Fatalf("unexpected error updating evaluation: %v", err)
	}
}
