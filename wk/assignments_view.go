package main

import (
	"fmt"
	"strings"
)

// Render assignments view
func (m mainModel) assignmentsView() string {
	if m.Assignments == nil {
		return fmt.Sprintf("loading...")
	}
	var b strings.Builder
	b.WriteString("  Assignments:\n")
	tally := map[string]int{}
	for _, assignment := range m.Assignments {
		var classification string
		switch stage := assignment.Data.SRSStage; {
		case stage < 5:
			classification = "apprentice"
		case stage < 7:
			classification = "guru"
		case stage < 8:
			classification = "master"
		case stage < 9:
			classification = "enlightened"
		}
		tally[classification] += 1
	}
	b.WriteString(fmt.Sprintf("Apprentice:  %d\n", tally["apprentice"]))
	b.WriteString(fmt.Sprintf("Guru:        %d\n", tally["guru"]))
	b.WriteString(fmt.Sprintf("Master:      %d\n", tally["master"]))
	b.WriteString(fmt.Sprintf("Enlightened: %d\n", tally["enlightened"]))

	return b.String() + "\n"
}
