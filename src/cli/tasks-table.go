package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeySelect      = "select"
	columnKeyID          = "id"
	columnKeyDescription = "description"
	columnKeyProject     = "project"
	columnKeyTags        = "tags"
	columnKeyAge         = "age"
	columnKeyPriority    = "priority"
)

func generateTaskRowsFromData(tasks []Task) []table.Row {
	taskRows := []table.Row{}
	priorityStyles := map[Priority]lipgloss.Style{}
	priorityStyles[Low] = lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1"))
	priorityStyles[Medium] = lipgloss.NewStyle().Foreground(lipgloss.Color("#F9E2AF"))
	priorityStyles[High] = lipgloss.NewStyle().Foreground(lipgloss.Color("#F5C2E7"))
	priorityStyles[Critical] = lipgloss.NewStyle().Foreground(lipgloss.Color("#F38BA8"))

	for _, task := range tasks {
		row := table.NewRow(table.RowData{
			columnKeyID:          task.id,
			columnKeyDescription: " " + task.description,
			columnKeyProject:     NoHighlightCell(task.projectTitle),
			columnKeyTags:        NoHighlightCell(task.tagsJoined),
			columnKeyAge:         NoHighlightCell(SmartTime(task.createdDate)),
			columnKeyPriority:    table.NewStyledCell(task.priority.String(), priorityStyles[task.priority].Background(lipgloss.NoColor{})),
		})

		if task.isStarted {
			row.Data[columnKeyDescription] = table.NewStyledCell(" "+task.description, lipgloss.NewStyle().Background(lipgloss.Color("#313244")).Bold(true))
		}

		taskRows = append(taskRows, row)
	}

	return taskRows
}

func buildTaskColumns() []table.Column {
	columns := []table.Column{
		table.NewFlexColumn(columnKeyDescription, " Description", 4),
		table.NewFlexColumn(columnKeyProject, "Project", 1),
		table.NewFlexColumn(columnKeyTags, "Tags", 1),
		table.NewColumn(columnKeyAge, "Age", 4),
		table.NewColumn(columnKeyPriority, "Prio", 4),
	}
	return columns
}
