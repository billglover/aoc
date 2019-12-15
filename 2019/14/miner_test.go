package miner

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		in  string
		out map[string]chemical
	}{
		{
			in: `10 ORE => 10 A
			1 ORE => 1 B
			7 A, 1 B => 1 C
			7 A, 1 C => 1 D
			7 A, 1 D => 1 E
			7 A, 1 E => 1 FUEL`,
			out: map[string]chemical{
				"A":    chemical{id: "A", out: 10, ingredients: map[string]int{"ORE": 10}},
				"B":    chemical{id: "B", out: 1, ingredients: map[string]int{"ORE": 1}},
				"C":    chemical{id: "C", out: 1, ingredients: map[string]int{"A": 7, "B": 1}},
				"D":    chemical{id: "D", out: 1, ingredients: map[string]int{"A": 7, "C": 1}},
				"E":    chemical{id: "E", out: 1, ingredients: map[string]int{"A": 7, "D": 1}},
				"FUEL": chemical{id: "FUEL", out: 1, ingredients: map[string]int{"A": 7, "E": 1}},
			},
		},
	}

	for _, tc := range tcs {
		reactions := parse(tc.in)
		if reflect.DeepEqual(reactions, tc.out) == false {
			t.Errorf("\n%v !=\n%v", reactions, tc.out)
		}
	}
}

func TestOreRequred(t *testing.T) {
	tcs := []struct {
		in  string
		out int
	}{
		{
			in: `10 ORE => 10 A
			1 ORE => 1 B
			7 A, 1 B => 1 C
			7 A, 1 C => 1 D
			7 A, 1 D => 1 E
			7 A, 1 E => 1 FUEL`,
			out: 31,
		},
		{
			in: `9 ORE => 2 A
			8 ORE => 3 B
			7 ORE => 5 C
			3 A, 4 B => 1 AB
			5 B, 7 C => 1 BC
			4 C, 1 A => 1 CA
			2 AB, 3 BC, 4 CA => 1 FUEL`,
			out: 165,
		},
		{
			in: `157 ORE => 5 NZVS
			165 ORE => 6 DCFZ
			44 XJWVT, 5 KHKGT, 1 QDVJ, 29 NZVS, 9 GPVTF, 48 HKGWZ => 1 FUEL
			12 HKGWZ, 1 GPVTF, 8 PSHF => 9 QDVJ
			179 ORE => 7 PSHF
			177 ORE => 5 HKGWZ
			7 DCFZ, 7 PSHF => 2 XJWVT
			165 ORE => 2 GPVTF
			3 DCFZ, 7 NZVS, 5 HKGWZ, 10 PSHF => 8 KHKGT`,
			out: 13312,
		},
		{
			in: `2 VPVL, 7 FWMGM, 2 CXFTF, 11 MNCFX => 1 STKFG
			17 NVRVD, 3 JNWZP => 8 VPVL
			53 STKFG, 6 MNCFX, 46 VJHF, 81 HVMC, 68 CXFTF, 25 GNMV => 1 FUEL
			22 VJHF, 37 MNCFX => 5 FWMGM
			139 ORE => 4 NVRVD
			144 ORE => 7 JNWZP
			5 MNCFX, 7 RFSQX, 2 FWMGM, 2 VPVL, 19 CXFTF => 3 HVMC
			5 VJHF, 7 MNCFX, 9 VPVL, 37 CXFTF => 6 GNMV
			145 ORE => 6 MNCFX
			1 NVRVD => 8 CXFTF
			1 VJHF, 6 MNCFX => 4 RFSQX
			176 ORE => 6 VJHF`,
			out: 180697,
		},
		{
			in: `171 ORE => 8 CNZTR
			7 ZLQW, 3 BMBT, 9 XCVML, 26 XMNCP, 1 WPTQ, 2 MZWV, 1 RJRHP => 4 PLWSL
			114 ORE => 4 BHXH
			14 VRPVC => 6 BMBT
			6 BHXH, 18 KTJDG, 12 WPTQ, 7 PLWSL, 31 FHTLT, 37 ZDVW => 1 FUEL
			6 WPTQ, 2 BMBT, 8 ZLQW, 18 KTJDG, 1 XMNCP, 6 MZWV, 1 RJRHP => 6 FHTLT
			15 XDBXC, 2 LTCX, 1 VRPVC => 6 ZLQW
			13 WPTQ, 10 LTCX, 3 RJRHP, 14 XMNCP, 2 MZWV, 1 ZLQW => 1 ZDVW
			5 BMBT => 4 WPTQ
			189 ORE => 9 KTJDG
			1 MZWV, 17 XDBXC, 3 XCVML => 2 XMNCP
			12 VRPVC, 27 CNZTR => 2 XDBXC
			15 KTJDG, 12 BHXH => 5 XCVML
			3 BHXH, 2 VRPVC => 7 MZWV
			121 ORE => 7 VRPVC
			7 XCVML => 6 RJRHP
			5 BHXH, 4 VRPVC => 5 LTCX`,
			out: 2210736,
		},
	}

	for _, tc := range tcs {
		r := parse(tc.in)
		ore := ore(1, r)
		if ore != tc.out {
			t.Errorf("%d != %d", ore, tc.out)
		}
	}
}

func TestLotsOfOre(t *testing.T) {
	tcs := []struct {
		in  string
		out int
	}{
		{
			in: `157 ORE => 5 NZVS
			165 ORE => 6 DCFZ
			44 XJWVT, 5 KHKGT, 1 QDVJ, 29 NZVS, 9 GPVTF, 48 HKGWZ => 1 FUEL
			12 HKGWZ, 1 GPVTF, 8 PSHF => 9 QDVJ
			179 ORE => 7 PSHF
			177 ORE => 5 HKGWZ
			7 DCFZ, 7 PSHF => 2 XJWVT
			165 ORE => 2 GPVTF
			3 DCFZ, 7 NZVS, 5 HKGWZ, 10 PSHF => 8 KHKGT`,
			out: 82892753,
		},
		{
			in: `2 VPVL, 7 FWMGM, 2 CXFTF, 11 MNCFX => 1 STKFG
			17 NVRVD, 3 JNWZP => 8 VPVL
			53 STKFG, 6 MNCFX, 46 VJHF, 81 HVMC, 68 CXFTF, 25 GNMV => 1 FUEL
			22 VJHF, 37 MNCFX => 5 FWMGM
			139 ORE => 4 NVRVD
			144 ORE => 7 JNWZP
			5 MNCFX, 7 RFSQX, 2 FWMGM, 2 VPVL, 19 CXFTF => 3 HVMC
			5 VJHF, 7 MNCFX, 9 VPVL, 37 CXFTF => 6 GNMV
			145 ORE => 6 MNCFX
			1 NVRVD => 8 CXFTF
			1 VJHF, 6 MNCFX => 4 RFSQX
			176 ORE => 6 VJHF`,
			out: 5586022,
		},
		{
			in: `171 ORE => 8 CNZTR
			7 ZLQW, 3 BMBT, 9 XCVML, 26 XMNCP, 1 WPTQ, 2 MZWV, 1 RJRHP => 4 PLWSL
			114 ORE => 4 BHXH
			14 VRPVC => 6 BMBT
			6 BHXH, 18 KTJDG, 12 WPTQ, 7 PLWSL, 31 FHTLT, 37 ZDVW => 1 FUEL
			6 WPTQ, 2 BMBT, 8 ZLQW, 18 KTJDG, 1 XMNCP, 6 MZWV, 1 RJRHP => 6 FHTLT
			15 XDBXC, 2 LTCX, 1 VRPVC => 6 ZLQW
			13 WPTQ, 10 LTCX, 3 RJRHP, 14 XMNCP, 2 MZWV, 1 ZLQW => 1 ZDVW
			5 BMBT => 4 WPTQ
			189 ORE => 9 KTJDG
			1 MZWV, 17 XDBXC, 3 XCVML => 2 XMNCP
			12 VRPVC, 27 CNZTR => 2 XDBXC
			15 KTJDG, 12 BHXH => 5 XCVML
			3 BHXH, 2 VRPVC => 7 MZWV
			121 ORE => 7 VRPVC
			7 XCVML => 6 RJRHP
			5 BHXH, 4 VRPVC => 5 LTCX`,
			out: 460664,
		},
	}

	for _, tc := range tcs {
		r := parse(tc.in)
		ore := lotsOfOre(r, 1000000000000)
		if ore != tc.out {
			t.Errorf("%d != %d", ore, tc.out)
		}
	}
}

func TestPartOne(t *testing.T) {
	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	r := parse(string(d))
	ore := ore(1, r)
	fmt.Println("Part One:", ore)
}

func TestPartTwo(t *testing.T) {
	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	r := parse(string(d))
	ore := lotsOfOre(r, 1000000000000)
	fmt.Println("Part Two:", ore)
}
