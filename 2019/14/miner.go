package miner

import (
	"fmt"
	"sort"
	"strings"
)

type chemical struct {
	id          string
	out         int
	ingredients map[string]int
}

func parse(in string) map[string]chemical {
	var quant int
	var chem string
	reactions := map[string]chemical{}

	lines := strings.Split(in, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		parts := strings.Split(l, " => ")
		inputs := strings.Split(parts[0], ",")

		output := strings.TrimSpace(parts[1])
		fmt.Sscanf(output, "%d %s", &quant, &chem)
		c := chemical{id: chem, out: quant, ingredients: map[string]int{}}

		for _, input := range inputs {
			input = strings.TrimSpace(input)
			fmt.Sscanf(input, "%d %s", &quant, &chem)
			c.ingredients[chem] = quant
		}

		reactions[c.id] = c
	}

	return reactions
}

func ore(fuel int, r map[string]chemical) int {
	c := map[string]int{"FUEL": fuel}
	for {
		for e := range c {

			// Keep a track of how much or each element we still need to
			// produce. We retain details of elements we overproduce (-ve)
			// in case we subsequently need more of these elements.
			if e != "ORE" && c[e] > 0 {
				a := (c[e]-1)/r[e].out + 1
				c[e] -= r[e].out * a

				for i := range r[e].ingredients {
					c[i] += r[e].ingredients[i] * a
				}
			}
		}

		// We keep producing elements until we have all the elements we need.
		done := true
		for e := range c {
			if c[e] > 0 && e != "ORE" {
				done = false
			}
		}
		if done {
			return c["ORE"]
		}
	}
}

func lotsOfOre(r map[string]chemical, o int) int {
	max := 0
	fuel := 1
	for max < o {

		max = ore(fuel, r)
		fuel *= 2
	}

	i := sort.Search(fuel, func(i int) bool {
		output := ore(i, r)
		return output >= o
	})

	return i - 1
}
