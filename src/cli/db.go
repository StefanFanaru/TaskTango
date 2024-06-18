package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	_ "github.com/mattn/go-sqlite3"
)

type dbUpdateFinishedTasksMsg int
type dbUpdateFinishedProjectsMsg int
type queryTaskListMsg []Task
type queryProjectListMsg []ProjectRow
type queryTaskTagsMsg Tags
type queryAllTagsMsg Tags

type Tags map[int]string

type Task struct {
	id           int
	description  string
	projectTitle string
	projectId    string
	tagsJoined   string
	tags         Tags
	isStarted    bool
	priority     Priority
	createdDate  string
}

type ProjectRow struct {
	id        int
	title     string
	taskCount int
}

func startTask(taskId int) tea.Cmd {
	return func() tea.Msg {
		db, err := sql.Open("sqlite3", dataBasePath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// update task IsStarted to 1 if it's not started
		// if it's already started, set IsStarted to 0
		_, err = db.Exec(`
			update Tasks
			set IsStarted = case
				when IsStarted = 0 then 1
				else 0
			end
			where Id = ?
		`, taskId)

		if err != nil {
			log.Fatal(err)
		}

		return dbUpdateFinishedTasksMsg(0)
	}
}

func queryProjectList() tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	projectRows := []ProjectRow{}
	rows, err := db.Query(`
		select p.Id, p.Title, count(t.Id) as TaskCount
		from Projects p
		left join (select * from Tasks where DeletedDate IS NULL) t on p.Id = t.ProjectId
		where p.DeletedDate IS NULL
		and p.Context = ?
		group by p.Id
		order by p.Title asc
	`, AppState.ActiveContext)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var projectRow ProjectRow
		err := rows.Scan(&projectRow.id, &projectRow.title, &projectRow.taskCount)
		if err != nil {
			fmt.Println("Scan error:", err)
			return queryProjectListMsg(projectRows)
		}
		projectRows = append(projectRows, projectRow)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return queryProjectListMsg(projectRows)
}

func queryAllTags() tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select Id, Text from Tags order by Text asc")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	tags := Tags{}

	for rows.Next() {
		var id int
		var text string
		scanErr := rows.Scan(&id, &text)
		if scanErr != nil {
			fmt.Println("Scan error:", err)
			return queryTaskTagsMsg(tags)
		}
		tags[id] = text
	}

	return queryAllTagsMsg(tags)
}

func queryTaskTags(taskId int) tea.Cmd {
	return func() tea.Msg {
		db, err := sql.Open("sqlite3", dataBasePath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query(`
			select tag.Id, tag.Text
			from TaskTags tt
			join Tags tag on tt.TagId = tag.Id
			where tt.TaskId = ?
		`, taskId)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		taskTags := map[int]string{}

		for rows.Next() {
			var tagId int
			var tagText string
			scanErr := rows.Scan(&tagId, &tagText)
			if scanErr != nil {
				fmt.Println("Scan error:", err)
				return queryTaskTagsMsg(taskTags)
			}
			taskTags[tagId] = tagText
		}

		return queryTaskTagsMsg(taskTags)
	}
}

func updateProject(projectId int, form *huh.Form) tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	title := form.GetString("title")

	_, err = db.Exec("update Projects set Title = ? where Id = ?", title, projectId)
	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedProjectsMsg(0)
}

func updateTask(taskId int, form *huh.Form) tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	description := form.GetString("description")
	projectId := form.GetString("project")
	priority := form.GetString("priority")
	tags := form.Get("tags").([]string)

	_, err = db.Exec("update Tasks set Description = ?, ProjectId = ?, Priority = ? where Id = ?", description, projectId, priority, taskId)
	if err != nil {
		log.Fatal(err)
	}

	// delete all tags for this task
	_, err = db.Exec("delete from TaskTags where TaskId = ?", taskId)
	if err != nil {
		log.Fatal(err)
	}

	for _, tag := range tags {
		_, err = db.Exec("insert into TaskTags (TaskId, TagId) values (?, ?)", taskId, tag)
		if err != nil {
			log.Fatal(err)
		}
	}

	return dbUpdateFinishedTasksMsg(0)
}

func undeleteProject(id int) tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("update Projects set DeletedDate = NULL where Id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedProjectsMsg(0)
}

func undeleteTask(id int) tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("update Tasks set DeletedDate = NULL where Id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedTasksMsg(0)
}

func deleteProject(id int) tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deletedDate := time.Now().UTC().Format(time.RFC3339)
	_, err = db.Exec("update Projects set DeletedDate = ? where Id = ?", deletedDate, id)
	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedProjectsMsg(0)
}

func deleteTask(id int) tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deletedDate := time.Now().UTC().Format(time.RFC3339)
	_, err = db.Exec("update Tasks set DeletedDate = ? where Id = ?", deletedDate, id)
	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedTasksMsg(0)
}

func addProject(title string) tea.Msg {
	// trim whitespaces
	title = strings.TrimSpace(title)
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("insert into Projects (Title, Context) values (?, ?)", title, AppState.ActiveContext)
	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedProjectsMsg(0)
}

func changeTaskContext(taskId int, context string) tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("update Tasks set Context = ? where Id = ?", context, taskId)
	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedTasksMsg(0)
}

func addTask(description string) tea.Msg {
	// trim whitespaces
	description = strings.TrimSpace(description)
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	creteadDate := time.Now().UTC().Format(time.RFC3339)
	_, err = db.Exec("insert into Tasks (Description, CreatedDate, Priority, Context) values (?, ?, 1, ?)", description, creteadDate, AppState.ActiveContext)

	var taskId int
	err = db.QueryRow("select last_insert_rowid()").Scan(&taskId)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("insert into TaskTags (TaskId, TagId) values (?, 5)", taskId)

	if err != nil {
		log.Fatal(err)
	}

	return dbUpdateFinishedTasksMsg(0)
}

func queryTaskList() tea.Msg {
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	query := `
		select t.Id, t.Description, t.IsStarted, t.Priority, t.CreatedDate, p.Id || ') ' || p.Title as ProjectTitle, p.Id as ProjectId,
			   CASE
				   WHEN COUNT(tag.Text) > 0 THEN group_concat(tag.Text, ', ')
				   ELSE NULL
			   END as Tags
		from Tasks t
		left join Projects p on t.ProjectId = p.Id
		left join TaskTags tt on t.Id = tt.TaskId
		left join Tags tag on tt.TagId = tag.Id
		where t.DeletedDate IS NULL and p.DeletedDate IS NULL
		and t.Context = ?
		`

	if AppState.Filter.Tag != "" {
		query += " and tag.Text = '" + AppState.Filter.Tag + "'"
	}

	if AppState.Filter.ProjectId != "" {
		query += " and p.Id = '" + AppState.Filter.ProjectId + "'"
	}

	if AppState.Filter.Search != "" {
		query += " and t.Description like '%" + AppState.Filter.Search + "%'"
	}

	if AppState.Filter.Priority != "" {
		query += " and t.Priority = " + AppState.Filter.Priority
	}
	query = query + `
		group by t.Id
		order by t.Priority desc, t.IsStarted desc, t.CreatedDate asc
		`
	rows, err := db.Query(query, AppState.ActiveContext)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var task Task
		var tags sql.NullString
		var projectTitle sql.NullString
		var projectId sql.NullString
		var priority int

		err := rows.Scan(&task.id, &task.description, &task.isStarted, &priority, &task.createdDate, &projectTitle, &projectId, &tags)
		if err != nil {
			fmt.Println("Scan error:", err)
			return queryTaskListMsg(tasks)
		}

		if tags.Valid {
			task.tagsJoined = tags.String
		} else {
			task.tagsJoined = "-"
		}
		if projectTitle.Valid {
			task.projectTitle = projectTitle.String
		} else {
			task.projectTitle = "-"
		}
		if projectId.Valid {
			task.projectId = projectId.String
		}

		task.priority = Priority(priority)
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return queryTaskListMsg(tasks)
}
