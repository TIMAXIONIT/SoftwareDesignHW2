package domain

import "testing"

func TestNewCriterion(t *testing.T) {
	criterion, err := NewCriterion("c1", "Практика", 10, 2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if criterion.ID() != "c1" {
		t.Errorf("expected ID c1, got %s", criterion.ID())
	}

	if criterion.MaxScore() != 10 {
		t.Errorf("expected max score 10, got %f", criterion.MaxScore())
	}

	if criterion.Weight() != 2 {
		t.Errorf("expected weight 2, got %f", criterion.Weight())
	}
}

func TestNewCriterionWithInvalidValues(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		title    string
		maxScore float64
		weight   float64
		wantErr  error
	}{
		{
			name:     "empty ID",
			id:       "",
			title:    "Практика",
			maxScore: 10,
			weight:   1,
			wantErr:  ErrEmptyCriterionID,
		},
		{
			name:     "empty name",
			id:       "c1",
			title:    " ",
			maxScore: 10,
			weight:   1,
			wantErr:  ErrEmptyCriterionName,
		},
		{
			name:     "invalid max score",
			id:       "c1",
			title:    "Практика",
			maxScore: 0,
			weight:   1,
			wantErr:  ErrInvalidMaxScore,
		},
		{
			name:     "invalid weight",
			id:       "c1",
			title:    "Практика",
			maxScore: 10,
			weight:   0,
			wantErr:  ErrInvalidWeight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCriterion(
				tt.id,
				tt.title,
				tt.maxScore,
				tt.weight,
			)

			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNewScore(t *testing.T) {
	score, err := NewScore(8.5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if score.Value() != 8.5 {
		t.Errorf("expected 8.5, got %f", score.Value())
	}
}

func TestNewScoreWithNegativeValue(t *testing.T) {
	_, err := NewScore(-1)

	if err != ErrNegativeScore {
		t.Errorf("expected ErrNegativeScore, got %v", err)
	}
}
