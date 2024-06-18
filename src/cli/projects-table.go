package main

import (
	"fmt"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyProjectTitle     = "title"
	columnKeyProjectTaskCount = "taskCount"
)

func generateProjectRowsFromData(projects []ProjectRow) []table.Row {
	taskRows := []table.Row{}
	for _, project := range projects {
		taskRows = append(taskRows, table.NewRow(table.RowData{
			columnKeyID:               project.id,
			columnKeyProjectTitle:     " " + project.title,
			columnKeyProjectTaskCount: NoHighlightCell(fmt.Sprint(project.taskCount)),
		}))
	}

	return taskRows
}

func buildProjectColumns() []table.Column {
	columns := []table.Column{
		table.NewColumn(columnKeyProjectTitle, " Title", 20),
		table.NewColumn(columnKeyProjectTaskCount, "Tasks", 7),
	}
	return columns
}
