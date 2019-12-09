package fuelcalc

import (
	"fmt"
	"os"
	"testing"
)

func TestAoCSolution(t *testing.T) {
	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal("unable to open file:", err)
	}

	result, err := Total(f)
	if err != nil {
		t.Error("unable to calculate the fuel:", err)
	}
	fmt.Println("part one:", result)

	f.Seek(0, 0)

	result, err = TotalRecursive(f)
	if err != nil {
		t.Error("unable to calculate the fuel:", err)
	}
	fmt.Println("part two:", result)

	err = f.Close()
	if err != nil {
		t.Fatal("unable to close file:", err)
	}
}

func TestMassToFuel(t *testing.T) {
	tc := []struct {
		mass int
		fuel int
	}{
		{mass: 12, fuel: 2},
		{mass: 14, fuel: 2},
		{mass: 1969, fuel: 654},
		{mass: 100756, fuel: 33583},
	}

	for _, c := range tc {
		got := massToFuel(c.mass)
		if got != c.fuel {
			t.Errorf("mass %d: %d != %d", c.mass, got, c.fuel)
		}
	}
}

func TestMassToFuelTotal(t *testing.T) {
	tc := []struct {
		mass int
		fuel int
	}{
		{mass: 14, fuel: 2},
		{mass: 1969, fuel: 966},
		{mass: 100756, fuel: 50346},
	}

	for _, c := range tc {
		got := massToFuelTotal(c.mass)
		if got != c.fuel {
			t.Errorf("mass %d: %d != %d", c.mass, got, c.fuel)
		}
	}
}
