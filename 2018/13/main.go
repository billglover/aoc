package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var left = map[string]string{"^": "<", "<": "v", "v": ">", ">": "^"}
var right = map[string]string{"^": ">", "<": "^", "v": "<", ">": "v"}

var partB = flag.Bool("partB", false, "return result for part B")

func main() {
	flag.Parse()

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
	g1 := make([][]string, len(lines))

	for l := range lines {
		g1[l] = strings.Split(lines[l], "")
	}

	track, carts = extractCarts(g1)

	t := 0

	for len(carts) > 1 {

		// sort the carts so that we start with the top left
		order = cartOrder{}
		for c := range carts {
			order = append(order, c)
		}
		sort.Sort(order)

		// loop over the carts and increment the position
		for i := range order {
			cart := carts[order[i]]
			if cart != nil {
				cart.Tick(t)
			}
		}

		t++
	}

	for pos := range carts {
		fmt.Printf("Remaining cart at: (%d, %d)\n", pos.y, pos.x)
	}
}

func extractCarts(g [][]string) ([][]string, cartData) {
	// ASSUME: no carts start at an intersection
	carts := cartData{}
	for x := range g {
		for y := range g[x] {
			var c cart
			switch g[x][y] {
			case "<":
				c = cart{pos: xy{x: x, y: y}, turns: 0, dir: "<"}
				g[x][y] = "-"
			case ">":
				c = cart{pos: xy{x: x, y: y}, turns: 0, dir: ">"}
				g[x][y] = "-"
			case "^":
				c = cart{pos: xy{x: x, y: y}, turns: 0, dir: "^"}
				g[x][y] = "|"
			case "v":
				c = cart{pos: xy{x: x, y: y}, turns: 0, dir: "v"}
				g[x][y] = "|"
			default:
				continue
			}
			carts[c.pos] = &c
		}
	}
	return g, carts
}

func printGrid(g [][]string) {
	for x := range g {
		for y := range g[x] {
			fmt.Print(g[x][y])
		}
		fmt.Println()
	}
}

func printGridWithCarts(g [][]string, cs cartData, t int) {
	f, _ := os.Create(fmt.Sprintf("data/map_%05d.txt", t))
	defer f.Close()

	for x := range g {
		for y := range g[x] {

			found := false
			for c := range cs {
				// if cs[c].removed == true {
				// 	continue
				// }
				if cs[c].pos.x == x && cs[c].pos.y == y {
					fmt.Fprint(f, cs[c].dir)
					found = true
				}
			}

			if !found {
				fmt.Fprint(f, g[x][y])
			}
		}
		fmt.Fprintln(f)
	}
}
