package uom

import (
	"fmt"
	"os"
	"strings"
)

type body struct {
	id       string
	children []*body
	parent   *body
}

type bodies map[string]*body

func parseInput(i string) bodies {
	bodies := bodies{}

	lines := strings.Split(i, "\n")
	for _, line := range lines {
		ids := strings.Split(line, ")")
		if len(ids) != 2 {
			fmt.Printf("unexpected body count: %d, %s", len(ids), line)
			os.Exit(1)
		}

		id1 := ids[0]
		b1, ok := bodies[id1]
		if !ok {
			b1 = &body{id: id1, children: []*body{}}
			bodies[id1] = b1
		}

		id2 := ids[1]
		b2, ok := bodies[id2]
		if !ok {
			b2 = &body{id: id2, children: []*body{}}
			bodies[id2] = b2
		}

		b1.children = append(b1.children, b2)
		b2.parent = b1
	}

	_, ok := bodies["COM"]
	if !ok {
		fmt.Println("no root node", bodies)
		os.Exit(1)
	}
	return bodies
}

func orbits(i string) int {
	bodies := parseInput(i)

	c := 0
	for _, b := range bodies {
		c += countParents(b)
	}
	return c
}

func findSanta(i string) int {
	bodies := parseInput(i)

	youToCOM := map[string]int{}

	c := -1
	for b := bodies["YOU"]; ; b = b.parent {
		youToCOM[b.id] = c
		c++
		if b.id == "COM" {
			break
		}
	}

	c = -1
	for b := bodies["SAN"]; ; b = b.parent {
		yc, ok := youToCOM[b.id]
		if ok {
			return yc + c
		}
		c++
		if b.id == "COM" {
			break
		}
	}

	return -1
}

func countChildren(b *body) int {
	count := 0
	for c := range b.children {
		count++
		count += countChildren(b.children[c])
	}
	return count
}

func countParents(b *body) int {
	if b.parent == nil {
		return 0
	}
	return 1 + countParents(b.parent)
}
