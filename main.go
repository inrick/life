package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	const (
		ScreenWidth  = 640
		ScreenHeight = 480
		CellSize     = 10
		Width        = ScreenWidth / CellSize
		Height       = ScreenHeight / CellSize
	)
	// NrFramesPerStep can be changed during runtime.
	NrFramesPerStep := int32(10)
	rl.InitWindow(ScreenWidth, ScreenHeight, "Life")
	rl.SetTargetFPS(60)
	board := NewRect(Width, Height)
	board.Init()
	boardCopy := NewRect(Width, Height)
	shouldClose := false
	state := StateRunning
	var frameNr uint
	var mousePos rl.Vector2
	var selected Rect
	for !shouldClose {
		// State handling
		shouldClose = rl.WindowShouldClose() || rl.IsKeyPressed(rl.KeyQ)
		mousePos = rl.GetMousePosition()
		mx, my := int32(mousePos.X/CellSize), int32(mousePos.Y/CellSize)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !(0 <= mx && mx <= board.width && 0 <= my && my <= board.height) {
				// Should never trigger
				panic([]int32{mx, my})
			}
			if len(selected.buf) == 0 {
				*board.AtPtr(my, mx) = (board.At(my, mx) + 1) % 2
				state = StatePaused
			} else {
				board.Add(selected, mx, my)
				selected = Rect{}
				state = StateRunning
			}
		}
		switch {
		case rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeySpace):
			state = (state + 1) % 2
		case rl.IsKeyPressed(rl.KeyMinus):
			NrFramesPerStep = Min(NrFramesPerStep+1, 60)
		case rl.IsKeyPressed(rl.KeyEqual):
			NrFramesPerStep = Max(NrFramesPerStep-1, 1)
		case rl.IsKeyPressed(rl.KeyZero):
			board.Zero()
		case rl.IsKeyPressed(rl.KeyOne):
			selected = Creatures[GliderRight]
			state = StatePaused
		case rl.IsKeyPressed(rl.KeyTwo):
			selected = Creatures[GliderLeft]
			state = StatePaused
		case rl.IsKeyPressed(rl.KeyThree):
			selected = Creatures[Spaceship]
			state = StatePaused
		case rl.IsKeyPressed(rl.KeyFour):
			selected = Creatures[PrePulsar]
			state = StatePaused
		case rl.IsKeyPressed(rl.KeyFive):
			selected = Creatures[Pulsar]
			state = StatePaused
		case rl.IsKeyPressed(rl.KeySix):
			selected = Creatures[QueenBee]
			state = StatePaused
		}

		// Drawing
		rl.BeginDrawing()
		rl.ClearBackground(Colors[0])

		if state == StateRunning {
			if frameNr%uint(NrFramesPerStep) == 0 {
				for y := int32(0); y < board.height; y++ {
					for x := int32(0); x < board.width; x++ {
						alive := board.At(x, y)
						n := board.Neighbors(x, y)
						if (alive == 1 && (n == 2 || n == 3)) || (alive == 0 && n == 3) {
							*boardCopy.AtPtr(x, y) = 1
						} else {
							*boardCopy.AtPtr(x, y) = 0
						}
					}
				}
				copy(board.buf, boardCopy.buf)
			}
			frameNr++
		}

		for y := int32(0); y < board.height; y++ {
			for x := int32(0); x < board.width; x++ {
				cellColor := Colors[board.At(x, y)]
				x0 := x * CellSize
				y0 := y * CellSize
				rl.DrawRectangle(int32(x0), int32(y0), CellSize, CellSize, cellColor)
			}
		}

		if len(selected.buf) > 0 {
			for y := int32(0); y < selected.height; y++ {
				for x := int32(0); x < selected.width; x++ {
					cellColor := Colors[selected.At(x, y)]
					x0 := ((x + mx) % board.width) * CellSize
					y0 := ((y + my) % board.height) * CellSize
					rl.DrawRectangle(int32(x0), int32(y0), CellSize, CellSize, cellColor)
				}
			}
		}

		if state == StatePaused {
			rl.DrawText(
				"Press Space or P to continue running",
				ScreenWidth/2-190,
				ScreenHeight/2-10,
				20,
				rl.Color{240, 10, 10, 190},
			)
		}
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

var Colors = [...]rl.Color{
	{240, 240, 240, 255}, // Background/dead
	{50, 80, 80, 255},    // Alive
}

type State int32

const (
	StatePaused State = iota
	StateRunning
)

func Min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int32) int32 {
	if x < y {
		return y
	}
	return x
}

// Rect is used both for the board structure and the stored creatures.
type Rect struct {
	buf    []int32
	width  int32
	height int32
}

func NewRect(width, height int32) Rect {
	return Rect{make([]int32, width*height), width, height}
}

func (r Rect) At(x, y int32) int32 {
	// We add the width/height before calculating the remainder to handle the
	// case when x or y is negative.
	x = (x + r.width) % r.width
	y = (y + r.height) % r.height
	return r.buf[y*r.width+x]
}

func (r Rect) AtPtr(x, y int32) *int32 {
	x = (x + r.width) % r.width
	y = (y + r.height) % r.height
	return &r.buf[y*r.width+x]
}

func (r Rect) Neighbors(x, y int32) int32 {
	sum := int32(0)
	for yy := y - 1; yy <= y+1; yy++ {
		for xx := x - 1; xx <= x+1; xx++ {
			if yy == y && xx == x {
				continue
			}
			sum += r.At(xx, yy)
		}
	}
	return sum
}

func (r Rect) Zero() {
	for i := range r.buf {
		r.buf[i] = 0
	}
}

func (r Rect) Init() {
	r.Add(Creatures[GliderRight], 0, 0)
	r.Add(Creatures[GliderRight], 100, 100)
	r.Add(Creatures[GliderLeft], r.width-6, 0)
	r.Add(Creatures[PrePulsar], r.width/2-7, r.height-24)
	r.Add(Creatures[Pulsar], r.width-20, r.height-20)
}

func (r Rect) Add(c Rect, x0, y0 int32) {
	for y := int32(0); y < c.height; y++ {
		for x := int32(0); x < c.width; x++ {
			*r.AtPtr(x0+x, y0+y) = c.At(x, y)
		}
	}
}

type CreatureKind int32

const (
	GliderRight CreatureKind = iota
	GliderLeft
	Spaceship
	PrePulsar
	Pulsar
	QueenBee
	CreatureCount
)

var Creatures = [...]Rect{
	GliderRight: {
		buf: []int32{
			1, 0, 0,
			0, 1, 1,
			1, 1, 0,
		},
		width:  3,
		height: 3,
	},
	GliderLeft: {
		buf: []int32{
			0, 0, 1,
			1, 1, 0,
			0, 1, 1,
		},
		width:  3,
		height: 3,
	},
	Spaceship: {
		buf: []int32{
			1, 0, 0, 1, 0,
			0, 0, 0, 0, 1,
			0, 0, 0, 0, 1,
			0, 1, 1, 1, 1,
		},
		width:  5,
		height: 4,
	},
	PrePulsar: {
		buf: []int32{
			0, 1, 0,
			1, 1, 1,
			1, 0, 1,
			1, 1, 1,
			0, 1, 0,
		},
		width:  3,
		height: 5,
	},
	Pulsar: {
		buf: []int32{
			0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
			0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0,
			1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0,
		},
		width:  13,
		height: 13,
	},
	QueenBee: {
		buf: []int32{
			1, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
			0, 0, 0, 1,
			0, 0, 0, 1,
			0, 0, 1, 0,
			1, 1, 0, 0,
		},
		width:  4,
		height: 7,
	},
}
