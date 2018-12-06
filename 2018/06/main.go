package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type coord struct {
	x int
	y int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("unable to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("unable to read file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	coords := make([]coord, len(lines))
	maxX, maxY := 0, 0
	for l := range lines {
		var x, y int
		fmt.Sscanf(lines[l], "%d, %d", &y, &x)
		coords[l] = coord{x: x, y: y}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	maxBlastRadius := partOne(maxX, maxY, coords)
	fmt.Println("Part One:", maxBlastRadius)

	regionSize := partTwo(maxX, maxY, coords, 10000)
	fmt.Println("Part Two:", regionSize)
}

// PartOne solves the first part of the puzzle.
func partOne(maxX, maxY int, coords []coord) int {

	// Create a grid where the value at each location is the index of the
	// nearest coordinate in the coords slice. An index of -1 is used if
	// no single coordinate is nearer than any other.
	grid := make([][]int, maxX+1)
	for x := range grid {
		grid[x] = make([]int, maxY+1)
		for y := range grid[x] {
			grid[x][y] = closestPoint(x, y, coords)
		}
	}

	// Identify coordinates that are unbounded, i.e. as the grid extends to
	// infinity, these coordinates remain closest to all future grid locations.
	infCoords := map[int]bool{}
	for _, x := range []int{0, len(grid) - 1} {
		for y := range grid[x] {
			if c := grid[x][y]; c >= 0 {
				infCoords[c] = true
			}
		}
	}

	for x := range grid[0] {
		for _, y := range []int{0, len(grid[x]) - 1} {
			if c := grid[x][y]; c >= 0 {
				infCoords[c] = true
			}
		}
	}

	// Count the number of times each coordinate index appears in the grid
	// to identify the coordinate that is the shortest distance to the largest
	// area. Ignore unbounded coordinates (infCoords) and locations where no
	// individual coordinate is closer than any other.
	coordCount := make([]int, len(coords))
	maxCoord := 0
	for x := range grid {
		for y := range grid[x] {
			c := grid[x][y]
			if _, ok := infCoords[c]; ok == false && c >= 0 {
				coordCount[c]++
				if coordCount[c] > maxCoord {
					maxCoord = coordCount[c]
				}
			}
		}
	}

	return maxCoord
}

// PartTwo solves the second part of the puzzle.
func partTwo(maxX, maxY int, coords []coord, threshold int) int {

	// Create a grid where the value at each location represents the sum
	// of the distances from that location to each of the provided coordinates.
	grid := make([][]int, maxX+1)
	for x := range grid {
		grid[x] = make([]int, maxY+1)
		for y := range grid[x] {
			grid[x][y] = totalDist(x, y, coords)
		}
	}

	// We are looking for locations where the total distance is less than the
	// magic threshold provided in the problem.
	count := 0
	for x := range grid {
		for y := range grid[x] {
			if grid[x][y] < threshold {
				count++
			}
		}
	}
	return count
}

// ClosestPoint takes an x and a y coordinate and returns the index of the closes
// coordinate in the provided []coord slice. An index of -1 is used if no single
// coordinate is nearer than any other.
func closestPoint(x, y int, coords []coord) int {
	minDistance := math.MaxInt64
	minCount := 0
	minLoc := -1
	for c := range coords {
		d := distance(x, y, coords[c].x, coords[c].y)
		if d == minDistance {
			minCount++
			continue
		}
		if d < minDistance {
			minDistance = d
			minLoc = c
			minCount = 1
		}
	}

	if minCount > 1 {
		return -1
	}
	return minLoc
}

// TotalDist takes an x and y coordinate and returns the sum of the distance to
// each of the coordinates in the provided []coord slice.
func totalDist(x, y int, coords []coord) int {
	d := 0
	for c := range coords {
		d += distance(x, y, coords[c].x, coords[c].y)
	}

	return d
}

// Distance returns the Manhattan distance between two points.
func distance(x1, y1, x2, y2 int) int {
	dX := x2 - x1
	if dX < 0 {
		dX = -dX
	}
	dY := y2 - y1
	if dY < 0 {
		dY = -dY
	}
	return dX + dY
}
