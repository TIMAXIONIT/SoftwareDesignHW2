package domain

import (
	"errors"
	"math"
	"strings"
)

var (
	ErrEmptyCriterionID   = errors.New("criterion ID cannot be empty")
	ErrEmptyCriterionName = errors.New("criterion name cannot be empty")
	ErrInvalidMaxScore    = errors.New("maximum score must be positive")
	ErrInvalidWeight      = errors.New("criterion weight must be positive")
)

type Criterion struct {
	id       string
	name     string
	maxScore float64
	weight   float64
}

func NewCriterion(id, name string, maxScore, weight float64) (*Criterion, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)

	if id == "" {
		return nil, ErrEmptyCriterionID
	}

	if name == "" {
		return nil, ErrEmptyCriterionName
	}

	if math.IsNaN(maxScore) || math.IsInf(maxScore, 0) || maxScore <= 0 {
		return nil, ErrInvalidMaxScore
	}

	if math.IsNaN(weight) || math.IsInf(weight, 0) || weight <= 0 {
		return nil, ErrInvalidWeight
	}

	return &Criterion{
		id:       id,
		name:     name,
		maxScore: maxScore,
		weight:   weight,
	}, nil
}

func (c *Criterion) ID() string {
	return c.id
}

func (c *Criterion) Name() string {
	return c.name
}

func (c *Criterion) MaxScore() float64 {
	return c.maxScore
}

func (c *Criterion) Weight() float64 {
	return c.weight
}
