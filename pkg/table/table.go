package table

import (
	"strings"
)

type Table = SimpleTable

func NewTable(headers []string) *Table {
	return NewSimpleTable(headers)
}

func formatCell(cell string) string {
	cell = strings.TrimSpace(cell)
	if len(cell) > 50 {
		return cell[:47] + "..."
	}
	return cell
}

func GetStatusColor(status string) int {
	return GetSimpleStatusColor(status)
}