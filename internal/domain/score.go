package domain

import (
	"errors"
	"math"
)

var ErrNegativeScore = errors.New("score cannot be negative")

type Score struct {
	value float64
}

func NewScore(value float64) (Score, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return Score{}, ErrNegativeScore
	}

	return Score{value: value}, nil
}

func (s Score) Value() float64 {
	return s.value
}
