package password

func Check(i int, partB bool) bool {
	c := i

	decreasing := true
	seq := 0
	double := false
	pd := 10

	for c > 0 {
		d := c % 10

		if d > pd {
			decreasing = false
		}
		if !partB {
			if pd == d {
				double = true
			}
		} else {
			if pd != d {
				if seq == 2 {
					double = true
				}
				seq = 1
			} else {
				seq++
			}
		}
		pd = d
		c /= 10
	}

	if partB && seq == 2 {
		double = true
	}

	return double && decreasing
}
