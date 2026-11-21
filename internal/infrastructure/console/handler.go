package console

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"pos/internal/application"
)

type Handler struct {
	service *application.Service
	scanner *bufio.Scanner
	output  io.Writer
}

func NewHandler(service *application.Service, input io.Reader, output io.Writer) *Handler {
	return &Handler{
		service: service,
		scanner: bufio.NewScanner(input),
		output:  output,
	}
}

func (h *Handler) readLine(prompt string) (string, error) {
	fmt.Fprint(h.output, prompt)

	if !h.scanner.Scan() {
		if err := h.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}

	return strings.TrimSpace(h.scanner.Text()), nil
}

func (h *Handler) readFloat(prompt string) (float64, error) {
	value, err := h.readLine(prompt)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(value, 64)
}

func (h *Handler) Run() error {
	for {
		fmt.Fprintln(h.output, "\n=== ПОС ===")
		fmt.Fprintln(h.output, "1. Создать студента")
		fmt.Fprintln(h.output, "2. Создать критерий")
		fmt.Fprintln(h.output, "3. Выставить оценку")
		fmt.Fprintln(h.output, "4. Показать студентов")
		fmt.Fprintln(h.output, "5. Показать критерии")
		fmt.Fprintln(h.output, "6. Показать оценки студента")
		fmt.Fprintln(h.output, "7. Рассчитать итоговую оценку")
		fmt.Fprintln(h.output, "8. Изменить оценку")
		fmt.Fprintln(h.output, "9. Изменить комментарий")
		fmt.Fprintln(h.output, "10. Удалить оценку")
		fmt.Fprintln(h.output, "0. Выход")

		choice, err := h.readLine("Выберите действие: ")
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		var actionErr error

		switch choice {
		case "1":
			actionErr = h.createStudent()
		case "2":
			actionErr = h.createCriterion()
		case "3":
			actionErr = h.createEvaluation()
		case "4":
			actionErr = h.showStudents()
		case "5":
			actionErr = h.showCriteria()
		case "6":
			actionErr = h.showEvaluations()
		case "7":
			actionErr = h.showResult()
		case "8":
			actionErr = h.updateScore()
		case "9":
			actionErr = h.updateComment()
		case "10":
			actionErr = h.deleteEvaluation()
		case "0":
			fmt.Fprintln(h.output, "Работа завершена.")
			return nil
		default:
			fmt.Fprintln(h.output, "Неизвестная команда.")
		}

		if actionErr != nil {
			fmt.Fprintf(h.output, "Ошибка: %v\n", actionErr)
		}
	}
}

func (h *Handler) createStudent() error {
	id, err := h.readLine("ID студента: ")
	if err != nil {
		return err
	}

	name, err := h.readLine("Имя студента: ")
	if err != nil {
		return err
	}

	student, err := h.service.CreateStudent(id, name)
	if err != nil {
		return err
	}

	fmt.Fprintf(h.output, "Студент создан: %s (%s)\n", student.Name(), student.ID())
	return nil
}

func (h *Handler) createCriterion() error {
	id, err := h.readLine("ID критерия: ")
	if err != nil {
		return err
	}

	name, err := h.readLine("Название критерия: ")
	if err != nil {
		return err
	}

	maxScore, err := h.readFloat("Максимальный балл: ")
	if err != nil {
		return err
	}

	weight, err := h.readFloat("Вес критерия: ")
	if err != nil {
		return err
	}

	criterion, err := h.service.CreateCriterion(id, name, maxScore, weight)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		h.output,
		"Критерий создан: %s (максимум: %.2f, вес: %.2f)\n",
		criterion.Name(),
		criterion.MaxScore(),
		criterion.Weight(),
	)
	return nil
}

func (h *Handler) createEvaluation() error {
	id, err := h.readLine("ID оценки: ")
	if err != nil {
		return err
	}

	studentID, err := h.readLine("ID студента: ")
	if err != nil {
		return err
	}

	criterionID, err := h.readLine("ID критерия: ")
	if err != nil {
		return err
	}

	score, err := h.readFloat("Выставленный балл: ")
	if err != nil {
		return err
	}

	comment, err := h.readLine("Комментарий: ")
	if err != nil {
		return err
	}

	evaluation, err := h.service.CreateEvaluation(
		id,
		studentID,
		criterionID,
		score,
		comment,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		h.output,
		"Оценка выставлена. ID: %s, балл: %.2f\n",
		evaluation.ID(),
		evaluation.Score().Value(),
	)
	return nil
}

func (h *Handler) showStudents() error {
	students, err := h.service.GetAllStudents()
	if err != nil {
		return err
	}

	if len(students) == 0 {
		fmt.Fprintln(h.output, "Студентов пока нет.")
		return nil
	}

	for _, student := range students {
		fmt.Fprintf(h.output, "%s — %s\n", student.ID(), student.Name())
	}
	return nil
}

func (h *Handler) showCriteria() error {
	criteria, err := h.service.GetAllCriteria()
	if err != nil {
		return err
	}

	if len(criteria) == 0 {
		fmt.Fprintln(h.output, "Критериев пока нет.")
		return nil
	}

	for _, criterion := range criteria {
		fmt.Fprintf(
			h.output,
			"%s — %s (максимум: %.2f, вес: %.2f)\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore(),
			criterion.Weight(),
		)
	}
	return nil
}

func (h *Handler) showEvaluations() error {
	studentID, err := h.readLine("ID студента: ")
	if err != nil {
		return err
	}

	evaluations, err := h.service.GetStudentEvaluations(studentID)
	if err != nil {
		return err
	}

	if len(evaluations) == 0 {
		fmt.Fprintln(h.output, "У этого студента пока нет оценок.")
		return nil
	}

	for _, evaluation := range evaluations {
		fmt.Fprintf(
			h.output,
			"Оценка %s | критерий: %s | балл: %.2f | комментарий: %s\n",
			evaluation.ID(),
			evaluation.CriterionID(),
			evaluation.Score().Value(),
			evaluation.Comment(),
		)
	}
	return nil
}

func (h *Handler) showResult() error {
	studentID, err := h.readLine("ID студента: ")
	if err != nil {
		return err
	}

	result, err := h.service.CalculateStudentResult(studentID)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		h.output,
		"Итоговая оценка студента %s: %.2f из 10\n",
		result.StudentID(),
		result.Score(),
	)
	return nil
}

func (h *Handler) updateScore() error {
	id, err := h.readLine("ID оценки: ")
	if err != nil {
		return err
	}

	score, err := h.readFloat("Новый балл: ")
	if err != nil {
		return err
	}

	if err := h.service.UpdateEvaluationScore(id, score); err != nil {
		return err
	}

	fmt.Fprintln(h.output, "Оценка обновлена.")
	return nil
}

func (h *Handler) updateComment() error {
	id, err := h.readLine("ID оценки: ")
	if err != nil {
		return err
	}

	comment, err := h.readLine("Новый комментарий: ")
	if err != nil {
		return err
	}

	if err := h.service.UpdateEvaluationComment(id, comment); err != nil {
		return err
	}

	fmt.Fprintln(h.output, "Комментарий обновлён.")
	return nil
}

func (h *Handler) deleteEvaluation() error {
	id, err := h.readLine("ID оценки: ")
	if err != nil {
		return err
	}

	if err := h.service.DeleteEvaluation(id); err != nil {
		return err
	}

	fmt.Fprintln(h.output, "Оценка удалена.")
	return nil
}
