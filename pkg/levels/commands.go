package levels

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"wk/pkg/wanikani"
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
