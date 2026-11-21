package console

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"pos/internal/application"
	"pos/internal/infrastructure/export"
)

type Console struct {
	app      *application.Facade
	reader   *bufio.Scanner
	writer   io.Writer
	exporter *export.CSVExporter
}

func New(
	app *application.Facade,
	exporter *export.CSVExporter,
	input io.Reader,
	output io.Writer,
) *Console {
	return &Console{
		app:      app,
		exporter: exporter,
		reader:   bufio.NewScanner(input),
		writer:   output,
	}
}

func (c *Console) readLine(prompt string) (string, error) {
	fmt.Fprint(c.writer, prompt)

	if !c.reader.Scan() {
		if err := c.reader.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}

	return c.reader.Text(), nil
}

func (c *Console) readFloat(prompt string) (float64, error) {
	value, err := c.readLine(prompt)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(value, 64)
}

func (c *Console) Run() {
	for {
		fmt.Fprintln(c.writer, "\n=== ПОС ===")
		fmt.Fprintln(c.writer, "1. Создать студента")
		fmt.Fprintln(c.writer, "2. Показать студентов")
		fmt.Fprintln(c.writer, "3. Создать критерий")
		fmt.Fprintln(c.writer, "4. Показать критерии")
		fmt.Fprintln(c.writer, "5. Выставить оценку")
		fmt.Fprintln(c.writer, "6. Показать отчёт студента")
		fmt.Fprintln(c.writer, "7. Изменить оценку")
		fmt.Fprintln(c.writer, "8. Изменить комментарий")
		fmt.Fprintln(c.writer, "9. Отменить оценку")
		fmt.Fprintln(c.writer, "10. Экспортировать результаты в CSV")
		fmt.Fprintln(c.writer, "0. Выход")

		choice, err := c.readLine("Выберите действие: ")
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintln(c.writer, "Ошибка ввода:", err)
			return
		}

		switch choice {
		case "1":
			c.createStudent()
		case "2":
			c.listStudents()
		case "3":
			c.createCriterion()
		case "4":
			c.listCriteria()
		case "5":
			c.createEvaluation()
		case "6":
			c.showReport()
		case "7":
			c.updateScore()
		case "8":
			c.updateComment()
		case "9":
			c.deleteEvaluation()
		case "10":
			c.exportReports()
		case "0":
			fmt.Fprintln(c.writer, "Завершение работы.")
			return
		default:
			fmt.Fprintln(c.writer, "Неизвестная команда.")
		}
	}
}

func (c *Console) createStudent() {
	id, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	name, err := c.readLine("Имя студента: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	student, err := c.app.CreateStudent(id, name)
	if err != nil {
		fmt.Fprintln(c.writer, "Не удалось создать студента:", err)
		return
	}

	fmt.Fprintf(c.writer, "Студент создан: %s (%s)\n", student.Name(), student.ID())
}

func (c *Console) listStudents() {
	students, err := c.app.GetAllStudents()
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка:", err)
		return
	}

	if len(students) == 0 {
		fmt.Fprintln(c.writer, "Список студентов пуст.")
		return
	}

	for _, student := range students {
		fmt.Fprintf(c.writer, "%s — %s\n", student.ID(), student.Name())
	}
}

func (c *Console) createCriterion() {
	id, err := c.readLine("ID критерия: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	name, err := c.readLine("Название критерия: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	maxScore, err := c.readFloat("Максимальный балл: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Некорректное число:", err)
		return
	}

	weight, err := c.readFloat("Вес критерия: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Некорректное число:", err)
		return
	}

	criterion, err := c.app.CreateCriterion(id, name, maxScore, weight)
	if err != nil {
		fmt.Fprintln(c.writer, "Не удалось создать критерий:", err)
		return
	}

	fmt.Fprintf(c.writer, "Критерий создан: %s (%s)\n", criterion.Name(), criterion.ID())
}

func (c *Console) listCriteria() {
	criteria, err := c.app.GetAllCriteria()
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка:", err)
		return
	}

	if len(criteria) == 0 {
		fmt.Fprintln(c.writer, "Список критериев пуст.")
		return
	}

	for _, criterion := range criteria {
		fmt.Fprintf(
			c.writer,
			"%s — %s (максимум: %.2f, вес: %.2f)\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore(),
			criterion.Weight(),
		)
	}
}

func (c *Console) createEvaluation() {
	id, err := c.readLine("ID оценки: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	criterionID, err := c.readLine("ID критерия: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	score, err := c.readFloat("Балл: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Некорректное число:", err)
		return
	}

	comment, err := c.readLine("Комментарий: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	_, err = c.app.CreateEvaluation(id, studentID, criterionID, score, comment)
	if err != nil {
		fmt.Fprintln(c.writer, "Не удалось выставить оценку:", err)
		return
	}

	fmt.Fprintln(c.writer, "Оценка успешно выставлена.")
}

func (c *Console) showReport() {
	studentID, err := c.readLine("ID студента: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	report, err := c.app.GetStudentReport(studentID)
	if err != nil {
		fmt.Fprintln(c.writer, "Не удалось получить отчёт:", err)
		return
	}

	fmt.Fprintf(c.writer, "\nОтчёт по студенту: %s (%s)\n",
		report.Student.Name(), report.Student.ID())

	fmt.Fprintln(c.writer, "Оценки:")
	for _, evaluation := range report.Evaluations {
		fmt.Fprintf(
			c.writer,
			"Критерий %s: %.2f, комментарий: %s\n",
			evaluation.CriterionID(),
			evaluation.Score().Value(),
			evaluation.Comment(),
		)
	}

	fmt.Fprintf(c.writer, "Итоговый результат: %.2f / 10\n", report.Result.Score())
}

func (c *Console) updateScore() {
	id, err := c.readLine("ID оценки: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	score, err := c.readFloat("Новая оценка: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Некорректное число:", err)
		return
	}

	if err := c.app.UpdateEvaluationScore(id, score); err != nil {
		fmt.Fprintln(c.writer, "Не удалось изменить оценку:", err)
		return
	}

	fmt.Fprintln(c.writer, "Оценка изменена.")
}

func (c *Console) updateComment() {
	id, err := c.readLine("ID оценки: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	comment, err := c.readLine("Новый комментарий: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	if err := c.app.UpdateEvaluationComment(id, comment); err != nil {
		fmt.Fprintln(c.writer, "Не удалось изменить комментарий:", err)
		return
	}

	fmt.Fprintln(c.writer, "Комментарий изменён.")
}

func (c *Console) deleteEvaluation() {
	id, err := c.readLine("ID оценки: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	if err := c.app.DeleteEvaluation(id); err != nil {
		fmt.Fprintln(c.writer, "Не удалось отменить оценку:", err)
		return
	}

	fmt.Fprintln(c.writer, "Оценка отменена.")
}

func (c *Console) exportReports() {
	path, err := c.readLine("Путь для сохранения CSV: ")
	if err != nil {
		fmt.Fprintln(c.writer, "Ошибка ввода:", err)
		return
	}

	if err := c.exporter.Export(path); err != nil {
		fmt.Fprintln(c.writer, "Ошибка экспорта:", err)
		return
	}

	fmt.Fprintln(c.writer, "Результаты успешно экспортированы.")
}
