package application

import "pos/internal/domain"

type StudentReport struct {
	Student     *domain.Student
	Evaluations []*domain.Evaluation
	Result      *domain.Result
}

type Facade struct {
	service *Service
}

func NewFacade(service *Service) *Facade {
	return &Facade{
		service: service,
	}
}

func (f *Facade) CreateStudent(id, name string) (*domain.Student, error) {
	return f.service.CreateStudent(id, name)
}

func (f *Facade) CreateCriterion(
	id, name string,
	maxScore, weight float64,
) (*domain.Criterion, error) {
	return f.service.CreateCriterion(id, name, maxScore, weight)
}

func (f *Facade) CreateEvaluation(
	id, studentID, criterionID string,
	scoreValue float64,
	comment string,
) (*domain.Evaluation, error) {
	return f.service.CreateEvaluation(
		id,
		studentID,
		criterionID,
		scoreValue,
		comment,
	)
}

func (f *Facade) UpdateEvaluationScore(
	evaluationID string,
	scoreValue float64,
) error {
	return f.service.UpdateEvaluationScore(evaluationID, scoreValue)
}

func (f *Facade) UpdateEvaluationComment(
	evaluationID, comment string,
) error {
	return f.service.UpdateEvaluationComment(evaluationID, comment)
}

func (f *Facade) DeleteEvaluation(evaluationID string) error {
	return f.service.DeleteEvaluation(evaluationID)
}

func (f *Facade) GetStudentReport(studentID string) (*StudentReport, error) {
	student, err := f.service.GetStudent(studentID)
	if err != nil {
		return nil, err
	}

	evaluations, err := f.service.GetStudentEvaluations(studentID)
	if err != nil {
		return nil, err
	}

	result, err := f.service.CalculateStudentResult(studentID)
	if err != nil {
		return nil, err
	}

	return &StudentReport{
		Student:     student,
		Evaluations: evaluations,
		Result:      result,
	}, nil
}

func (f *Facade) GetAllStudents() ([]*domain.Student, error) {
	return f.service.GetAllStudents()
}

func (f *Facade) GetAllCriteria() ([]*domain.Criterion, error) {
	return f.service.GetAllCriteria()
}

func (f *Facade) GetAllStudentReports() ([]*StudentReport, error) {
	return f.service.GetAllStudentReports()
}
