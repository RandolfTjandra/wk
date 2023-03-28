package main

import (
	"fmt"
	"strings"
)

func (m mainModel) levelsView() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%#v\n", m.Levels))
	return b.String()
}
