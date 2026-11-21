package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyStudentID   = errors.New("student ID cannot be empty")
	ErrEmptyStudentName = errors.New("student name cannot be empty")
)

type Student struct {
	id   string
	name string
}

func NewStudent(id, name string) (*Student, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)

	if id == "" {
		return nil, ErrEmptyStudentID
	}

	if name == "" {
		return nil, ErrEmptyStudentName
	}

	return &Student{
		id:   id,
		name: name,
	}, nil
}

func (s *Student) ID() string {
	return s.id
}

func (s *Student) Name() string {
	return s.name
}

func (s *Student) Rename(newName string) error {
	newName = strings.TrimSpace(newName)

	if newName == "" {
		return ErrEmptyStudentName
	}

	s.name = newName
	return nil
}
