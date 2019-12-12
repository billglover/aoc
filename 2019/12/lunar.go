package lunar

// Pos is the position of a point in 3D space.
type Pos struct {
	x int
	y int
	z int
}

// Vel is the velocity of a point in 3D space.
type Vel struct {
	x int
	y int
	z int
}

// A Moon has both position and velocity.
type Moon struct {
	P Pos
	V Vel
}

// Move takes a starting position and a velocity and returns the resulting ending
// position after a single time step has elapsed.
func move(sP Pos, v Vel) Pos {
	eP := Pos{
		x: sP.x + v.x,
		y: sP.y + v.y,
		z: sP.z + v.z,
	}
	return eP
}

// ApplyGravity takes a slice of Moons and updates their velocities according to
// the laws of gravity.
func ApplyGravity(ms []Moon, steps int) {

	g := map[int]Vel{}

	for m := range ms {
		g[m] = ms[m].V
	}

	for s := 0; s < steps; s++ {
		for m1 := range ms {
			for m2 := range ms {
				if m2 <= m1 {
					continue
				}
				v1 := g[m1]
				v2 := g[m2]

				vX := 0
				if ms[m1].P.x > ms[m2].P.x {
					vX = -1
				}
				if ms[m1].P.x < ms[m2].P.x {
					vX = 1
				}

				v1.x += vX
				v2.x -= vX

				vY := 0
				if ms[m1].P.y > ms[m2].P.y {
					vY = -1
				}
				if ms[m1].P.y < ms[m2].P.y {
					vY = 1
				}
				v1.y += vY
				v2.y -= vY

				vZ := 0
				if ms[m1].P.z > ms[m2].P.z {
					vZ = -1
				}
				if ms[m1].P.z < ms[m2].P.z {
					vZ = 1
				}
				v1.z += vZ
				v2.z -= vZ

				g[m1] = v1
				g[m2] = v2
			}
		}
		for m := 0; m < len(g); m++ {
			ms[m].V = g[m]

			ms[m].P.x += ms[m].V.x
			ms[m].P.y += ms[m].V.y
			ms[m].P.z += ms[m].V.z
		}
	}
}

// TotalEnergy returns the sum of the Potential and Kinetic Energy in the solar system.
func TotalEnergy(ms []Moon) int {
	te := 0
	for i := range ms {
		pe := abs(ms[i].P.x) + abs(ms[i].P.y) + abs(ms[i].P.z)
		ke := abs(ms[i].V.x) + abs(ms[i].V.y) + abs(ms[i].V.z)
		te += (pe * ke)
	}
	return te
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// RepeatCount returns the number of steps required before the universe starts
// repeating itself.
func RepeatCount(ms []Moon) int {

	initial := make([]Moon, len(ms))
	copy(initial, ms)

	g := map[int]Vel{}

	for m := range ms {
		g[m] = ms[m].V
	}

	latchX, latchY, latchZ := 0, 0, 0

	for s := 1; s < 10000000000; s++ {
		for m1 := range ms {
			for m2 := range ms {

				if m2 <= m1 {
					continue
				}
				v1 := g[m1]
				v2 := g[m2]

				vX := 0
				if ms[m1].P.x > ms[m2].P.x {
					vX = -1
				}
				if ms[m1].P.x < ms[m2].P.x {
					vX = 1
				}

				v1.x += vX
				v2.x -= vX

				vY := 0
				if ms[m1].P.y > ms[m2].P.y {
					vY = -1
				}
				if ms[m1].P.y < ms[m2].P.y {
					vY = 1
				}
				v1.y += vY
				v2.y -= vY

				vZ := 0
				if ms[m1].P.z > ms[m2].P.z {
					vZ = -1
				}
				if ms[m1].P.z < ms[m2].P.z {
					vZ = 1
				}
				v1.z += vZ
				v2.z -= vZ

				g[m1] = v1
				g[m2] = v2
			}
		}

		// loop over all the planets and update positions
		for m := 0; m < len(g); m++ {
			ms[m].V = g[m]

			ms[m].P.x += ms[m].V.x
			ms[m].P.y += ms[m].V.y
			ms[m].P.z += ms[m].V.z
		}

		// assume that all axes are repeating until we prove othrewise
		repeatX := true
		repeatY := true
		repeatZ := true

		// loop over all planets and look for differences to original position
		for m := 0; m < len(ms); m++ {

			// check for repeating state
			if !(ms[m].P.x == initial[m].P.x && ms[m].V.x == initial[m].V.x) {
				repeatX = false
			}
			if ms[m].P.y != initial[m].P.y || ms[m].V.y != initial[m].V.y {
				repeatY = false
			}
			if ms[m].P.z != initial[m].P.z || ms[m].V.z != initial[m].V.z {
				repeatZ = false
			}

		}
		if repeatX && latchX == 0 {
			latchX = s
		}
		if repeatY && latchY == 0 {
			latchY = s
		}
		if repeatZ && latchZ == 0 {
			latchZ = s
		}

		if latchX != 0 && latchY != 0 && latchZ != 0 {
			return cycle(latchX, latchY, latchZ)
		}
	}
	return 0
}

func cycle(x, y, z int) int {
	l := lcm(x, y)
	l = lcm(l, z)
	return l
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}
