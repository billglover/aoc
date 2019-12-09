package sif

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewLayer(t *testing.T) {
	tcs := []struct {
		data   string
		width  int
		height int
		out    Layer
	}{
		{data: "123456", width: 3, height: 2, out: Layer{[]int{1, 2, 3}, []int{4, 5, 6}}},
		{data: "789012", width: 3, height: 2, out: Layer{[]int{7, 8, 9}, []int{0, 1, 2}}},
	}

	for _, tc := range tcs {
		got, err := NewLayer(tc.width, tc.height, tc.data)
		if err != nil {
			t.Error(err)
		}

		if reflect.DeepEqual(got, tc.out) == false {
			t.Errorf("%v != %v", got, tc.out)
		}
	}
}

func TestCountVals(t *testing.T) {
	tcs := []struct {
		l     Layer
		count map[int]int
	}{
		{
			l:     Layer{[]int{7, 7, 0}, []int{0, 1, 2}},
			count: map[int]int{0: 2, 1: 1, 2: 1, 7: 2},
		},
		{
			l:     Layer{[]int{7, 7, 0}, []int{0, 1, 2}},
			count: map[int]int{0: 2, 1: 1, 2: 1, 7: 2},
		},
	}

	for _, tc := range tcs {
		count := tc.l.countVals()
		if reflect.DeepEqual(count, tc.count) == false {
			t.Errorf("%v != %v", count, tc.count)
		}
	}
}

func TestNewImage(t *testing.T) {
	tcs := []struct {
		data   string
		width  int
		height int
		out    Image
	}{
		{
			data: "123456789012", width: 3, height: 2,
			out: Image{
				Layer{[]int{1, 2, 3}, []int{4, 5, 6}},
				Layer{[]int{7, 8, 9}, []int{0, 1, 2}},
			},
		},
	}
	for _, tc := range tcs {
		got, err := NewImage(tc.width, tc.height, tc.data)
		if err != nil {
			t.Error(err)
		}

		if reflect.DeepEqual(got, tc.out) == false {
			t.Errorf("%v != %v", got, tc.out)
		}
	}
}

func TestMinZeros(t *testing.T) {
	tcs := []struct {
		i        []Layer
		minLayer int
		sum      int
	}{
		{
			i: []Layer{
				Layer{[]int{1, 2, 3}, []int{4, 5, 6}},
				Layer{[]int{7, 8, 9}, []int{0, 1, 2}},
			},
			minLayer: 0,
			sum:      1,
		},
		{
			i: []Layer{
				Layer{[]int{1, 2, 3}, []int{1, 1, 0}},
				Layer{[]int{0, 8, 9}, []int{0, 1, 2}},
			},
			minLayer: 0,
			sum:      3,
		},
	}

	for _, tc := range tcs {
		l, sum := minZeros(tc.i)
		if l != tc.minLayer {
			t.Errorf("%v != %v", l, tc.minLayer)
		}

		if sum != tc.sum {
			t.Errorf("%v != %v", sum, tc.sum)
		}
	}
}

func TestPartOne(t *testing.T) {
	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	data := string(b)
	i, err := NewImage(25, 6, data)
	if err != nil {
		t.Fatal(err)
	}

	_, sum := minZeros(i)
	fmt.Println("Part One:", sum)
}

func TestRender(t *testing.T) {
	tcs := []struct {
		i Image
		l Layer
	}{
		{
			i: Image{
				Layer{[]int{0, 2}, []int{2, 2}},
				Layer{[]int{1, 1}, []int{2, 2}},
				Layer{[]int{2, 2}, []int{1, 2}},
				Layer{[]int{0, 0}, []int{0, 0}}},
			l: Layer{[]int{0, 1}, []int{1, 0}},
		},
	}

	for _, tc := range tcs {
		tc.i.Render()
	}
}

func TestPartTwo(t *testing.T) {
	f, err := os.Open("input.txt")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	data := string(b)
	i, err := NewImage(25, 6, data)
	if err != nil {
		t.Fatal(err)
	}

	l := i.Render()
	fmt.Println("Part Two:")
	l.Print()
}
