// Package tab provides simple functions for tabulated output on the terminal.
package tab

import (
	"bufio"
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

	var b strings.Builder
	// minwidth=1, tabwidth=0, padding=2, padchar=' ', flags=0
	w := tabwriter.NewWriter(&b, 1, 0, 2, ' ', 0)

	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
	return b.String(), nil
}

// SplitColumns reads lines fromt the input string, splits them into columns based on whitespace, and returns the rows.
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
