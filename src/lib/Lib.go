package lib

import (
	"strings"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var RESOURCES_PATH string
var CONFIG_PATH string
var VAR_PATH string

var Screen tcell.Screen
var DefaultStyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 255, 0))
var SelectedStyle = tcell.StyleDefault.Background(color.NewRGBColor(0, 255, 0)).Foreground(color.NewRGBColor(0, 0, 0))
var GameOverStyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 0, 0))

var BackgroundPieceStyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 255, 255))
var BackgroundShadowStyles = [10]tcell.Style{
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 220, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 200, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 180, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 160, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 140, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 100, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 80, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 60, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 40, 0)),
	tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 20, 0)),
}

var (
	Height int
	Width  int
	Mute   = false
)

func DrawString(x, y int, text []string, style tcell.Style) {
	for iR, r := range text {
		Screen.PutStrStyled(x, y+iR, r, style)
	}
}

func DrawChars(x, y int, text []string, styles [][]tcell.Style) {
	for iR, r := range text {
		for iC, c := range r {
			if c != ' ' {
				Screen.Put(x+iC, y+iR, string(c), styles[iR][iC])
			}
		}
	}
}

func DrawWithFrame(x, y int, text []string, style tcell.Style) {
	height := len(text)
	var width int
	for _, line := range text {
		length := len(line)
		if length > width {
			width = length
		}
	}

	Screen.Put(x, y, "┌", DefaultStyle)
	Screen.Put(x+width+1, y, "┐", DefaultStyle)
	Screen.Put(x, y+height+1, "└", DefaultStyle)
	Screen.Put(x+width+1, y+height+1, "┘", DefaultStyle)

	for i := 0; i < height+1; i++ {
		if i != 0 && i != height+1 {
			Screen.Put(x, y+i, "¦", DefaultStyle)
			Screen.Put(x+width+1, y+i, "¦", DefaultStyle)
		}

		for j := 0; j < width+1; j++ {
			if j != 0 && j != width+1 {
				Screen.Put(x+j, y, "-", DefaultStyle)
				Screen.Put(x+j, y+height+1, "-", DefaultStyle)
			}
		}
	}

	DrawString(x+1, y+1, text, style)
}

func AppendBlank(text string, length int) string {
	if text == "0" {
		return strings.Repeat(" ", length)
	}
	for len(text) < length {
		text = " " + text
	}
	return text
}
