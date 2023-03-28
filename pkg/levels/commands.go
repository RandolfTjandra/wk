package levels

import (
	"context"
	"wk/pkg/wanikani"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// returns []*wanikaniapi.GetLevelProgression
func GetLevelProgressions() tea.Msg {
	progressions, err := wanikani.GetLevelProgressions(context.Background(), wanikani.Client)
	if err != nil {
		return errMsg{err}
	}

	return progressions
}
