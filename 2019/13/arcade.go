package arcade

import (
	"github.com/billglover/aoc/2019/intcode"
)

// Sprite represents a graphic in the arcade game.
type Sprite int

// There are a fixed number of sprites used in the game.
const (
	Empty Sprite = iota
	Wall
	Block
	Paddle
	Ball
)

// XY is a point in 2D space.
type XY struct {
	x int
	y int
}

// Frame is a single 2D frame of sprites.
type Frame map[XY]Sprite

// DrawFrame takes an IntCode program and renders a frame of the arcade game and
// returns the score to be shown on the segment display.
func DrawFrame(p []int, initC int) (Frame, int) {
	m := intcode.Mem(p)
	f := Frame{}
	score := 0

	in := make(chan int, 1)
	out, done := intcode.Run(m, in)

	c := 0
	op := [3]int{}

	for {
		select {
		case o := <-out:
			op[c] = o
			c++

			if c != 3 {
				continue
			}

			x, y, s := op[0], op[1], op[2]
			if x == -1 && y == 0 {
				score = s
			} else {
				f[XY{x, y}] = Sprite(s)
			}
			c = 0

		case <-done:
			return f, score
		}
	}
}

// PlayGame runs the arcade game.
func PlayGame(p []int) int {
	m := intcode.Mem(p)
	f := Frame{}
	score := 0

	// set memory location 0 to 2 to play for free
	m[0] = 2

	in := make(chan int)
	out, done := intcode.Run(m, in)

	// joystick movements
	// neutral: 0
	// left: -1
	// right: 1

	c := 0
	op := [3]int{}
	prevPos := XY{}
	ballPos := XY{}
	paddlePos := XY{}
	joystick := 0

	for {
		joystick = movePaddle(prevPos, ballPos, paddlePos)
		select {
		case o := <-out:
			op[c] = o
			c++

			if c != 3 {
				continue
			}

			x, y, s := op[0], op[1], op[2]
			if x == -1 && y == 0 {
				score = s
			} else {
				sp := Sprite(s)
				f[XY{x, y}] = sp
				switch sp {
				case Ball:
					prevPos = ballPos
					ballPos = XY{x, y}
				case Paddle:
					paddlePos = XY{x, y}
				}
			}
			c = 0

		case in <- joystick:
		case <-done:
			return score
		}
	}
}

func movePaddle(ba, bb, p XY) int {
	switch {
	case p.x > bb.x:
		return -1
	case p.x < bb.x:
		return 1
	}
	return 0
}
