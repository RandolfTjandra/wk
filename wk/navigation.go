package main

import (
	"wk/pkg/levels"

	tea "github.com/charmbracelet/bubbletea"
)

func (m mainModel) handleIndexKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.currentPage = IndexView
		return m, nil
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursors[m.currentPage] > 0 {
			m.cursors[m.currentPage]--
		}
	case "down", "j":
		switch m.currentPage {
		case IndexView:
			if m.cursors[m.currentPage] < len(m.navChoices)-1 {
				m.cursors[m.currentPage]++
			}
		}
	case "enter":
		m.currentPage = m.navChoices[m.cursors[IndexView]]
		switch m.currentPage {
		case LevelsView:
			return m, levels.GetLevelProgressions
		}
	}
	return m, nil
}

func (m mainModel) handleDefaultKeyPress(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.currentPage = IndexView
		return m, nil
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursors[m.currentPage] > 0 {
			m.cursors[m.currentPage]--
		}
	case "down", "j":
		switch m.currentPage {
		case IndexView:
			if m.cursors[m.currentPage] < len(m.navChoices)-1 {
				m.cursors[m.currentPage]++
			}
		}
	}
	return m, nil
}
