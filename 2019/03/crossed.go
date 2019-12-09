package crossed

import (
	"math"
	"strconv"
	"strings"
)

type xy struct {
	x int
	y int
}

func Move(o xy, m string, g map[xy]int) (xy, int) {
	g[o] = 0
	d := 0

	ms := strings.Split(m, ",")
	for _, m := range ms {
		dir := m[:1]
		n, err := strconv.Atoi(m[1:])
		if err != nil {
			panic(err)
		}

		for s := 0; s < n; s++ {
			switch dir {
			case "R":
				o.x++
			case "U":
				o.y++
			case "L":
				o.x--
			case "D":
				o.y--
			}
			d++
			if d > g[o] {
				g[o] = d
			}
		}

	}

	return o, d
}

func Check(ps []string) (int, int, error) {
	gs := make([]map[xy]int, len(ps))
	d := 0
	l := 0

	for i, p := range ps {
		gs[i] = map[xy]int{}
		Move(xy{0, 0}, p, gs[i])
	}

	minx, maxx := math.MaxInt64, math.MinInt64
	miny, maxy := math.MaxInt64, math.MinInt64
	for _, g := range gs {
		for coord := range g {
			if coord.x > maxx {
				maxx = coord.x
			}
			if coord.x < minx {
				minx = coord.x
			}
			if coord.y > maxy {
				maxy = coord.y
			}
			if coord.y < miny {
				miny = coord.y
			}
		}
	}

	for x := minx; x < maxx; x++ {
		for y := miny; y < maxy; y++ {
			found := true
			for _, g := range gs {
				if g[xy{x, y}] == 0 {
					found = false
				}
			}
			if found == false {
				continue
			}

			mh := AbsInt(x) + AbsInt(y)
			if mh == 0 {
				continue
			}

			if d == 0 {
				d = mh
			}

			if mh < d {
				d = mh
			}

			tl := 0
			for _, g := range gs {
				tl += g[xy{x, y}]
			}

			if l == 0 {
				l = tl
			}

			if tl < l {
				l = tl
			}
		}
	}

	return d, l, nil
}

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
