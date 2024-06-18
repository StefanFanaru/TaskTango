package main

func (m *Model) pushDeletionStack() {
	row := m.tableModel.HighlightedRow()
	rowValue := row.Data[columnKeyID]
	rowIdInt := rowValue.(int)
	m.undoStack.recentlyDeleted[rowIdInt] = m.viewMode
}

func (m *Model) popDeletionStack() int {
	lastDeleted := 0
	for key := range m.undoStack.recentlyDeleted {
		lastDeleted = key
	}
	delete(m.undoStack.recentlyDeleted, lastDeleted)
	return lastDeleted
}
