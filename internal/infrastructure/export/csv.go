package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"pos/internal/application"
	"pos/internal/domain"
)

type CSVExporter struct {
	app *application.Facade
}

func NewCSVExporter(app *application.Facade) *CSVExporter {
	return &CSVExporter{
		app: app,
	}
}

func (e *CSVExporter) Export(path string) error {
	reports, err := e.app.GetAllStudentReports()
	if err != nil {
		return err
	}

	criteria, err := e.app.GetAllCriteria()
	if err != nil {
		return err
	}

	criteriaByID := make(map[string]*domain.Criterion, len(criteria))
	for _, criterion := range criteria {
		criteriaByID[criterion.ID()] = criterion
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"student_id",
		"student_name",
		"criterion_id",
		"criterion_name",
		"score",
		"max_score",
		"weight",
		"comment",
		"result",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, report := range reports {
		for _, evaluation := range report.Evaluations {
			criterion := criteriaByID[evaluation.CriterionID()]
			if criterion == nil {
				return fmt.Errorf(
					"criterion %q not found",
					evaluation.CriterionID(),
				)
			}

			record := []string{
				report.Student.ID(),
				report.Student.Name(),
				criterion.ID(),
				criterion.Name(),
				strconv.FormatFloat(evaluation.Score().Value(), 'f', 2, 64),
				strconv.FormatFloat(criterion.MaxScore(), 'f', 2, 64),
				strconv.FormatFloat(criterion.Weight(), 'f', 2, 64),
				evaluation.Comment(),
				strconv.FormatFloat(report.Result.Score(), 'f', 2, 64),
			}

			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}

	writer.Flush()
	return writer.Error()
}
