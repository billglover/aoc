package main

type xy struct {
	x int
	y int
}

type cart struct {
	pos   xy
	turns int
	dir   string
}

func (c *cart) Tick(t int) {
	// remove the cart from the current map while in motion
	delete(carts, c.pos)

	// move the cart to the next location
	switch c.dir {
	case "^":
		c.pos.x--
	case "v":
		c.pos.x++
	case "<":
		c.pos.y--
	case ">":
		c.pos.y++
	}

	// look for a crash at the new cart location
	if c2, crash := carts[c.pos]; crash {
		delete(carts, c2.pos)
		return
	}

	// update cart direction based on track layout
	switch track[c.pos.x][c.pos.y] {
	case "\\":
		if c.dir == "^" || c.dir == "v" {
			c.dir = left[c.dir]
		} else {
			c.dir = right[c.dir]
		}
	case "/":
		if c.dir == "^" || c.dir == "v" {
			c.dir = right[c.dir]
		} else {
			c.dir = left[c.dir]
		}
	case "+":
		t := c.turns % 3
		switch t {
		case 0:
			c.dir = left[c.dir]
		case 1:

		case 2:
			c.dir = right[c.dir]
		}
		c.turns++
	}

	// add the cart back to the map in new location
	carts[c.pos] = c
}

type cartData map[xy]*cart
type cartOrder []xy

func (o cartOrder) Len() int {
	return len(o)
}

func (o cartOrder) Less(a, b int) bool {
	if o[a].x <= o[b].x {
		return o[a].y < o[b].y
	}
	return false
}

func (o cartOrder) Swap(a, b int) {
	o[a], o[b] = o[b], o[a]
}

// Global variales used because...
var carts cartData
var order cartOrder
var track = [][]string{}
