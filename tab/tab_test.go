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
