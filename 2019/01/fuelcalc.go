package fuelcalc

import (
	"bufio"
	"io"
	"strconv"
)

func Total(r io.Reader) (int, error) {
	var f int

	s := bufio.NewScanner(r)

	for s.Scan() {
		m, err := strconv.Atoi(s.Text())
		if err != nil {
			return f, err
		}
		f += massToFuel(m)
	}

	if err := s.Err(); err != nil {
		return f, err
	}

	return f, nil
}

func TotalRecursive(r io.Reader) (int, error) {
	var f int

	s := bufio.NewScanner(r)

	for s.Scan() {
		m, err := strconv.Atoi(s.Text())
		if err != nil {
			return f, err
		}
		f += massToFuelTotal(m)
	}

	if err := s.Err(); err != nil {
		return f, err
	}

	return f, nil
}

func massToFuel(m int) int {
	f := m/3 - 2
	return f
}

func massToFuelTotal(m int) int {
	f := massToFuel(m)
	nfm := f
	for {
		nfm = massToFuel(nfm)
		if nfm < 0 {
			break
		}
		f += nfm
	}
	return f
}
