package tab_test

import (
	"testing"

	"github.com/grimdork/climate/tab"
)

func TestSplitColumnsTSVPreservesEmpties(t *testing.T) {
	input := "a\tb\tc\n1\t\t3\n4\t5\t\n"
	rows, err := tab.SplitColumnsTSV(input)
	if err != nil {
		t.Fatalf("SplitColumnsTSV failed: %v", err)
	}

	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}

	expected := [][]string{
		{"a", "b", "c"},
		{"1", "", "3"},
		{"4", "5", ""},
	}

	for i := range expected {
		if len(rows[i]) != len(expected[i]) {
			t.Fatalf("row %d: expected %d columns, got %d", i, len(expected[i]), len(rows[i]))
		}
		for j := range expected[i] {
			if rows[i][j] != expected[i][j] {
				t.Fatalf("row %d col %d: expected %q, got %q", i, j, expected[i][j], rows[i][j])
			}
		}
	}
}

func TestTabulate(t *testing.T) {
	input := "a b c\nd e f"
	got, err := tab.Tabulate(input, false)
	if err != nil {
		t.Fatalf("Tabulate failed: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTabulateTwoColMode(t *testing.T) {
	input := "name Alice\nrole Engineer"
	got, err := tab.Tabulate(input, true)
	if err != nil {
		t.Fatalf("Tabulate twoCol failed: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTabulateEmpty(t *testing.T) {
	got, err := tab.Tabulate("", false)
	if err != nil {
		t.Fatalf("Tabulate empty failed: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}

func TestTabulateCSV(t *testing.T) {
	input := "name,role,dept\nAlice,Engineer,IT\nBob,Manager,Sales"
	got, err := tab.TabulateCSV(input)
	if err != nil {
		t.Fatalf("TabulateCSV failed: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTabulateCSVEmpty(t *testing.T) {
	got, err := tab.TabulateCSV("")
	if err != nil {
		t.Fatalf("TabulateCSV empty failed: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty output, got %q", got)
	}
}

func TestTabulateCSVVariableColumns(t *testing.T) {
	input := "a,b\n1,2\n3,4"
	got, err := tab.TabulateCSV(input)
	if err != nil {
		t.Fatalf("TabulateCSV variable cols failed: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTabulateTSV(t *testing.T) {
	input := "name\trole\nAlice\tEngineer"
	got, err := tab.TabulateTSV(input)
	if err != nil {
		t.Fatalf("TabulateTSV failed: %v", err)
	}
	if got == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestSplitColumns(t *testing.T) {
	rows, err := tab.SplitColumns("a b c\nd e f", false)
	if err != nil {
		t.Fatalf("SplitColumns failed: %v", err)
	}
	if len(rows) != 2 || len(rows[0]) != 3 {
		t.Fatalf("expected 2x3, got %dx%d", len(rows), len(rows[0]))
	}
}

func TestSplitColumnsTwoColMode(t *testing.T) {
	rows, err := tab.SplitColumns("name Alice\nrole Engineer", true)
	if err != nil {
		t.Fatalf("SplitColumns twoCol failed: %v", err)
	}
	if len(rows) != 2 || len(rows[0]) != 2 {
		t.Fatalf("expected 2x2, got %dx%d", len(rows), len(rows[0]))
	}
	if rows[0][0] != "name" || rows[0][1] != "Alice" {
		t.Fatalf("expected name/Alice, got %q/%q", rows[0][0], rows[0][1])
	}
}

func TestSplitColumnsBlankLines(t *testing.T) {
	rows, err := tab.SplitColumns("a b\n\nc d", false)
	if err != nil {
		t.Fatalf("SplitColumns failed: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
}

func TestSplitColumnsTSV(t *testing.T) {
	input := "a\tb\n1\t2"
	rows, err := tab.SplitColumnsTSV(input)
	if err != nil {
		t.Fatalf("SplitColumnsTSV failed: %v", err)
	}
	if len(rows) != 2 || len(rows[0]) != 2 {
		t.Fatalf("expected 2x2, got %dx%d", len(rows), len(rows[0]))
	}
}

func TestSplitColumnsTSVBlankLines(t *testing.T) {
	input := "a\tb\n\nc\td"
	rows, err := tab.SplitColumnsTSV(input)
	if err != nil {
		t.Fatalf("SplitColumnsTSV failed: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
}
