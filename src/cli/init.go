package main

import (
	"database/sql"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var configPath string
var databaseName = "tasktango.db"
var dataBasePath string

func initTaskTango() tea.Cmd {
	return func() tea.Msg {
		initConfigFolder()
		initDatabase()
		ReadAppState()
		seed()
		return queryTaskList()
	}
}

func initConfigFolder() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configDirectory := os.Getenv("TASKTANGO_CONFIG")
	if configDirectory == "" {
		configPath = homeDir + "/.tasktango"
	} else {
		configPath = configDirectory
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.Mkdir(configPath, 0755)
		log.Println("Created config directory at " + configPath)
	}
}

func initDatabase() {
	dataBasePath = configPath + "/" + databaseName
	if _, err := os.Stat(dataBasePath); err == nil {
		return
	}
	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table Tasks (
		Id integer not null primary key,
		Description text,
		ProjectId integer,
		CreatedDate datetime,
		Priority integer,
		IsStarted INTEGER DEFAULT 0 CHECK (IsStarted IN (0, 1)),
		DeletedDate datetime,
		Context text
	);
	create table Projects (
		Id integer not null primary key,
		Title text,
		DeletedDate datetime,
		Context text
	);
	create table Tags (
		Id integer not null primary key,
		Text text
	);
	create table TaskTags (
		TaskId integer,
		TagId integer,
		primary key (TaskId, TagId)
	);
	insert into Tags (Text) values
		("Today"),
		("Blocked"),
		("Delegate"),
		("Maybe"),
		("Review");
	`
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
