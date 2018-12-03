package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type claim struct {
	id int
	x  int
	y  int
	w  int
	h  int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("unable to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	claims, err := parseClaims(f)
	if err != nil {
		fmt.Println("unable to parse claims:", err)
		os.Exit(1)
	}

	var grid [1000][1000]int

	for _, c := range claims {
		for x := 0; x < c.w; x++ {
			for y := 0; y < c.h; y++ {
				grid[c.x+x][c.y+y] = grid[c.x+x][c.y+y] + 1
			}
		}
	}

	count := 0
	for x := range grid {
		for y := range grid[x] {
			if grid[x][y] > 1 {
				count++
			}
		}
	}
	fmt.Println("overlapping squares:", count)

	clear := true
	for _, c := range claims {
		for x := 0; x < c.w; x++ {
			for y := 0; y < c.h; y++ {
				if grid[c.x+x][c.y+y] > 1 {
					clear = false
				}
			}
		}
		if clear {
			fmt.Println("claim with no conflict:", c.id)
		}
		clear = true
	}
}

func parseClaims(r io.Reader) ([]claim, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	claims := make([]claim, len(lines))

	for i, l := range lines {
		id, x, y, w, h := 0, 0, 0, 0, 0
		_, err := fmt.Sscanf(l, "#%d @ %d,%d: %dx%d\n", &id, &x, &y, &w, &h)
		if err != nil {
			return nil, err
		}
		claims[i] = claim{id: id, x: x, y: y, w: w, h: h}
	}

	return claims, nil
}
