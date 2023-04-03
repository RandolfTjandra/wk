package main

import (
	"fmt"
	"os"

	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"

	"wk/pkg/db"
	"wk/pkg/wanikani"
)

var (
	WKClient *wanikaniapi.Client
)

func main() {
	log.SetLevel(log.DebugLevel)
	database, err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	db.InitSubjectRepo(database)

	apiKey, ok := os.LookupEnv("WK_KEY")
	if !ok {
		fmt.Printf("no API key: %v", err)
		os.Exit(1)
	}
	wanikani.Init(apiKey)

	mainModel := initialMainModel(IndexView)
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting: %v", err)
		os.Exit(1)
	}
}
