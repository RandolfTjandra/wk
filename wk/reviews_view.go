package main

import (
	"fmt"
	"strings"
)

// Render reviews view
func (m mainModel) reviewsView() string {
	if m.Reviews == nil {
		return "loading..."
	}
	var b strings.Builder
	b.WriteString("  Reviews:\n")
	for i, review := range m.Reviews.Data {
		b.WriteString(fmt.Sprintf("%d: %d\n", i, review.Data.SubjectID))
	}

	return b.String() + "\n"
}
