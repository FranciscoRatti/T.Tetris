package obj

import (
	"TTetris/src/lib"
	"math/rand"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

var ShadowStyle = tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(120, 120, 120))

type Piece struct {
	x, y, r int
	id      byte
	style   tcell.Style
	sprite  [][]string
}

var (
	Pieces = []Piece{
		{0, 0, 0, 0,
			tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 122, 0)),
			[][]string{
				{
					"[]",
					"[][][]",
				}, {
					"[][]",
					"[]",
					"[]",
				}, {
					"[][][]",
					"    []",
				}, {
					"  []",
					"  []",
					"[][]",
				},
			}},
		{0, 0, 0, 1,
			tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 0, 255)),
			[][]string{
				{
					"    []",
					"[][][]",
				}, {
					"[]",
					"[]",
					"[][]",
				}, {
					"[][][]",
					"[]",
				}, {
					"[][]",
					"  []",
					"  []",
				},
			}},
		{0, 0, 0, 2,
			tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 0, 0)),
			[][]string{
				{
					"[][]",
					"  [][]",
				}, {
					"  []",
					"[][]",
					"[]",
				},
			}},
		{0, 0, 0, 3,
			tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 255, 0)),
			[][]string{
				{
					"  [][]",
					"[][]",
				}, {
					"[]",
					"[][]",
					"  []",
				},
			}},
		{0, 0, 0, 4,
			tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 255, 0)),
			[][]string{
				{
					"[][]",
					"[][]",
				},
			}},
		{0, 0, 0, 5,
			tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(255, 0, 255)),
			[][]string{
				{
					"  []",
					"[][][]",
				}, {
					"[]",
					"[][]",
					"[]",
				}, {
					"[][][]",
					"  []",
				}, {
					"  []",
					"[][]",
					"  []",
				},
			}},
		{0, 0, 0, 6,
			tcell.StyleDefault.Background(color.Reset).Foreground(color.NewRGBColor(0, 255, 255)),
			[][]string{
				{
					"[][][][]",
				}, {
					"[]",
					"[]",
					"[]",
					"[]",
				},
			}},
	}
)

// Getters, Setters y Constructor --------------------------------------------------------------------------------------

func NewPiece(nextPieces [3]*Piece) *Piece {

	id := byte(rand.Intn(7))
	for nextPieces[0].id == id || nextPieces[1].id == id || nextPieces[2].id == id {
		id = byte(rand.Intn(7))
	}

	newPiece := &Pieces[id]
	return newPiece
}

func (p *Piece) GetId() byte {
	return p.id
}

func (p *Piece) GetActualSprite() []string {
	return p.sprite[p.r]
}

func (p *Piece) GetSprite(r int) []string {
	switch p.id {
	case 0, 1, 5:
		return p.sprite[r]
	case 2, 3, 6:
		return p.sprite[r/2]
	default:
		return p.sprite[0]
	}

}

func (p *Piece) GetStyle() tcell.Style {
	return p.style
}

func (p *Piece) SetPos(x, y int) {
	p.x = x
	p.y = y
}

func (p *Piece) GetY() int {
	return p.y
}
func (p *Piece) GetX() int {
	return p.x
}

func (p *Piece) SetRotation(r int) {
	if p.id != 4 {
		p.r = r
	}
}

func IsEmpty(stationaryPieces []string, pos [][]int) bool {
	for i := range pos {
		if pos[i][0] >= 20 || pos[i][1] >= 20 || pos[i][0] < 0 {
			return false
		}

		if pos[i][1] >= 0 {
			c := stationaryPieces[pos[i][1]][pos[i][0]]
			if c != ' ' {
				return false
			}
		}
	}
	return true
}

// Draw ----------------------------------------------------------------------------------------------------------------

func (p *Piece) Draw(stationaryPieces []string) {
	width, height := lib.Screen.Size()

	// Sombra
	if Config.Game.Shadow {
		y := p.GetFloorRow(stationaryPieces)

		for iR, r := range p.sprite[p.r] {
			for iC, c := range r {
				if c != ' ' {
					lib.Screen.Put(width/2-10+p.x+iC, height/2-11+y+iR, string(c), ShadowStyle)
				}
			}
		}
	}

	// Pieza
	for iR, r := range p.sprite[p.r] {
		for iC, c := range r {
			if c != ' ' {
				lib.Screen.Put(width/2-10+p.x+iC, height/2-11+p.y+iR, string(c), p.style)
			}
		}
	}
}

// Movimiento ----------------------------------------------------------------------------------------------------------

func (p *Piece) MoveRight(stationaryPieces []string) {
	switch p.id {
	case 0:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x + 6, p.y + 1},
			}) {
				p.x += 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 2, p.y + 1},
				{p.x + 2, p.y + 2},
			}) {
				p.x += 2
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 6, p.y},
				{p.x + 6, p.y + 1},
			}) {
				p.x += 2
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 4, p.y + 1},
				{p.x + 4, p.y + 2},
			}) {
				p.x += 2
			}
		}
	case 1:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 6, p.y},
				{p.x + 6, p.y + 1},
			}) {
				p.x += 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x + 2, p.y + 1},
				{p.x + 4, p.y + 2},
			}) {
				p.x += 2
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 6, p.y},
				{p.x + 2, p.y + 1},
			}) {
				p.x += 2
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 4, p.y + 1},
				{p.x + 4, p.y + 2},
			}) {
				p.x += 2
			}
		}
	case 2:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 6, p.y + 1},
			}) {
				p.x += 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 4, p.y + 1},
				{p.x + 2, p.y + 2},
			}) {
				p.x += 2
			}
		}
	case 3:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 6, p.y},
				{p.x + 4, p.y + 1},
			}) {
				p.x += 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x + 4, p.y + 1},
				{p.x + 4, p.y + 2},
			}) {
				p.x += 2
			}
		}
	case 4:
		if IsEmpty(stationaryPieces, [][]int{
			{p.x + 4, p.y},
			{p.x + 4, p.y + 1},
		}) {
			p.x += 2
		}
	case 5:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 6, p.y + 1},
			}) {
				p.x += 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x + 4, p.y + 1},
				{p.x + 2, p.y + 2},
			}) {
				p.x += 2
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 6, p.y},
				{p.x + 4, p.y + 1},
			}) {
				p.x += 2
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 4, p.y + 1},
				{p.x + 4, p.y + 2},
			}) {
				p.x += 2
			}
		}
	case 6:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 8, p.y},
			}) {
				p.x += 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x + 2, p.y + 1},
				{p.x + 2, p.y + 2},
				{p.x + 2, p.y + 3},
			}) {
				p.x += 2
			}
		}
	}
}

func (p *Piece) MoveLeft(stationaryPieces []string) {
	switch p.id {
	case 0:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
			}) {
				p.x -= 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
				{p.x - 2, p.y + 2},
			}) {
				p.x -= 2
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x + 2, p.y + 1},
			}) {
				p.x -= 2
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y},
				{p.x, p.y + 1},
				{p.x - 2, p.y + 2},
			}) {
				p.x -= 2
			}
		}
	case 1:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x - 2, p.y + 1},
			}) {
				p.x -= 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
				{p.x - 2, p.y + 2},
			}) {
				p.x -= 2
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
			}) {
				p.x -= 2
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x, p.y + 1},
				{p.x, p.y + 2},
			}) {
				p.x -= 2
			}
		}
	case 2:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x, p.y + 1},
			}) {
				p.x -= 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y},
				{p.x - 2, p.y + 1},
				{p.x - 2, p.y + 2},
			}) {
				p.x -= 2
			}
		}
	case 3:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y},
				{p.x - 2, p.y + 1},
			}) {
				p.x -= 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
				{p.x, p.y + 2},
			}) {
				p.x -= 2
			}
		}
	case 4:
		if IsEmpty(stationaryPieces, [][]int{
			{p.x - 2, p.y},
			{p.x - 2, p.y + 1},
		}) {
			p.x -= 2
		}
	case 5:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y},
				{p.x - 2, p.y + 1},
			}) {
				p.x -= 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
				{p.x - 2, p.y + 2},
			}) {
				p.x -= 2
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x, p.y + 1},
			}) {
				p.x -= 2
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y},
				{p.x - 2, p.y + 1},
				{p.x, p.y + 2},
			}) {
				p.x -= 2
			}
		}
	case 6:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
			}) {
				p.x -= 2
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
				{p.x - 2, p.y + 2},
				{p.x - 2, p.y + 3},
			}) {
				p.x -= 2
			}
		}
	}
}

func (p *Piece) MoveDown(stationaryPieces []string) bool {
	switch p.id {
	case 0:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 2},
				{p.x + 2, p.y + 2},
				{p.x + 4, p.y + 2},
			}) {
				p.y++
				return true
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 3},
				{p.x + 2, p.y + 1},
			}) {
				p.y++
				return true
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 1},
				{p.x + 2, p.y + 1},
				{p.x + 4, p.y + 2},
			}) {
				p.y++
				return true
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 3},
				{p.x + 2, p.y + 3},
			}) {
				p.y++
				return true
			}
		}
	case 1:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 2},
				{p.x + 2, p.y + 2},
				{p.x + 4, p.y + 2},
			}) {
				p.y++
				return true
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 3},
				{p.x + 2, p.y + 3},
			}) {
				p.y++
				return true
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 2},
				{p.x + 2, p.y + 1},
				{p.x + 4, p.y + 1},
			}) {
				p.y++
				return true
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 1},
				{p.x + 2, p.y + 3},
			}) {
				p.y++
				return true
			}
		}
	case 2:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 1},
				{p.x + 2, p.y + 2},
				{p.x + 4, p.y + 2},
			}) {
				p.y++
				return true
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 3},
				{p.x + 2, p.y + 2},
			}) {
				p.y++
				return true
			}
		}
	case 3:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 2},
				{p.x + 2, p.y + 2},
				{p.x + 4, p.y + 1},
			}) {
				p.y++
				return true
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 2},
				{p.x + 2, p.y + 3},
			}) {
				p.y++
				return true
			}
		}
	case 4:
		if IsEmpty(stationaryPieces, [][]int{
			{p.x, p.y + 2},
			{p.x + 2, p.y + 2},
		}) {
			p.y++
			return true
		}
	case 5:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 2},
				{p.x + 2, p.y + 2},
				{p.x + 4, p.y + 2},
			}) {
				p.y++
				return true
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 3},
				{p.x + 2, p.y + 2},
			}) {
				p.y++
				return true
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 1},
				{p.x + 2, p.y + 2},
				{p.x + 4, p.y + 1},
			}) {
				p.y++
				return true
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 2},
				{p.x + 2, p.y + 3},
			}) {
				p.y++
				return true
			}
		}
	case 6:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 1},
				{p.x + 2, p.y + 1},
				{p.x + 4, p.y + 1},
				{p.x + 6, p.y + 1},
			}) {
				p.y++
				return true
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 4},
			}) {
				p.y++
				return true
			}
		}
	}
	return false
}

func (p *Piece) MoveFloor(stationaryPieces []string) {
	p.y = p.GetFloorRow(stationaryPieces)
}

func (p *Piece) GetFloorRow(stationaryPieces []string) int {
	y := p.y
	switch p.id {
	case 0:
		switch p.r {
		case 0:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 2},
				{p.x + 2, y + 2},
				{p.x + 4, y + 2},
			}) {
				y++
			}
		case 1:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 3},
				{p.x + 2, y + 1},
			}) {
				y++
			}
		case 2:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 1},
				{p.x + 2, y + 1},
				{p.x + 4, y + 2},
			}) {
				y++
			}
		case 3:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 3},
				{p.x + 2, y + 3},
			}) {
				y++
			}
		}
	case 1:
		switch p.r {
		case 0:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 2},
				{p.x + 2, y + 2},
				{p.x + 4, y + 2},
			}) {
				y++
			}
		case 1:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 3},
				{p.x + 2, y + 3},
			}) {
				y++
			}
		case 2:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 2},
				{p.x + 2, y + 1},
				{p.x + 4, y + 1},
			}) {
				y++
			}
		case 3:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 1},
				{p.x + 2, y + 3},
			}) {
				y++
			}
		}
	case 2:
		switch p.r {
		case 0:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 1},
				{p.x + 2, y + 2},
				{p.x + 4, y + 2},
			}) {
				y++
			}
		case 1:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 3},
				{p.x + 2, y + 2},
			}) {
				y++
			}
		}
	case 3:
		switch p.r {
		case 0:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 2},
				{p.x + 2, y + 2},
				{p.x + 4, y + 1},
			}) {
				y++
			}
		case 1:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 2},
				{p.x + 2, y + 3},
			}) {
				y++
			}
		}
	case 4:
		for IsEmpty(stationaryPieces, [][]int{
			{p.x, y + 2},
			{p.x + 2, y + 2},
		}) {
			y++
		}
	case 5:
		switch p.r {
		case 0:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 2},
				{p.x + 2, y + 2},
				{p.x + 4, y + 2},
			}) {
				y++
			}
		case 1:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 3},
				{p.x + 2, y + 2},
			}) {
				y++
			}
		case 2:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 1},
				{p.x + 2, y + 2},
				{p.x + 4, y + 1},
			}) {
				y++
			}
		case 3:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 2},
				{p.x + 2, y + 3},
			}) {
				y++
			}
		}
	case 6:
		switch p.r {
		case 0:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 1},
				{p.x + 2, y + 1},
				{p.x + 4, y + 1},
				{p.x + 6, y + 1},
			}) {
				y++
			}
		case 1:
			for IsEmpty(stationaryPieces, [][]int{
				{p.x, y + 4},
			}) {
				y++
			}
		}
	}

	return y
}

func (p *Piece) RotateRight(stationaryPieces []string) {
	switch p.id {
	case 0:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x + 4, p.y},
				{p.x + 2, p.y + 2},
			}) {
				p.x += 2
				p.r++
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
			}) {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 2, p.y + 1},
				}) {
					p.x -= 2
					p.r++
				}
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 4, p.y},
					{p.x + 4, p.y + 1},
				}) {
					p.r++
				}
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y + 1},
				{p.x + 2, p.y + 2},
				{p.x, p.y + 2},
			}) {
				p.r++
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y + 1},
			}) {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x, p.y},
					{p.x, p.y + 1},
				}) {
					p.r = 0
				}
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x - 2, p.y},
					{p.x - 2, p.y + 1},
					{p.x, p.y + 1},
				}) {
					p.x -= 2
					p.r = 0
				}
			}
		}
	case 1:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y},
				{p.x + 2, p.y + 2},
				{p.x + 4, p.y + 2},
			}) {
				p.x += 2
				p.r++
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y},
				{p.x - 2, p.y + 1},
			}) {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 2, p.y},
				}) {
					p.x -= 2
					p.r++
				}
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 2, p.y},
					{p.x + 4, p.y},
				}) {
					p.r++
				}
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y + 1},
				{p.x + 2, p.y + 2},
			}) {
				p.r++
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y},
				{p.x + 4, p.y + 1},
			}) {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x, p.y + 1},
				}) {
					p.r = 0
				}
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x, p.y + 1},
					{p.x - 2, p.y + 1},
				}) {
					p.x -= 2
					p.r = 0
				}
			}
		}
	case 2:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x, p.y + 1},
				{p.x, p.y + 2},
			}) {
				p.r++
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y + 1},
			}) {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x, p.y},
				}) {
					p.r = 0
				}
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x, p.y},
					{p.x - 2, p.y},
				}) {
					p.x -= 2
					p.r = 0
				}
			}
		}
	case 3:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y + 1},
				{p.x + 4, p.y + 2},
			}) {
				p.x += 2
				p.r++
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y + 1},
			}) {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 2, p.y},
				}) {
					p.x -= 2
					p.r = 0
				}
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 2, p.y},
					{p.x + 4, p.y},
				}) {
					p.r = 0
				}
			}
		}
	case 5:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y + 2},
			}) {
				p.x += 2
				p.r++
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y + 1},
			}) {
				p.x -= 2
				p.y++
				p.r++
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 4, p.y + 1},
					{p.x + 2, p.y + 2},
				}) {
					p.y++
					p.r++
				}
			}
		case 2:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y - 1},
			}) {
				p.y--
				p.r++
			}
		case 3:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 4, p.y + 1},
			}) {
				p.r = 0
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x, p.y},
					{p.x - 2, p.y + 1},
				}) {
					p.x -= 2
					p.r = 0
				}
			}
		}
	case 6:
		switch p.r {
		case 0:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x + 2, p.y - 1},
				{p.x + 2, p.y + 1},
				{p.x + 2, p.y + 2},
			}) {
				p.x += 2
				p.y--
				p.r++
			}
		case 1:
			if IsEmpty(stationaryPieces, [][]int{
				{p.x - 2, p.y + 1},
			}) {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 4, p.y + 1},
					{p.x + 2, p.y + 1},
				}) {
					p.x -= 2
					p.y++
					p.r = 0
				} else {
					if IsEmpty(stationaryPieces, [][]int{
						{p.x - 4, p.y + 1},
					}) {
						if IsEmpty(stationaryPieces, [][]int{
							{p.x + 2, p.y + 1},
						}) {
							p.x -= 4
							p.y++
							p.r = 0
						} else {
							if IsEmpty(stationaryPieces, [][]int{
								{p.x - 6, p.y + 1},
							}) {
								p.x -= 6
								p.y++
								p.r = 0
							}
						}
					}
				}
			} else {
				if IsEmpty(stationaryPieces, [][]int{
					{p.x + 2, p.y + 1},
					{p.x + 4, p.y + 1},
					{p.x + 6, p.y + 1},
				}) {
					p.y++
					p.r = 0
				}
			}
		}
	}
}
