package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyEvaluationID = errors.New("evaluation ID cannot be empty")
	ErrNilStudent        = errors.New("student cannot be nil")
	ErrNilCriterion      = errors.New("criterion cannot be nil")
	ErrScoreExceedsMax   = errors.New("score exceeds criterion maximum")
)

type Evaluation struct {
	id          string
	studentID   string
	criterionID string
	score       Score
	date        time.Time
	comment     string
}

func NewEvaluation(
	id string,
	student *Student,
	criterion *Criterion,
	score Score,
	comment string,
) (*Evaluation, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return nil, ErrEmptyEvaluationID
	}

	if student == nil {
		return nil, ErrNilStudent
	}

	if criterion == nil {
		return nil, ErrNilCriterion
	}

	if score.Value() > criterion.MaxScore() {
		return nil, ErrScoreExceedsMax
	}

	return &Evaluation{
		id:          id,
		studentID:   student.ID(),
		criterionID: criterion.ID(),
		score:       score,
		date:        time.Now(),
		comment:     comment,
	}, nil
}

func (e *Evaluation) ID() string {
	return e.id
}

func (e *Evaluation) StudentID() string {
	return e.studentID
}

func (e *Evaluation) CriterionID() string {
	return e.criterionID
}

func (e *Evaluation) Score() Score {
	return e.score
}

func (e *Evaluation) Date() time.Time {
	return e.date
}

func (e *Evaluation) Comment() string {
	return e.comment
}

func (e *Evaluation) UpdateScore(
	score Score,
	criterion *Criterion,
) error {
	if criterion == nil {
		return ErrNilCriterion
	}

	if criterion.ID() != e.criterionID {
		return ErrNilCriterion
	}

	if score.Value() > criterion.MaxScore() {
		return ErrScoreExceedsMax
	}

	e.score = score
	e.date = time.Now()

	return nil
}

func (e *Evaluation) UpdateComment(comment string) {
	e.comment = comment
}
