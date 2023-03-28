package main

import (
	"fmt"
	"strings"

	"wk/pkg/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m mainModel) accountView() string {
	var b strings.Builder
	activeStatus := ui.BatsuMark
	if m.User.Data.Subscription.Active {
		activeStatus = ui.CheckMark
	}
	userDataContent := fmt.Sprintf(
		"Subscription\n"+
			"Active: %s\n"+
			"Type:   %s",
		activeStatus,
		m.User.Data.Subscription.Type,
	)
	userData := lipgloss.NewStyle().Render(userDataContent)
	b.WriteString(userData)
	return b.String()
}
