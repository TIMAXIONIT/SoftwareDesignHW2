package export

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"pos/internal/application"
	"pos/internal/domain"
	"pos/internal/infrastructure/persistence"
)

func createTestExporter(t *testing.T) *CSVExporter {
	t.Helper()

	repo := persistence.NewInMemoryRepository()

	calculator := domain.NewResultDomainService(
		domain.WeightedAverageStrategy{},
	)

	factory := application.NewEvaluationFactory(repo, repo)

	service := application.NewService(
		repo,
		repo,
		repo,
		calculator,
		factory,
	)

	app := application.NewFacade(service)

	return NewCSVExporter(app)
}

func TestCSVExporterExport(t *testing.T) {
	exporter := createTestExporter(t)

	if _, err := exporter.app.CreateStudent("student-1", "Ivan Ivanov"); err != nil {
		t.Fatalf("failed to create student: %v", err)
	}

	if _, err := exporter.app.CreateCriterion("criterion-1", "Programming", 10, 2); err != nil {
		t.Fatalf("failed to create criterion: %v", err)
	}

	if _, err := exporter.app.CreateEvaluation(
		"evaluation-1",
		"student-1",
		"criterion-1",
		8,
		"Good work",
	); err != nil {
		t.Fatalf("failed to create evaluation: %v", err)
	}

	path := filepath.Join(t.TempDir(), "report.csv")

	if err := exporter.Export(path); err != nil {
		t.Fatalf("failed to export CSV: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("CSV file was not created: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	expectedHeader := []string{
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

	if !reflect.DeepEqual(records[0], expectedHeader) {
		t.Fatalf("unexpected CSV header:\nexpected: %v\ngot: %v", expectedHeader, records[0])
	}

	if len(records) != 2 {
		t.Fatalf("expected 2 CSV records (header + 1 evaluation), got %d", len(records))
	}

	row := records[1]

	expectedFields := []string{
		"student-1",
		"Ivan Ivanov",
		"criterion-1",
		"Programming",
		"8.00",
		"10.00",
		"2.00",
		"Good work",
	}

	for i, expected := range expectedFields {
		if row[i] != expected {
			t.Errorf("unexpected value in column %d: expected %q, got %q", i, expected, row[i])
		}
	}

	result, err := strconv.ParseFloat(row[8], 64)
	if err != nil {
		t.Fatalf("invalid result value in CSV: %v", err)
	}

	if result != 8 {
		t.Errorf("expected result 8, got %f", result)
	}
}
