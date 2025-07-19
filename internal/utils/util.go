package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const RED_BG_WHITE_FG = "\033[37;41m"
const RESETCOLOUR = "\033[0m"

const BAR = "█"
const VER = "11.9.2 alpha"

func BarColor(char string, count int, color string) string {
	if count < 0 {
		count = 0
	}
	if count > 10 {
		count = 10
	}
	bars := strings.Repeat(char, count)
	if count == 0 {
		bars = " "
	}
	coloredBars := fmt.Sprintf("[%s]%s[-]", color, bars)
	return coloredBars + strings.Repeat(" ", 10-count)
}

func Between(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Split(text, " ")
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine) == 0 {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

func TruncateText(text string, maxWidth int) string {
	if maxWidth <= 3 {
		return "..."[:maxWidth]
	}
	if len(text) > maxWidth {
		return text[:maxWidth-3] + "..."
	}
	return text
}

func SafePrintText(screen tcell.Screen, text string, x, y, w, h int, color tcell.Color) int {
	if w <= 2 || h <= 1 {
		return y
	}

	lines := WrapText(text, w-2)
	currentY := y

	for _, line := range lines {
		if currentY >= y+h-1 {
			break
		}
		tview.Print(screen, line, x+1, currentY, w-2, tview.AlignLeft, color)
		currentY++
	}

	return currentY
}

func FormatTime(timeInt int64) string {
	t := time.Unix(int64(timeInt), 0)
	return t.Format("15:04:05")
}

func formatMap(m map[string]interface{}, indent string, b *strings.Builder) {
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			b.WriteString(fmt.Sprintf("%s%s:\n", indent, k))
			formatMap(val, indent+"  ", b)
		case []interface{}:
			b.WriteString(fmt.Sprintf("%s%s:\n", indent, k))
			formatSlice(val, indent+"  ", b)
		default:
			b.WriteString(fmt.Sprintf("%s%s: %v\n", indent, k, val))
		}
	}
}

func formatSlice(s []interface{}, indent string, b *strings.Builder) {
	for _, v := range s {
		switch val := v.(type) {
		case map[string]interface{}:
			formatMap(val, indent+"  ", b)
		case []interface{}:
			formatSlice(val, indent+"  ", b)
		default:
			b.WriteString(fmt.Sprintf("%s- %v\n", indent, val))
		}
	}
}

func FormatJSONToString(rawJSON string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(rawJSON), &obj); err != nil {
		return "Cannot retrieve info!"
	}

	var b strings.Builder
	formatMap(obj, "", &b)
	return b.String()
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}

func RepeatString(str string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(str, count)
}
