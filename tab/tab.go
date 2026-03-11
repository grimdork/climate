// Package tab provides simple functions for tabulated output on the terminal.
package tab

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"strings"
	"text/tabwriter"
)

// Tabulate takes the input string, splits it into lines and columns, and returns it in a tabulated string.
// If twoColMode is true, it treats the first word as the first column and the rest of the line as the second column.
func Tabulate(input string, twoColMode bool) (string, error) {
	rows, err := SplitColumns(input, twoColMode)
	if err != nil {
		return "", err
	}

	return formatRows(rows)
}

// TabulateCSV parses CSV input and returns it as aligned columns.
// The first row is treated as a header and separated from the data by a dashed line.
func TabulateCSV(input string) (string, error) {
	r := csv.NewReader(strings.NewReader(input))
	rows, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "", nil
	}

	// Find max width per column across all rows
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Build separator dashes matching the widest cell in each column
	dashes := make([]string, len(widths))
	for i, w := range widths {
		dashes[i] = strings.Repeat("-", w)
	}

	var b strings.Builder
	w := tabwriter.NewWriter(&b, 1, 0, 2, ' ', 0)

	fmt.Fprintln(w, strings.Join(rows[0], "\t"))
	fmt.Fprintln(w, strings.Join(dashes, "\t"))

	for _, row := range rows[1:] {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
	return b.String(), nil
}

// SplitColumns reads lines from the input string, splits them into columns based on whitespace, and returns the rows.
// If twoColMode is true, it treats the first word as the first column and the rest of the line as the second column.
func SplitColumns(input string, twoColMode bool) ([][]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var rows [][]string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var fields []string
		if twoColMode {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				if len(parts) > 1 {
					fields = []string{parts[0], strings.Join(parts[1:], " ")}
				} else {
					fields = []string{parts[0], ""}
				}
			}
		} else {
			fields = strings.Fields(line)
		}

		if len(fields) > 0 {
			rows = append(rows, fields)
		}
	}

	return rows, scanner.Err()
}

func formatRows(rows [][]string) (string, error) {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 1, 0, 2, ' ', 0)

	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
	return b.String(), nil
}
