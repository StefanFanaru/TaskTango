package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func seed() {
	// Seed the database with some initial data

	db, err := sql.Open("sqlite3", dataBasePath)
	if err != nil {
		log.Fatal(err)
	}

	// if table tasks contains data, do not seed
	rows, err := db.Query("select count(*) from Tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	if count > 0 {
		log.Println("Database already seeded")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// Insert some initial data
	_, err = tx.Exec(`
		insert into Projects (Title, Context) values
			("Project 1", ""),
			("Project 2", ""),
			("Project 3", ""),
			("Project 4", ""),
			("Project 5", "")
	`)
	if err != nil {
		log.Fatal(err)
	}

	// ("Task 8", 4, datetime('now'), 2, 0, ""),
	query := `
		insert into Tasks (Description, ProjectId, CreatedDate, Priority, IsStarted, Context) values
	`
	for i := 1; i <= 50; i++ {
		query += ` ("Task ` + fmt.Sprint(i) + `", 4, datetime('now'), 2, 0, "")`
		// Add a comma if this isn't the last iteration
		if i < 50 {
			query += ", "
		}
	}

	_, err = tx.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(`
		insert into TaskTags (TaskId, TagId) values
			(1, 1),
			(1, 2),
			(2, 1),
			(2, 2),
			(3, 3),
			(3, 4),
			(4, 3),
			(4, 4),
			(5, 2),
			(5, 1),
			(6, 2),
			(6, 1),
			(7, 2),
			(7, 3),
			(8, 2),
			(8, 3),
			(9, 4),
			(9, 5),
			(10, 4),
			(10, 5)
	`)

	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
