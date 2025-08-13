package table

import (
	"fmt"
	"strings"
)

type SimpleTable struct {
	headers []string
	rows    [][]string
}

func NewSimpleTable(headers []string) *SimpleTable {
	return &SimpleTable{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

func (t *SimpleTable) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

func (t *SimpleTable) AddRowWithColors(row []string, colors []int) {
	t.rows = append(t.rows, row)
}

func (t *SimpleTable) Render() {
	if len(t.headers) == 0 {
		return
	}

	colWidths := make([]int, len(t.headers))
	for i, header := range t.headers {
		colWidths[i] = len(header)
	}

	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	printSeparator(colWidths)
	printRow(t.headers, colWidths, true)
	printSeparator(colWidths)

	for _, row := range t.rows {
		printRow(row, colWidths, false)
	}

	printSeparator(colWidths)
}

func printSeparator(colWidths []int) {
	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2))
		fmt.Print("+")
	}
	fmt.Println()
}

func printRow(row []string, colWidths []int, isHeader bool) {
	fmt.Print("|")
	for i, cell := range row {
		if i < len(colWidths) {
			if isHeader {
				fmt.Printf(" %-*s ", colWidths[i], strings.ToUpper(cell))
			} else {
				fmt.Printf(" %-*s ", colWidths[i], cell)
			}
			fmt.Print("|")
		}
	}
	fmt.Println()
}

func GetSimpleStatusColor(status string) int {
	return 0 // No color for simple version
}
