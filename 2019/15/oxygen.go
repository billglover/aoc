package oxygen

import (
	"fmt"
	"math"

	"github.com/billglover/aoc/2019/intcode"
)

type xy struct {
	x int
	y int
}

const (
	Wall = iota
	Space
	Target
)

var visited = map[xy]int{}
var grid = map[xy]int{}

func locate(p []int) int {
	m := intcode.Mem(p)
	in := make(chan int)
	out, done := intcode.Run(m, in)

	// start at the origin
	pos := xy{0, 0}
	visited[pos] = 0
	grid[pos] = Space

	dir := 4
	step := 0
	for step < 10000 {

		// make the next move
		//fmt.Println(step, "direction", dir)
		in <- dir

		nextPos := pos
		switch dir {
		case 1:
			nextPos.y--
		case 2:
			nextPos.y++
		case 3:
			nextPos.x--
		case 4:
			nextPos.x++
		}

		status := <-out
		//fmt.Println(step, "status:", status)

		switch status {
		case 0:
			grid[nextPos] = Wall
			//fmt.Println(step, "turning left")
			dir = turnLeft(dir)
		case 1:
			pos = nextPos
			grid[nextPos] = Space

			// if the next square hasn't been visited, store the number of steps
			// if it has use the lowest number of steps
			// update the step counter to be the lowest number of steps
			nextStep, seen := visited[nextPos]
			if !seen {
				visited[nextPos] = step
			} else {
				if nextStep < step {
					step = nextStep
				}
				visited[nextPos] = step
			}
			//fmt.Println(step, "turning right")
			dir = turnRight(dir)
			step++
		case 2:
			//fmt.Println(step, "found target")
			step++
			return step
		}
		//print(grid, visited, pos)
		//fmt.Println()
		//time.Sleep(time.Millisecond * 25)
	}

	close(in)
	<-done

	return 0
}

func turnRight(dir int) int {
	switch dir {
	case 1:
		return 4
	case 2:
		return 3
	case 3:
		return 1
	case 4:
		return 2
	}
	return -1
}

func turnLeft(dir int) int {
	switch dir {
	case 1:
		return 3
	case 2:
		return 4
	case 3:
		return 2
	case 4:
		return 1
	}
	return -1
}

func print(g, v map[xy]int, pos xy) {
	minX, maxX, minY, maxY := math.MaxInt64, math.MinInt64, math.MaxInt64, math.MinInt64
	for p := range g {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	if pos.x < minX {
		minX = pos.x
	}
	if pos.x > maxX {
		maxX = pos.x
	}
	if pos.y < minY {
		minY = pos.y
	}
	if pos.y > maxY {
		maxY = pos.y
	}

	fmt.Printf("Grid: {%d, %d}, {%d, %d}\n", minX, minY, maxX, maxY)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if x == pos.x && y == pos.y {
				fmt.Print("D")
				continue
			}

			switch g[xy{x, y}] {
			case 0:
				if _, ok := v[xy{x, y}]; !ok {
					fmt.Print(" ")
					continue
				}
				fmt.Print("#")
			case 1:
				fmt.Print(".")
			case 2:
				fmt.Print("O")
			default:
				fmt.Print("?")
			}
			fmt.Print()
		}
		fmt.Println()
	}
}

func flood(p []int) (map[xy]int, xy) {
	m := intcode.Mem(p)
	in := make(chan int)
	out, done := intcode.Run(m, in)

	// start at the origin
	pos := xy{}
	target := xy{}
	visited[pos] = 0
	grid[pos] = Space

	dir := 4
	step := 0
	for step < 10000 {
		in <- dir

		nextPos := pos
		switch dir {
		case 1:
			nextPos.y--
		case 2:
			nextPos.y++
		case 3:
			nextPos.x--
		case 4:
			nextPos.x++
		}

		status := <-out

		switch status {
		case 0:
			grid[nextPos] = Wall
			visited[nextPos] = math.MaxInt64
			dir = turnLeft(dir)
		case 1, 2:
			pos = nextPos
			grid[nextPos] = status
			nextStep, seen := visited[nextPos]
			if !seen {
				visited[nextPos] = step
			} else {
				if nextStep < step {
					step = nextStep
				}
				visited[nextPos] = step
			}
			dir = turnRight(dir)
			step++

			if status == 2 {
				target = pos
			}

			if nextPos.x == 0 && nextPos.y == 0 {
				return grid, target
			}
		}
	}

	close(in)
	<-done

	return grid, target
}

func fill(g map[xy]int, pos xy) int {
	space := 0
	unvisited := map[xy]bool{}
	distance := map[xy]int{}
	for l, s := range grid {
		if s == Space {
			unvisited[l] = true
			distance[l] = math.MaxInt64
			space++
		}
	}
	distance[pos] = 0

	// loop here
	for iter := 0; iter < 1000; iter++ {
		neighbours := unvisitedNeighbours(unvisited, pos)
		for _, nb := range neighbours {
			nbDist := distance[pos] + 1
			if nbDist < distance[nb] {
				distance[nb] = nbDist
			}
		}

		if len(unvisited) == 1 {
			return distance[pos]
		}

		grid[pos] = 2
		delete(unvisited, pos)

		minDist := math.MaxInt64
		for u := range unvisited {
			if distance[u] < minDist {
				minDist = distance[u]
				pos = u
			}
		}
	}

	return 0
}

func unvisitedNeighbours(unvisited map[xy]bool, pos xy) []xy {
	neighbours := []xy{}
	if unvisited[xy{pos.x, pos.y - 1}] {
		neighbours = append(neighbours, xy{pos.x, pos.y - 1})
	}
	if unvisited[xy{pos.x, pos.y + 1}] {
		neighbours = append(neighbours, xy{pos.x, pos.y + 1})
	}
	if unvisited[xy{pos.x - 1, pos.y}] {
		neighbours = append(neighbours, xy{pos.x - 1, pos.y})
	}
	if unvisited[xy{pos.x + 1, pos.y}] {
		neighbours = append(neighbours, xy{pos.x + 1, pos.y})
	}

	return neighbours
}
