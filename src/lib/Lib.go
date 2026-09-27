package lib

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var RESOURCES_PATH string
var CONFIG_PATH string

var Screen tcell.Screen
var DefaultStyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 255, 0))
var SelectedStyle = tcell.StyleDefault.Background(color.NewRGBColor(0, 255, 0)).Foreground(color.NewRGBColor(0, 0, 0))
var GameOverStyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 0, 0))

var (
	Height int
	Width  int
	Mute   = false
)

func DrawString(x, y int, text []string, style tcell.Style, s tcell.Screen) {
	for iR, r := range text {
		s.PutStrStyled(x, y+iR, r, style)
	}
}

func DrawChars(x, y int, text []string, styles [][]tcell.Style, s tcell.Screen) {
	for iR, r := range text {
		for iC, c := range r {
			if c != ' ' {
				s.Put(x+iC, y+iR, string(c), styles[iR][iC])
			}
		}
	}
}

func DrawWithFrame(x, y int, text []string, style tcell.Style, s tcell.Screen) {
	height := len(text)
	var width int
	for _, line := range text {
		length := len(line)
		if length > width {
			width = length
		}
	}

	s.Put(x, y, "┌", DefaultStyle)
	s.Put(x+width+1, y, "┐", DefaultStyle)
	s.Put(x, y+height+1, "└", DefaultStyle)
	s.Put(x+width+1, y+height+1, "┘", DefaultStyle)

	for i := 0; i < height+1; i++ {
		if i != 0 && i != height+1 {
			s.Put(x, y+i, "¦", DefaultStyle)
			s.Put(x+width+1, y+i, "¦", DefaultStyle)
		}

		for j := 0; j < width+1; j++ {
			if j != 0 && j != width+1 {
				s.Put(x+j, y, "-", DefaultStyle)
				s.Put(x+j, y+height+1, "-", DefaultStyle)
			}
		}
	}

	DrawString(x+1, y+1, text, style, s)
}

func AppendCero(text string, length int) string {
	for len(text) < length {
		text = "0" + text
	}
	return text
}
