package charge

const gridSize = 300

type xy struct {
	x int
	y int
}

func maxPower(size int, serial int) (int, xy) {
	grid := make([][]int, gridSize)

	for x := 1; x <= gridSize; x++ {
		grid[x-1] = make([]int, gridSize)
		for y := 1; y <= gridSize; y++ {
			grid[x-1][y-1] = cellPower(xy{x, y}, serial)
		}
	}

	power := 0
	pos := xy{}
	for x := 1; x <= gridSize-size+1; x++ {
		for y := 1; y <= gridSize-size+1; y++ {
			sqPower := 0
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					sqPower += grid[x-1+i][y-1+j]
				}
			}
			if sqPower > power {
				power, pos = sqPower, xy{x, y}
			}
		}
	}

	return power, pos
}

// MaxPowerSize is functionally correct but slow.
// TODO:
// - cache calculations to reduce loop count
// - consider concurrency
func maxPowerSize(serial int) (int, xy, int) {
	power := 0
	pos := xy{}
	size := 1
	for s := 1; s <= 300; s++ {
		sizePower, sizePos := maxPower(s, serial)
		if sizePower > power {
			power, pos, size = sizePower, sizePos, s
		}
	}
	return power, pos, size
}

func cellPower(pos xy, serial int) int {
	rackID := pos.x + 10
	power := rackID * pos.y
	power += serial
	power *= rackID

	power = (power / 100) % 10
	power -= 5
	return power
}
