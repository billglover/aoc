package battlefield

import (
	"io"
	"io/ioutil"
	"math"
	"strings"

	"github.com/pkg/errors"
)

// Battlefield is a 2D plane where Elves and Goblins do battle. The value at
// grid location is either:
// - # a wall
// - . an open cavern
// - G a Goblin
// - E an Elf
type Battlefield struct {
	Map      [][]string
	Warriors Units
}

// NewFromReader reads a battlefield configuration from the io.Reader and
// returns a new Battlefield. An error is returned if unable to successfully
// parse the battlefield.
func NewFromReader(r io.Reader) (Battlefield, error) {

	bf := Battlefield{}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return bf, errors.Wrap(err, "unable to read input")
	}

	lines := strings.Split(string(data), "\n")
	bf.Map = make([][]string, len(lines))
	bf.Warriors = Units{}

	for l := range lines {
		squares := strings.Split(strings.TrimSpace(lines[l]), "")
		bf.Map[l] = make([]string, len(squares))

		for s := range squares {
			bf.Map[l][s] = squares[s]

			switch squares[s] {
			case "G":
				g := Unit{Loc: XY{s, l}, MemberOf: Goblin, Alive: true, AttackPower: 3, HitPoints: 200}
				bf.Warriors = append(bf.Warriors, &g)
			case "E":
				g := Unit{Loc: XY{s, l}, MemberOf: Elf, Alive: true, AttackPower: 3, HitPoints: 200}
				bf.Warriors = append(bf.Warriors, &g)
			}
		}
	}

	return bf, nil
}

// Enemies takes a warrior Unit and returns a slice of enemy Units.
func (bf *Battlefield) Enemies(u *Unit) Units {
	e := Units{}
	for w := range bf.Warriors {
		if u.MemberOf != bf.Warriors[w].MemberOf && bf.Warriors[w].Alive {
			e = append(e, bf.Warriors[w])
		}
	}
	return e
}

// InRange takes a warrior Unit and returns a slice of enemy Units that are in
// range for an attack.
func (bf *Battlefield) InRange(u *Unit) Units {
	e := Units{}
	for _, w := range bf.Enemies(u) {
		distX := (w.Loc.X - u.Loc.X)
		if distX < 0 {
			distX = -distX
		}
		distY := (w.Loc.Y - u.Loc.Y)
		if distY < 0 {
			distY = -distY
		}
		dist := distX + distY

		if dist == 1 {
			e = append(e, w)
		}
	}
	return e
}

// DistanceTo takes an XY coordinate and computes the distance to all other
// coordinates on the map. A negative distance indicates a map location is
// unreachable.
func (bf *Battlefield) DistanceTo(loc XY) [][]int {
	distance := make([][]int, len(bf.Map))

	visited := map[XY]int{}
	candidates := map[XY]int{}

	for y := range bf.Map {
		distance[y] = make([]int, len(bf.Map[y]))
		for x := range bf.Map[y] {
			distance[y][x] = int(math.MaxInt64)
		}
	}

	curLoc := loc
	candidates[curLoc] = 0

	for iter := 0; len(candidates) != 0; iter++ {
		visited[curLoc] = candidates[curLoc]
		distance[curLoc.Y][curLoc.X] = candidates[curLoc]
		delete(candidates, curLoc)

		eval := []XY{}

		if curLoc.Y-1 > 0 {
			eval = append(eval, XY{X: curLoc.X, Y: curLoc.Y - 1})
		}

		if curLoc.Y+1 < len(distance) {
			eval = append(eval, XY{X: curLoc.X, Y: curLoc.Y + 1})
		}

		if curLoc.X-1 > 0 {
			eval = append(eval, XY{X: curLoc.X - 1, Y: curLoc.Y})
		}

		if curLoc.X+1 < len(distance[curLoc.Y]) {
			eval = append(eval, XY{X: curLoc.X + 1, Y: curLoc.Y})
		}

		for _, nextLoc := range eval {

			if _, ok := visited[nextLoc]; ok {
				continue
			}

			nextDist := visited[curLoc] + 1

			// if the target location is occupied we can't reach it
			if bf.Map[nextLoc.Y][nextLoc.X] != "." {
				nextDist = -1
				visited[nextLoc] = nextDist
				continue
			}

			if curDist, ok := candidates[nextLoc]; ok {
				if curDist < nextDist {
					nextDist = curDist
				}
			}

			candidates[nextLoc] = nextDist
		}

		// pick the closest candidate for our next iteration
		shortestHop := int(math.MaxInt64)
		for loc, dist := range candidates {
			if dist == -1 {
				delete(candidates, loc)
				continue
			}

			if dist < shortestHop {
				shortestHop = dist
				curLoc = loc
			}
		}

	}

	for y := range distance {
		for x := range distance[y] {
			if distance[y][x] == int(math.MaxInt64) {
				distance[y][x] = -1
			}
		}
	}

	return distance
}

// Remaining returns the number of warriors who remain alive on the battlefield.
func (bf *Battlefield) Remaining() int {
	r := 0
	for w := range bf.Warriors {
		if bf.Warriors[w].Alive == true {
			r++
		}
	}
	return r
}

// SetElfPower allows us to tweak the power of the Elves.
func (bf *Battlefield) SetElfPower(p int) {
	for _, w := range bf.Warriors {
		if w.Alive && w.MemberOf == Elf {
			w.AttackPower = p
		}
	}
}
