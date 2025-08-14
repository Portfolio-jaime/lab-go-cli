package table

import (
	"strings"
)

type Table = SimpleTable

func NewTable(headers []string) *Table {
	return NewSimpleTable(headers)
}


func GetStatusColor(status string) int {
	return GetSimpleStatusColor(status)
}
