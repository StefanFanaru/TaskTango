package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/huh"
	"github.com/evertras/bubble-table/table"
)

type TextInputState struct {
	isActive      bool
	activeCommand string
	content       string
}

type UndoStack struct {
	recentlyDeleted map[int]string
}

type EditForm struct {
	description   string
	project       string
	tags          []string
	priority      int
	isFormVisible bool
}

type Model struct {
	tableModel        table.Model
	tablePageSize     int
	commandInput      textinput.Model
	commandInputState TextInputState
	totalWidth        int
	totalHeight       int
	undoStack         UndoStack
	editForm          *huh.Form
	editFormState     EditForm
	taskList          []Task
	projectList       []ProjectRow
	allTags           Tags
	viewMode          string
	filter            Filter
}
