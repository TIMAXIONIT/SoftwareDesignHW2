package domain

import "testing"

func TestNewStudent(t *testing.T) {
	student, err := NewStudent("s1", " Иван Иванов ")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if student.ID() != "s1" {
		t.Errorf("expected ID s1, got %s", student.ID())
	}

	if student.Name() != "Иван Иванов" {
		t.Errorf("expected trimmed name, got %s", student.Name())
	}
}

func TestNewStudentWithEmptyID(t *testing.T) {
	_, err := NewStudent("", "Иван Иванов")

	if err != ErrEmptyStudentID {
		t.Errorf("expected ErrEmptyStudentID, got %v", err)
	}
}

func TestNewStudentWithEmptyName(t *testing.T) {
	_, err := NewStudent("s1", " ")

	if err != ErrEmptyStudentName {
		t.Errorf("expected ErrEmptyStudentName, got %v", err)
	}
}

func TestStudentRename(t *testing.T) {
	student, _ := NewStudent("s1", "Иван Иванов")

	err := student.Rename("Пётр Петров")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if student.Name() != "Пётр Петров" {
		t.Errorf("expected renamed student, got %s", student.Name())
	}

	err = student.Rename(" ")
	if err != ErrEmptyStudentName {
		t.Errorf("expected ErrEmptyStudentName, got %v", err)
	}
}
