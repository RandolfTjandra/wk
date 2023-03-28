package ui

import (
	"encoding/json"
	"strings"

	"github.com/brandur/wanikaniapi"
)

// Renders a list of subjects
func RenderSubjects(subjects []*wanikaniapi.Subject) string {
	var b strings.Builder
	subjectCount := 0
	for _, subject := range subjects {
		if subject.KanjiData != nil {
			b.WriteString(Kanji(subject.KanjiData.Characters))
			if subjectCount < len(subjects)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else if subject.VocabularyData != nil {
			b.WriteString(Vocab(subject.VocabularyData.Characters))
			if subjectCount < len(subjects)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else if subject.RadicalData != nil && subject.RadicalData.Characters != nil {
			b.WriteString(Radical(*subject.RadicalData.Characters))
			if subjectCount < len(subjects)-1 {
				b.WriteString(", ")
			}
			subjectCount++
		} else {
			foo, _ := json.Marshal(subject)
			b.WriteString("\n\n" + string(foo) + "\n\n")
		}

		// hopefully can delete this if line wrap can work
		if subjectCount > 0 && subjectCount%10 == 0 {
			b.WriteString("\n  ")
		}
	}
	return b.String()
}
