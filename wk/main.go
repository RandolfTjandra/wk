package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"

	"wk/pkg/db"
	"wk/pkg/wanikani"
)

func main() {
	database, err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	subjectRepo := db.NewSubjectRepo(database)

	wkClient := wanikaniapi.NewClient(&wanikaniapi.ClientConfig{
		APIToken: wanikani.ApiKey,
	})

	commander := NewCommander(true, subjectRepo, wkClient)

	model := initialModel(commander, IndexView, subjectRepo)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting: %v", err)
		os.Exit(1)
	}
}
