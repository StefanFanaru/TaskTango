package main

// out of this datetime string 2024-06-13T20:52:41Z
// extract a smart time string like "10s", "5m", "2h", "3d", "1w", "2mo", "1y"
import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"time"
)

func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func SmartTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}

	d := time.Since(t)

	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}

	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}

	if d < time.Hour*24 {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}

	if d < time.Hour*24*7 {
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}

	if d < time.Hour*24*30 {
		return fmt.Sprintf("%dw", int(d.Hours()/24/7))
	}

	if d < time.Hour*24*365 {
		return fmt.Sprintf("%dmo", int(d.Hours()/24/30))
	}

	return fmt.Sprintf("%dy", int(d.Hours()/24/365))
}

func NoHighlightCell(s string) table.StyledCell {
	return table.NewStyledCell(s, lipgloss.NewStyle().Background(lipgloss.NoColor{}).Foreground(lipgloss.Color("#CDD6F4")))
}

func FindProjectById(selectedItemId int, projectRow []ProjectRow) ProjectRow {
	for _, project := range projectRow {
		if project.id == selectedItemId {
			return project
		}
	}
	return ProjectRow{}
}

func GetRowId(row table.Row) int {
	rowValue := row.Data[columnKeyID]
	return rowValue.(int)
}

func (m Model) GetTaskById(id int) Task {
	for _, task := range m.taskList {
		if task.id == id {
			return task
		}
	}
	return Task{}
}
