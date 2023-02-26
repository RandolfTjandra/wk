package main

import (
	"github.com/brandur/wanikaniapi"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	navChoices  []PageView
	currentPage PageView
	cursor      int
	response    []byte
	err         error

	Summary *wanikaniapi.Summary
	User    *wanikaniapi.User
}

func (m model) Init() tea.Cmd {
	return getUser
}

func initialModel(view PageView) model {
	return model{
		currentPage: view,
		navChoices: []PageView{
			SummaryView,
			LessonsView,
			AssignmentsView,
			ReviewsView,
			SettingsView,
		},
		cursor: 0,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *wanikaniapi.User:
		m.User = msg
		return m, nil
	case *wanikaniapi.Summary:
		m.Summary = msg
		return m, nil
	case errMsg:
		m.err = msg
		return m, tea.Quit
	// key press
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

func (m model) handleKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.currentPage = IndexView
		return m, nil
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.navChoices)-1 {
			m.cursor++
		}
	case "enter":
		m.currentPage = m.navChoices[m.cursor]
		switch m.currentPage {
		case SummaryView:
			return m, getSummary
		}
	}
	return m, nil
}
