package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type xy struct {
	x int
	y int
}

type sat struct {
	pos xy
	vel xy
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
	sats := make([]sat, len(lines))

	re := regexp.MustCompile(`<\s*(-?\d+),\s*(-?\d+)\s*>.*<\s*(-?\d+),\s*(-?\d+)\s*>`)

	for l := range lines {
		m := re.FindStringSubmatch(lines[l])

		x, err := strconv.Atoi(m[1])
		if err != nil {
			fmt.Println("unable to parse x:", err)
			os.Exit(1)
		}

		y, err := strconv.Atoi(m[2])
		if err != nil {
			fmt.Println("unable to parse y:", err)
			os.Exit(1)
		}

		dx, err := strconv.Atoi(m[3])
		if err != nil {
			fmt.Println("unable to parse dx:", err)
			os.Exit(1)
		}

		dy, err := strconv.Atoi(m[4])
		if err != nil {
			fmt.Println("unable to parse dy:", err)
			os.Exit(1)
		}

		s := sat{
			pos: xy{x: x, y: y},
			vel: xy{x: dx, y: dy},
		}

		sats[l] = s
	}

	min, max := bounds(sats)
	minArea := area(min, max)

	t := 0
	for {
		if t > 3000000 {
			break
		}

		sats = advance(sats, false)
		curMin, curMax := bounds(sats)
		curArea := area(curMin, curMax)
		if curArea < minArea {
			minArea = curArea
		}

		// We assume the message appears when satelites are closes together
		if curArea > minArea {
			sats = advance(sats, true)
			printSky(sats, curMin, curMax)
			fmt.Printf("after %ds\n", t)
			break
		}

		t++
	}
}

func advance(sats []sat, reverse bool) []sat {
	for s := range sats {
		if reverse {
			sats[s].pos.x -= sats[s].vel.x
			sats[s].pos.y -= sats[s].vel.y
		} else {
			sats[s].pos.x += sats[s].vel.x
			sats[s].pos.y += sats[s].vel.y
		}
	}
	return sats
}

func printSky(sats []sat, min, max xy) {
	coords := map[xy]bool{}

	for s := range sats {
		coords[sats[s].pos] = true
	}

	for y := min.y; y <= max.y; y++ {
		for x := min.x; x <= max.x; x++ {
			if _, ok := coords[xy{x: x, y: y}]; ok == true {
				fmt.Printf("# ")
				continue
			}
			fmt.Printf(". ")
		}
		fmt.Println()
	}
}

func bounds(sats []sat) (xy, xy) {
	maxX, maxY, minX, minY := math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64
	for s := range sats {
		if sats[s].pos.x > maxX {
			maxX = sats[s].pos.x
		}
		if sats[s].pos.x < minX {
			minX = sats[s].pos.x
		}
		if sats[s].pos.y > maxY {
			maxY = sats[s].pos.y
		}
		if sats[s].pos.y < minY {
			minY = sats[s].pos.y
		}
	}
	return xy{minX, minY}, xy{maxX, maxY}
}

func area(min, max xy) int {
	dx := max.x - min.x
	dy := max.y - min.y
	area := dx * dy
	if area < 0 {
		area = -area
	}
	return area
}
