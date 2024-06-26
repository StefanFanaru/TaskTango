package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

func BuildTextInput() textinput.Model {
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Prompt = ""
	return ti
}

func (m *Model) ExitTextInput() {
	m.commandInput.Reset()
	m.commandInput.Blur()
	m.commandInput.Prompt = ""
	m.commandInputState.activeCommand = ""
	m.commandInputState.isActive = false
	m.tableModel.WithKeyMap(table.DefaultKeyMap())
}

func doIfContextExists(context string, f func() tea.Msg) tea.Msg {
	if Contains(AppState.Contexts, context) {
		return f()
	}

	return tea.Msg(nil)
}

func (m *Model) HandleCommandInput(commandInput string, activeCommand string, viewMode string) tea.Cmd {
	return func() tea.Msg {
		type emptyMsg struct{}
		commandKeyWords := []string{"add", "ftag", "review", "fproj", "fprio", "nf", "c", "ac", "cc"}

		if commandInput == "" {
			return emptyMsg{}
		}

		if activeCommand == "" {
			activeCommand = strings.Split(commandInput, " ")[0]
		}

		if !Contains(commandKeyWords, activeCommand) {
			return emptyMsg{}
		}

		commandValue := strings.TrimPrefix(commandInput, activeCommand+" ")

		switch activeCommand {
		case "add":
			if viewMode == "projects" {
				return addProject(commandValue)
			}
			return addTask(commandValue)
		case "fprio":
			AppState.Filter.Priority = PriorityStingInt(commandValue)
			WriteAppState()
			return queryTaskList()
		case "ftag":
			AppState.Filter.Tag = strings.ToUpper(commandValue[:1]) + commandValue[1:]
			WriteAppState()
			return queryTaskList()
		case "fproj":
			AppState.Filter.ProjectId = commandValue
			WriteAppState()
			return queryTaskList()
		case "review":
			AppState.Filter.Tag = "Review"
			return queryTaskList()
		case "nf":
			AppState.Filter = Filter{}
			WriteAppState()
			return queryTaskList()
		case "cc":
			return doIfContextExists(commandValue, func() tea.Msg {
				rowId := GetRowId(m.tableModel.HighlightedRow())
				return changeTaskContext(rowId, commandValue)
			})
		case "ac":
			AddContext(commandValue)
		case "c":
			return doIfContextExists(commandValue, func() tea.Msg {
				SwitchToContext(commandValue)
				if viewMode == "projects" {
					return queryProjectList()
				} else {
					return queryTaskList()
				}
			})
		}
		return emptyMsg{}
	}
}
