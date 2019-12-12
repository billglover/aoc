package lunar

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMove(t *testing.T) {
	tcs := []struct {
		sP Pos
		v  Vel
		eP Pos
	}{
		{sP: Pos{1, 2, 3}, v: Vel{-2, 0, 3}, eP: Pos{-1, 2, 6}},
	}

	for _, tc := range tcs {
		got := move(tc.sP, tc.v)
		if got != tc.eP {
			t.Errorf("%v != %v", got, tc.eP)
		}
	}
}

func TestApplyGravity(t *testing.T) {
	tcs := []struct {
		Before []Moon
		Steps  int
		After  []Moon
	}{
		{
			Before: []Moon{
				Moon{P: Pos{-1, 0, 2}},
				Moon{P: Pos{2, -10, -7}},
				Moon{P: Pos{4, -8, 8}},
				Moon{P: Pos{3, 5, -1}},
			},
			Steps: 1,
			After: []Moon{
				Moon{P: Pos{2, -1, 1}, V: Vel{3, -1, -1}},
				Moon{P: Pos{3, -7, -4}, V: Vel{1, 3, 3}},
				Moon{P: Pos{1, -7, 5}, V: Vel{-3, 1, -3}},
				Moon{P: Pos{2, 2, 0}, V: Vel{-1, -3, 1}},
			},
		},
		{
			Before: []Moon{
				Moon{P: Pos{-1, 0, 2}},
				Moon{P: Pos{2, -10, -7}},
				Moon{P: Pos{4, -8, 8}},
				Moon{P: Pos{3, 5, -1}},
			},
			Steps: 10,
			After: []Moon{
				Moon{P: Pos{2, 1, -3}, V: Vel{-3, -2, 1}},
				Moon{P: Pos{1, -8, 0}, V: Vel{-1, 1, 3}},
				Moon{P: Pos{3, -6, 1}, V: Vel{3, 2, -3}},
				Moon{P: Pos{2, 0, 4}, V: Vel{1, -1, -1}},
			},
		},
	}

	for _, tc := range tcs {
		ApplyGravity(tc.Before, tc.Steps)
		if reflect.DeepEqual(tc.Before, tc.After) == false {
			t.Errorf("\n%+v !=\n%+v", tc.Before, tc.After)
		}
	}
}

func TestTotalEnergy(t *testing.T) {
	tcs := []struct {
		Before []Moon
		Steps  int
		Energy int
	}{
		{
			Before: []Moon{
				Moon{P: Pos{-1, 0, 2}},
				Moon{P: Pos{2, -10, -7}},
				Moon{P: Pos{4, -8, 8}},
				Moon{P: Pos{3, 5, -1}},
			},
			Steps:  10,
			Energy: 179,
		},
		{
			Before: []Moon{
				Moon{P: Pos{-8, -10, 0}},
				Moon{P: Pos{5, 5, 10}},
				Moon{P: Pos{2, -7, 3}},
				Moon{P: Pos{9, -8, -3}},
			},
			Steps:  100,
			Energy: 1940,
		},
	}

	for _, tc := range tcs {
		ApplyGravity(tc.Before, tc.Steps)
		got := TotalEnergy(tc.Before)
		if got != tc.Energy {
			t.Errorf("%v != %v", got, tc.Energy)
		}
	}
}

func TestRepeatCount(t *testing.T) {
	tcs := []struct {
		Before []Moon
		Steps  int
		Energy int
	}{
		{
			Before: []Moon{
				Moon{P: Pos{-1, 0, 2}},
				Moon{P: Pos{2, -10, -7}},
				Moon{P: Pos{4, -8, 8}},
				Moon{P: Pos{3, 5, -1}},
			},
			Steps: 2772,
		},
		{
			Before: []Moon{
				Moon{P: Pos{-8, -10, 0}},
				Moon{P: Pos{5, 5, 10}},
				Moon{P: Pos{2, -7, 3}},
				Moon{P: Pos{9, -8, -3}},
			},
			Steps: 4686774924,
		},
	}

	for _, tc := range tcs {
		got := RepeatCount(tc.Before)
		if got != tc.Steps {
			t.Errorf("%v != %v", got, tc.Steps)
		}
	}
}

func TestPartOne(t *testing.T) {
	moons := []Moon{
		Moon{P: Pos{14, 15, -2}},
		Moon{P: Pos{17, -3, 4}},
		Moon{P: Pos{6, 12, -13}},
		Moon{P: Pos{-2, 10, -8}},
	}
	ApplyGravity(moons, 1000)
	energy := TotalEnergy(moons)
	fmt.Println("Part One:", energy)
}

func TestPartTwo(t *testing.T) {
	moons := []Moon{
		Moon{P: Pos{14, 15, -2}},
		Moon{P: Pos{17, -3, 4}},
		Moon{P: Pos{6, 12, -13}},
		Moon{P: Pos{-2, 10, -8}},
	}
	cycle := RepeatCount(moons)
	fmt.Println("Part Two:", cycle)
}
