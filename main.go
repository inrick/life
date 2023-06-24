package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	CellSize     = 10
	Width        = ScreenWidth / CellSize
	Height       = ScreenHeight / CellSize
)

type State int32

const (
	StatePaused State = iota
	StateRunning
)

var (
	Board     [Height][Width]int
	BoardCopy [Height][Width]int

	NrFramesPerStep = 10

	BackgroundColor = rl.Color{240, 240, 240, 255}
	AliveColor      = rl.Color{50, 80, 80, 255}
)

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Life")
	rl.SetTargetFPS(60)
	BoardInit()
	shouldClose := false
	state := StateRunning
	var frameNr uint
	var mousePos rl.Vector2
	var selected Creature
	for !shouldClose {
		// State handling
		shouldClose = rl.WindowShouldClose() || rl.IsKeyPressed(rl.KeyQ)
		mousePos = rl.GetMousePosition()
		mx, my := int(mousePos.X/CellSize), int(mousePos.Y/CellSize)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			if !(0 <= mx && mx <= Width && 0 <= my && my <= Height) {
				panic(struct{ mx, my int }{mx, my})
			}
			if len(selected.buf) == 0 {
				Board[my][mx] = (Board[my][mx] + 1) % 2
				state = StatePaused
			} else {
				BoardAdd(selected, mx, my)
				selected = Creature{}
				state = StateRunning
			}
		}
		switch {
		case rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeySpace):
			state = State((int(state) + 1) % 2)
		case rl.IsKeyPressed(rl.KeyMinus):
			NrFramesPerStep = Min(NrFramesPerStep+1, 60)
		case rl.IsKeyPressed(rl.KeyEqual):
			NrFramesPerStep = Max(NrFramesPerStep-1, 1)
		case rl.IsKeyPressed(rl.KeyZero):
			BoardZero()
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
		rl.ClearBackground(BackgroundColor)

		if state == StateRunning {
			if frameNr%uint(NrFramesPerStep) == 0 {
				for y := 0; y < Height; y++ {
					for x := 0; x < Width; x++ {
						alive := CellAt(x, y)
						n := CellNeighbors(x, y)
						if (alive == 1 && (n == 2 || n == 3)) || (alive == 0 && n == 3) {
							BoardCopy[y][x] = 1
						} else {
							BoardCopy[y][x] = 0
						}
					}
				}
				copy(Board[:], BoardCopy[:])
			}
			frameNr++
		}

		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				var cellColor rl.Color
				//switch selected.buf[y*selected.width+x] {
				switch Board[y][x] {
				case 0:
					cellColor = BackgroundColor
				case 1:
					cellColor = AliveColor
				}
				x0 := x * CellSize
				y0 := y * CellSize
				rl.DrawRectangle(int32(x0), int32(y0), CellSize, CellSize, cellColor)
			}
		}

		if len(selected.buf) > 0 {
			for y := 0; y < selected.height; y++ {
				for x := 0; x < selected.width; x++ {
					var cellColor rl.Color
					switch selected.buf[y*selected.width+x] {
					case 0:
						cellColor = BackgroundColor
					case 1:
						cellColor = AliveColor
					}
					x0 := ((x + mx) % Width) * CellSize
					y0 := ((y + my) % Height) * CellSize
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

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

type Creature struct {
	buf    []int
	width  int
	height int
}

type CreatureKind int

const (
	GliderRight CreatureKind = iota
	GliderLeft
	Spaceship
	PrePulsar
	Pulsar
	QueenBee
	CreatureCount
)

var Creatures = []Creature{
	GliderRight: {
		buf: []int{
			1, 0, 0,
			0, 1, 1,
			1, 1, 0,
		},
		width:  3,
		height: 3,
	},
	GliderLeft: {
		buf: []int{
			0, 0, 1,
			1, 1, 0,
			0, 1, 1,
		},
		width:  3,
		height: 3,
	},
	Spaceship: {
		buf: []int{
			1, 0, 0, 1, 0,
			0, 0, 0, 0, 1,
			0, 0, 0, 0, 1,
			0, 1, 1, 1, 1,
		},
		width:  5,
		height: 4,
	},
	PrePulsar: {
		buf: []int{
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
		buf: []int{
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
		buf: []int{
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

func CellAt(x, y int) int {
	x = (x + Width) % Width
	y = (y + Height) % Height
	return Board[y][x]
}

func CellNeighbors(x, y int) int {
	sum := 0
	for yy := y - 1; yy <= y+1; yy++ {
		for xx := x - 1; xx <= x+1; xx++ {
			if yy == y && xx == x {
				continue
			}
			sum += CellAt(xx, yy)
		}
	}
	return sum
}

func BoardZero() {
	for y := range Board {
		for x := range Board[y] {
			Board[y][x] = 0
		}
	}
}

func BoardInit() {
	BoardAdd(Creatures[GliderRight], 0, 0)
	BoardAdd(Creatures[GliderRight], 100, 100)
	BoardAdd(Creatures[GliderLeft], Width-6, 0)
	BoardAdd(Creatures[PrePulsar], Width/2-7, Height-24)
	BoardAdd(Creatures[Pulsar], Width-20, Height-20)
}

func BoardAdd(c Creature, x0, y0 int) {
	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			Board[(y+y0)%Height][(x+x0)%Width] = c.buf[y*c.width+x]
		}
	}
}
