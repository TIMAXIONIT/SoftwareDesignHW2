package ports

import "pos/internal/domain"

type StudentRepository interface {
	SaveStudent(student *domain.Student) error
	FindStudentByID(id string) (*domain.Student, error)
	FindAllStudents() ([]*domain.Student, error)
}
