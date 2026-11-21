package domain

import "testing"

func TestWeightedAverageStrategy(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")

	criterion1, _ := NewCriterion("c1", "Практика", 10, 2)
	criterion2, _ := NewCriterion("c2", "Теория", 20, 1)

	score1, _ := NewScore(8)
	score2, _ := NewScore(15)

	evaluation1, _ := NewEvaluation("e1", student, criterion1, score1, "")
	evaluation2, _ := NewEvaluation("e2", student, criterion2, score2, "")

	criteria := map[string]*Criterion{
		"c1": criterion1,
		"c2": criterion2,
	}

	evaluations := []*Evaluation{
		evaluation1,
		evaluation2,
	}

	strategy := WeightedAverageStrategy{}

	result, err := strategy.Calculate(evaluations, criteria)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 10.0 * (0.8*2 + 0.75*1) / 3

	if difference := result - expected; difference > 0.000001 || difference < -0.000001 {
		t.Errorf("expected %.4f, got %.4f", expected, result)
	}
}

func TestResultDomainService(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")
	criterion, _ := NewCriterion("c1", "Практика", 10, 1)
	score, _ := NewScore(8)

	evaluation, _ := NewEvaluation("e1", student, criterion, score, "")

	strategy := WeightedAverageStrategy{}
	service := NewResultDomainService(strategy)

	criteria := map[string]*Criterion{
		"c1": criterion,
	}

	result, err := service.CalculateResult(
		"s1",
		[]*Evaluation{evaluation},
		criteria,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.StudentID() != "s1" {
		t.Errorf("expected student ID s1, got %s", result.StudentID())
	}

	if result.Score() != 8 {
		t.Errorf("expected result 8, got %f", result.Score())
	}
}

func TestResultDomainServiceWithoutEvaluations(t *testing.T) {
	service := NewResultDomainService(WeightedAverageStrategy{})

	_, err := service.CalculateResult(
		"s1",
		[]*Evaluation{},
		map[string]*Criterion{},
	)

	if err != ErrNoEvaluations {
		t.Errorf("expected ErrNoEvaluations, got %v", err)
	}
}

func TestResultDomainServiceWithDuplicateCriteria(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")
	criterion, _ := NewCriterion("c1", "Практика", 10, 1)

	score1, _ := NewScore(8)
	score2, _ := NewScore(9)

	evaluation1, _ := NewEvaluation("e1", student, criterion, score1, "")
	evaluation2, _ := NewEvaluation("e2", student, criterion, score2, "")

	service := NewResultDomainService(WeightedAverageStrategy{})

	_, err := service.CalculateResult(
		"s1",
		[]*Evaluation{evaluation1, evaluation2},
		map[string]*Criterion{"c1": criterion},
	)

	if err != ErrDuplicateEvaluation {
		t.Errorf("expected ErrDuplicateEvaluation, got %v", err)
	}
}
