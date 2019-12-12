package robot

import (
	"fmt"
	"math"

	intcode "github.com/billglover/aoc/2019/intcode"
)

type XY struct {
	x int
	y int
}

type Robot struct {
	d Direction
	p XY
}

type Direction XY
type Turn int

var (
	North = Direction{1, 0}
	East  = Direction{0, 1}
	South = Direction{-1, 0}
	West  = Direction{0, -1}
)

const (
	Left Turn = iota
	Right
)

func (r *Robot) Move(d int, t Turn) {
	switch t {
	case Left:
		r.d.x, r.d.y = r.d.y, -r.d.x
	case Right:
		r.d.x, r.d.y = -r.d.y, r.d.x
	}
	delta := XY{r.d.x * d, r.d.y * d}
	r.p.x, r.p.y = r.p.x+delta.x, r.p.y+delta.y
}

func Paint(p []int, initC int) map[XY]int {
	m := intcode.Mem(p)
	artwork := map[XY]int{}
	r := Robot{
		p: XY{0, 0},
		d: North,
	}

	artwork[r.p] = initC

	camera := make(chan int, 1)
	out, done := intcode.Run(m, camera)

	for {
		curC := artwork[r.p]
		camera <- curC

		newC := <-out
		artwork[r.p] = newC
		direction := <-out
		r.Move(1, Turn(direction))

		select {
		case <-done:
			return artwork
		default:
			continue
		}
	}
}

func Render(i map[XY]int) {
	minX, maxX := math.MaxInt64, math.MinInt64
	minY, maxY := math.MaxInt64, math.MinInt64
	for l := range i {
		if l.x > maxX {
			maxX = l.x
		}
		if l.y > maxY {
			maxY = l.y
		}
		if l.x < minX {
			minX = l.x
		}
		if l.y < minY {
			minY = l.y
		}
	}

	for x := maxX; x >= minX; x-- {
		for y := minY; y <= maxY; y++ {

			if i[XY{x, y}] == 1 {
				fmt.Print("⬜️")
			} else {
				fmt.Print("⬛️")
			}
		}
		fmt.Println()
	}
}
