package tab_test

import (
	"testing"

	"github.com/grimdork/climate/tab"
)

func TestTabulate(t *testing.T) {
	input := "apple banana cherry\ndate elderberry"

	// tabwriter default (tabwidth 8, padding 1):
	// "apple" (5) + 3 spaces to reach tab stop 8, then 1 padding space = 4 spaces total.
	// "date" (4) + 4 spaces to reach tab stop 8, then 1 padding space = 5 spaces total.
	expected := `apple  banana  cherry
date   elderberry
`

	result, err := tab.Tabulate(input, false)
	if err != nil {
		t.Fatalf("Tabulate failed: %v", err)
	}

	if result != expected {
		t.Errorf("Result mismatch.\nGot:  %q\nWant: %q", result, expected)
	}
}

func TestTabulateTwoColMode(t *testing.T) {
	input := "user1 Read write delete\nuser2 Read only"

	// Column 1 max width is "user1" (5).
	// To reach the next tab stop (8) + 1 padding, we expect 4 spaces after "user1".
	expected := `user1  Read write delete
user2  Read only
`

	result, err := tab.Tabulate(input, true)
	if err != nil {
		t.Fatalf("Tabulate failed: %v", err)
	}

	if result != expected {
		t.Errorf("Result mismatch.\nGot:  %q\nWant: %q", result, expected)
	}
}

func TestTabulateCSV(t *testing.T) {
	input := "Name,Role,Status\nAlice,Admin,Active\nBob,User,Offline"
	expected := `Name   Role   Status
-----  -----  -------
Alice  Admin  Active
Bob    User   Offline
`

	result, err := tab.TabulateCSV(input)
	if err != nil {
		t.Fatalf("TabulateCSV failed: %v", err)
	}

	if result != expected {
		t.Errorf("Result mismatch.\nGot:  %q\nWant: %q", result, expected)
	}
}

func TestTabulateCSVQuoted(t *testing.T) {
	input := `Name,Role
"Eve, Jr",Guest
Alice,Admin`

	result, err := tab.TabulateCSV(input)
	if err != nil {
		t.Fatalf("TabulateCSV failed: %v", err)
	}

	// Quoted commas should be preserved in the cell
	if result == "" {
		t.Fatal("Expected non-empty output")
	}

	// Verify "Eve, Jr" appears intact
	if !containsStr(result, "Eve, Jr") {
		t.Errorf("Expected 'Eve, Jr' in output, got:\n%s", result)
	}
}

func TestTabulateTSV(t *testing.T) {
	input := "Name\tRole\nAlice\tAdmin\nBob\tUser"
	expected := `Name   Role
Alice  Admin
Bob    User
`

	result, err := tab.Tabulate(input, false)
	if err != nil {
		t.Fatalf("Tabulate with TSV failed: %v", err)
	}

	if result != expected {
		t.Errorf("Result mismatch.\nGot:  %q\nWant: %q", result, expected)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
