package domain

import (
	"errors"
	"math"
)

var ErrInvalidResultScore = errors.New("result score must be between 0 and 10")

type Result struct {
	studentID string
	score     float64
}

func NewResult(studentID string, score float64) (*Result, error) {
	if studentID == "" {
		return nil, ErrEmptyStudentID
	}

	if math.IsNaN(score) || math.IsInf(score, 0) || score < 0 || score > 10 {
		return nil, ErrInvalidResultScore
	}

	return &Result{
		studentID: studentID,
		score:     score,
	}, nil
}

func (r *Result) StudentID() string {
	return r.studentID
}

func (r *Result) Score() float64 {
	return r.score
}
