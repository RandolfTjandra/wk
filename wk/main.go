package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"

	"wk/pkg/db"
)

func main() {
	database, err := sql.Open("sqlite3", "wk.db")
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	sts := `
CREATE TABLE IF NOT EXISTS subjects(id INTEGER PRIMARY KEY, subject TEXT);
`
	_, err = database.Exec(sts)

	subjectRepo := db.NewSubjectRepo(database)

	model := initialModel(IndexView, subjectRepo)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting: %v", err)
		os.Exit(1)
	}
}
