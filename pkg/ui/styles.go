package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

const (
	progressFullChar      = "█"
	progressEmptyChar     = "░"
	progressBarLeftColor  = "#B14FFF"
	progressBarRightColor = "#00FFA3"

	UIWidth   = 80
	UIXMargin = 2
	UIYMargin = 1
)

// styling
var (
	term = termenv.EnvColorProfile()

	Keyword       = makeFgStyle("211")
	Subtle        = makeFgStyle("241")
	Dot           = colorFg(" • ", "236")
	progressEmpty = Subtle(progressEmptyChar)

	Radical = makeFgBgStyle("#E5E5E5", "#00AAFF")
	Kanji   = makeFgBgStyle("#E5E5E5", "#FF00AA")
	Vocab   = makeFgBgStyle("#E5E5E5", "#9400FF")

	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warn      = lipgloss.AdaptiveColor{Light: "#BD3762", Dark: "#E04376"}

	Debug = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#DD0044"))

	H1Title = lipgloss.NewStyle().
		Width(UIWidth-UIXMargin).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).Margin(UIYMargin, UIXMargin)

	CheckMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	BatsuMark = lipgloss.NewStyle().SetString("x").
			Foreground(warn).
			PaddingRight(1).
			String()
)

func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

func makeFgBgStyle(fgColor string, bgColor string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(fgColor)).Background(term.Color(bgColor)).Styled
}

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

func Progressbar(width int, numerator float64, denominator float64) string {
	w := float64(width)
	percent := numerator / denominator

	fullSize := int(math.Round(w * percent))
	var fullCells string
	ramp := makeRamp(progressBarLeftColor, progressBarRightColor, w)
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(progressFullChar).Foreground(term.Color(ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(progressEmpty, emptySize)

	return fmt.Sprintf("%s%s %.0f/%.0f", fullCells, emptyCells, numerator, denominator)
}
