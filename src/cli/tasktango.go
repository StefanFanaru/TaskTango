// This file contains a full demo of most available features, for both testing

package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

func NewModel() Model {
	rows := []table.Row{}
	model := buildModel(rows)
	model.updateFooter()
	model.undoStack = UndoStack{}
	model.undoStack.recentlyDeleted = make(map[int]string)
	model.filter = AppState.Filter
	return model
}

func buildTableModel(rows []table.Row, viewMode string) table.Model {
	columns := []table.Column{}
	if viewMode == "tasks" {
		columns = buildTaskColumns()
	} else if viewMode == "projects" {
		columns = buildProjectColumns()
	}
	keys := table.DefaultKeyMap()
	return table.New(columns).
		WithRows(rows).
		HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1")).Bold(true)).
		HighlightStyle(lipgloss.NewStyle().Background(lipgloss.Color("#1E1E2E"))).
		Focused(true).
		Border(CustomBorder).
		WithKeyMap(keys).
		WithBaseStyle(
			lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#585B70")).
				Foreground(lipgloss.Color("#CDD6F4")).
				Align(lipgloss.Left),
		)
}

func buildModel(rows []table.Row) Model {
	model := Model{
		tableModel:   buildTableModel(rows, "tasks"),
		commandInput: BuildTextInput(),
		editForm:     BuildTaskEditForm(Task{}, Tags{}, []ProjectRow{}),
		viewMode:     "tasks",
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(initTaskTango())
}

func (m *Model) updateFooter() {
	footerText := ""
	if m.tablePageSize >= len(m.taskList) {
		footerText = fmt.Sprintf(m.commandInput.View())
	} else {
		footerText = fmt.Sprintf(
			"Page %d/%d \n%s",
			m.tableModel.CurrentPage(),
			m.tableModel.MaxPages(),
			m.commandInput.View(),
		)
	}

	m.tableModel = m.tableModel.WithStaticFooter(footerText)
}

func (m *Model) updateCommandInput(isCommandInit bool, msg tea.Msg, cmds []tea.Cmd) []tea.Cmd {
	if m.commandInputState.isActive && !isCommandInit {
		inputModel, cmd := m.commandInput.Update(msg)
		m.commandInput = inputModel
		if m.commandInputState.activeCommand == "search" {
			if m.commandInput.Value() != AppState.Filter.Search {
				AppState.Filter.Search = m.commandInput.Value()
				cmds = append(cmds, queryTaskList)
			}
		}
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (m *Model) updateEditForm(isFormInit bool, msg tea.Msg, cmds []tea.Cmd) []tea.Cmd {
	if m.editFormState.isFormVisible && !isFormInit {
		form, cmd := m.editForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.editForm = f
			if m.editForm.State == huh.StateCompleted {
				m.editFormState.isFormVisible = false
				cmds = append(cmds, handleFormCompletion(m))
			}
		}
		cmds = append(cmds, cmd)
	}

	return cmds
}
func (m *Model) updateTable(cmds []tea.Cmd, msg tea.Msg) []tea.Cmd {
	if !m.commandInputState.isActive && !m.editFormState.isFormVisible {
		table, cmd := m.tableModel.Update(msg)
		m.tableModel = table
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd           tea.Cmd
		cmds          []tea.Cmd
		isCommandInit bool
		isFormInit    bool
	)

	cmds = m.updateTable(cmds, msg)
	cmds = m.updateEditForm(isFormInit, msg, cmds)
	cmds = m.updateCommandInput(isCommandInit, msg, cmds)
	m.updateFooter()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		msgString := msg.String()

		if msgString == "ctrl+c" || msgString == "esc" {
			if m.editFormState.isFormVisible {
				m.editFormState.isFormVisible = false
			}
			cmds = append(cmds, tea.Quit)
			return m, tea.Batch(cmds...)
		}

		if m.commandInputState.isActive {
			switch msgString {
			case "ctrl+c", "esc":
				m.ExitTextInput()
				m.tableModel, cmd = m.tableModel.Update(msg)
				cmds = append(cmds, cmd)
			case "enter":
				if m.commandInput.Value() == "q" {
					cmds = append(cmds, tea.Quit)
					break
				}
				cmds = append(cmds, m.HandleCommandInput(m.commandInput.Value(), m.commandInputState.activeCommand, m.viewMode))
				m.ExitTextInput()
				m.tableModel, cmd = m.tableModel.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			if !m.editFormState.isFormVisible {
				switch msgString {
				case "q":
					cmds = append(cmds, tea.Quit)
				case "r":
					cmds = append(cmds, queryTaskList)
				case "d":
					m.pushDeletionStack()
					cmds = append(cmds, handleDeletion(m))
				case "u":
					if len(m.undoStack.recentlyDeleted) > 0 {
						id := m.popDeletionStack()
						cmds = append(cmds, handleUndelete(id, m.viewMode))
					}
				case "e":
					cmds = m.handleItemEdit(isFormInit, cmds)
				case "a":
					m.PrepareCommandInput("add", "Add: ")
					cmds = append(cmds, m.commandInput.Focus())
					isCommandInit = true
				case "t":
					m.switchView("tasks")
					cmds = append(cmds, queryTaskList)
				case "p":
					m.switchView("projects")
					cmds = append(cmds, queryProjectList)
				case "s":
					selectedItemId := GetRowId(m.tableModel.HighlightedRow())
					if m.viewMode == "tasks" {
						cmds = append(cmds, startTask(selectedItemId))
					}
				case ":":
					m.PrepareCommandInput("", ":")
					cmds = append(cmds, m.commandInput.Focus())
					isCommandInit = true
				case "/":
					m.PrepareCommandInput("search", "Search: ")
					m.commandInput.SetValue(AppState.Filter.Search)
					cmds = append(cmds, m.commandInput.Focus())
					isCommandInit = true
				}
			}
		}
	case queryAllTagsMsg:
		m.allTags = Tags(msg)
	case queryTaskTagsMsg:
		m.handleQueryTagsResult(msg)
		isFormInit = true
		cmds = append(cmds, m.editForm.Init())
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height
		m.recalculateTable()
	case queryProjectListMsg:
		m.projectList = msg
		if m.viewMode == "projects" {
			m.tableModel = m.tableModel.WithRows(generateProjectRowsFromData(msg))
		}
	case queryTaskListMsg:
		m.taskList = msg
		m.tableModel = m.tableModel.WithRows(generateTaskRowsFromData(msg))
		if m.projectList == nil || m.allTags == nil {
			cmds = append(cmds, queryProjectList, queryAllTags)
		}
	case dbUpdateFinishedTasksMsg:
		cmds = append(cmds, queryTaskList)
	case dbUpdateFinishedProjectsMsg:
		cmds = append(cmds, queryProjectList)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) handleItemEdit(isFormInit bool, cmds []tea.Cmd) []tea.Cmd {
	selectedItemId := GetRowId(m.tableModel.HighlightedRow())
	if m.viewMode == "projects" {
		project := FindProjectById(selectedItemId, m.projectList)
		isFormInit = true
		m.editForm = BuildProjectEditForm(project)
		m.editFormState.isFormVisible = true
		cmds = append(cmds, m.editForm.Init())
	} else {
		cmds = append(cmds, queryTaskTags(selectedItemId))
	}
	return cmds
}

func (m *Model) switchView(mode string) {
	m.viewMode = mode
	m.tableModel = buildTableModel([]table.Row{}, mode)
	m.recalculateTable()
}

func (m *Model) PrepareCommandInput(activeCommand string, prompt string) {
	m.commandInputState.isActive = true
	m.commandInputState.activeCommand = activeCommand
	m.tableModel.WithKeyMap(table.KeyMap{})
	m.commandInput.Prompt = prompt
}

func handleFormCompletion(m *Model) tea.Cmd {
	return func() tea.Msg {
		rowId := GetRowId(m.tableModel.HighlightedRow())
		if m.viewMode == "projects" {
			return updateProject(rowId, m.editForm)
		}
		return updateTask(rowId, m.editForm)
	}
}

func handleUndelete(id int, viewMode string) tea.Cmd {
	return func() tea.Msg {
		if viewMode == "projects" {
			return undeleteProject(id)
		}
		return undeleteTask(id)
	}
}

func handleDeletion(m Model) tea.Cmd {
	return func() tea.Msg {
		rowId := GetRowId(m.tableModel.HighlightedRow())
		if m.viewMode == "projects" {
			return deleteProject(rowId)
		}
		return deleteTask(rowId)
	}
}

func (m *Model) handleQueryTagsResult(msg queryTaskTagsMsg) {
	m.editFormState.isFormVisible = true
	selectedTaskId := GetRowId(m.tableModel.HighlightedRow())
	selectedTask := m.GetTaskById(selectedTaskId)
	selectedTask.tags = Tags(msg)
	m.editForm = BuildTaskEditForm(selectedTask, m.allTags, m.projectList)
}

func (m *Model) recalculateTable() {
	m.tablePageSize = m.totalHeight - 7
	m.tableModel = m.tableModel.
		WithTargetWidth(m.totalWidth).
		WithMinimumHeight(m.totalHeight).
		WithPageSize(m.tablePageSize)
}

func (m Model) View() string {
	body := strings.Builder{}
	if m.editForm.State == huh.StateCompleted {
		m.editFormState.isFormVisible = false
	}
	if m.editFormState.isFormVisible {
		body.WriteString("\n")
		body.WriteString(m.editForm.View())
	} else {
		body.WriteString(m.tableModel.View())
	}
	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
