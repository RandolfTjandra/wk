package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"

	"wk/pkg/db"
	"wk/pkg/summary"
)

func main() {
	database, err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	subjectRepo := db.NewSubjectRepo(database)

	apiKey, ok := os.LookupEnv("WK_KEY")
	if !ok {
		fmt.Printf("no API key: %v", err)
		os.Exit(1)
	}

	wkClient := wanikaniapi.NewClient(&wanikaniapi.ClientConfig{
		APIToken: apiKey,
	})

	commander := NewCommander(true, subjectRepo, wkClient)
	summaryCommander := summary.NewCommander(true, subjectRepo, wkClient)

	mainModel := initialMainModel(commander, summaryCommander, IndexView, subjectRepo)
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting: %v", err)
		os.Exit(1)
	}
}
