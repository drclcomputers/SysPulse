package utils

import (
	"github.com/gdamore/tcell/v2"
)

var ColorMap = map[string]tcell.Color{
	"black":   tcell.ColorBlack,
	"red":     tcell.ColorRed,
	"green":   tcell.ColorGreen,
	"yellow":  tcell.ColorYellow,
	"blue":    tcell.ColorBlue,
	"magenta": tcell.ColorFuchsia,
	"cyan":    tcell.ColorAqua,
	"white":   tcell.ColorWhite,
	"orange":  tcell.ColorOrange,
	"purple":  tcell.ColorPurple,
	"pink":    tcell.ColorPink,
	"lime":    tcell.ColorLime,
	"teal":    tcell.ColorTeal,
	"aqua":    tcell.ColorAqua,
	"navy":    tcell.ColorNavy,
	"gray":    tcell.ColorGray,
	"silver":  tcell.ColorSilver,
	"maroon":  tcell.ColorMaroon,
	"olive":   tcell.ColorOlive,
	"default": tcell.ColorDefault,
}

func GetColorFromName(colorName string) tcell.Color {
	if color, exists := ColorMap[colorName]; exists {
		return color
	}
	return tcell.ColorDefault
}
