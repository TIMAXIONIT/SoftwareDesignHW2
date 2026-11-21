package domain

import "testing"

func TestNewEvaluation(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")
	criterion, _ := NewCriterion("c1", "Практика", 10, 2)
	score, _ := NewScore(8)

	evaluation, err := NewEvaluation(
		"e1",
		student,
		criterion,
		score,
		"Хорошая работа",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if evaluation.StudentID() != "s1" {
		t.Errorf("expected student ID s1, got %s", evaluation.StudentID())
	}

	if evaluation.CriterionID() != "c1" {
		t.Errorf("expected criterion ID c1, got %s", evaluation.CriterionID())
	}

	if evaluation.Score().Value() != 8 {
		t.Errorf("expected score 8, got %f", evaluation.Score().Value())
	}

	if evaluation.Comment() != "Хорошая работа" {
		t.Errorf("unexpected comment: %s", evaluation.Comment())
	}
}

func TestNewEvaluationWithScoreExceedingMaximum(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")
	criterion, _ := NewCriterion("c1", "Практика", 10, 2)
	score, _ := NewScore(11)

	_, err := NewEvaluation("e1", student, criterion, score, "")

	if err != ErrScoreExceedsMax {
		t.Errorf("expected ErrScoreExceedsMax, got %v", err)
	}
}

func TestEvaluationUpdateScore(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")
	criterion, _ := NewCriterion("c1", "Практика", 10, 2)
	score, _ := NewScore(5)

	evaluation, _ := NewEvaluation("e1", student, criterion, score, "")

	newScore, _ := NewScore(9)

	err := evaluation.UpdateScore(newScore, criterion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if evaluation.Score().Value() != 9 {
		t.Errorf("expected score 9, got %f", evaluation.Score().Value())
	}
}

func TestEvaluationUpdateComment(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")
	criterion, _ := NewCriterion("c1", "Практика", 10, 2)
	score, _ := NewScore(5)

	evaluation, _ := NewEvaluation("e1", student, criterion, score, "Старый комментарий")

	evaluation.UpdateComment("Новый комментарий")

	if evaluation.Comment() != "Новый комментарий" {
		t.Errorf("unexpected comment: %s", evaluation.Comment())
	}
}
