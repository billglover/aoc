package uom

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestOrbitalMap(t *testing.T) {
	tcs := []struct {
		input  string
		orbits int
	}{
		{input: "COM)AAA", orbits: 1},
		{input: "COM)AAA\nCOM)BBB", orbits: 2},
		{input: "COM)AAA\nAAA)BBB", orbits: 3},
		{input: "COM)B\nB)C\nC)D", orbits: 6},
		{input: "COM)B\nB)C\nC)D\nD)E\nE)F\nB)G\nG)H\nD)I\nE)J\nJ)K\nK)L", orbits: 42},
	}

	for _, tc := range tcs {
		o := orbits(tc.input)
		if o != tc.orbits {
			t.Errorf("%d != %d", o, tc.orbits)
		}
	}
}

func TestFindSanta(t *testing.T) {
	input := `COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
K)YOU
I)SAN`

	o := findSanta(input)
	if o != 4 {
		t.Errorf("%d != %d", o, 4)
	}
}

func TestPartOne(t *testing.T) {
	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	o := orbits(string(input))
	fmt.Println("Part One:", o)
}

func TestPartTwo(t *testing.T) {
	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	o := findSanta(string(input))
	fmt.Println("Part Two:", o)
}

func TestCountChildren(t *testing.T) {
	b1 := body{id: "b1", children: []*body{}}
	b2 := body{id: "b2", children: []*body{}}
	b1.children = append(b1.children, &b2)
	b2.parent = &b1

	c := countChildren(&b1)
	if c != 1 {
		t.Errorf("%d != %d", c, 1)
	}
}

func TestCountParents(t *testing.T) {
	b1 := body{id: "b1", children: []*body{}}
	b2 := body{id: "b2", children: []*body{}}
	b1.children = append(b1.children, &b2)
	b2.parent = &b1

	c := countParents(&b2)
	if c != 1 {
		t.Errorf("%d != %d", c, 1)
	}
}
