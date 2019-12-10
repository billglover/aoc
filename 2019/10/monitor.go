package monitor

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strings"
)

// XY is a coordinate in the form {x, y} where x is the distance
// from the left and y is the distance from the top.
type XY struct {
	x int
	y int
}

// Grid is a map of objects in space.
type grid map[XY]bool

// Map is a map of space.
type Map struct {
	Range XY
	Grid  grid
}

func (m Map) String() string {
	buf := bytes.Buffer{}
	for y := 0; y < m.Range.y; y++ {
		for x := 0; x < m.Range.x; x++ {
			if _, ok := m.Grid[XY{x, y}]; ok {
				buf.WriteRune('#')
				continue
			}
			buf.WriteRune('.')
		}
		if y != m.Range.y-1 {
			buf.WriteRune('\n')
		}
	}

	buf.WriteString(fmt.Sprintf(" {%d, %d}", m.Range.x, m.Range.y))
	return buf.String()
}

// Scan takes an input and returns a map of space.
func Scan(i string) Map {
	m := Map{
		Range: XY{},
		Grid:  grid{},
	}

	// figure out the Range of our map
	l := strings.Split(i, "\n")
	m.Range.y = len(l)
	m.Range.x = len(l[00])

	for y := 0; y < m.Range.y; y++ {
		for x := 0; x < m.Range.x; x++ {
			if l[y][x] == '#' {
				m.Grid[XY{x, y}] = true
			}
		}
	}

	return m
}

func visible(m Map, a XY) (map[float64]float64, map[float64]XY) {
	dist := map[float64]float64{}
	vis := map[float64]XY{}

	for b := range m.Grid {
		if b == a {
			continue
		}
		th, h := vect(a, b)
		if d, ok := dist[th]; ok {
			if d < h {
				dist[th] = d
				vis[th] = b
				continue
			}
		}
		dist[th] = h
		vis[th] = b
	}
	return dist, vis
}

// Locate takes a map and identifies the best location to position an asteroid
// monitoring station. It returns the coordinates and number of asteroids that
// would be visible.
func Locate(m Map) (XY, int) {
	candidates := map[XY]int{}
	for a := range m.Grid {
		vis, _ := visible(m, a)
		candidates[a] = len(vis)
	}

	v := math.MinInt64
	l := XY{}
	for p := range candidates {
		if candidates[p] > v {
			v = candidates[p]
			l = p
		}
	}

	return l, v
}

func vect(a, b XY) (float64, float64) {
	dx := float64(b.x - a.x)
	dy := float64(b.y - a.y)
	th := math.Atan2(dx, dy)
	h := math.Sqrt(dx*dx + dy*dy)
	return th, h
}

// VaporiseN takes a map and the location of a laser. It returns the location
// of the nth asteroid to be vaporised.
func VaporiseN(m Map, p XY, n int) XY {
	var target XY

	for n > 0 {
		_, dist := visible(m, p)

		keys := make([]float64, 0, len(dist))
		for k := range dist {
			keys = append(keys, k)
		}
		sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

		for _, k := range keys {
			delete(m.Grid, dist[k])
			n--
			if n == 0 {
				return dist[k]
			}
		}

		break
	}

	return target
}
