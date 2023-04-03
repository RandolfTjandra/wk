package account

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"wk/pkg/wanikani"
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// returns *wanikaniapi.User
func GetUser() tea.Msg {
	user, err := wanikani.GetUser(context.Background())
	if err != nil {
		return errMsg{err}
	}
	return user
}
