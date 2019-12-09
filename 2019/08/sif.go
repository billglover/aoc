package sif

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Layer [][]int
type Image []Layer

func NewLayer(width, height int, data string) (Layer, error) {
	l := make(Layer, height)
	if len(data) != width*height {
		return l, errors.New("unexpected data length")
	}
	for h := 0; h < height; h++ {
		l[h] = make([]int, width)
		for w := 0; w < width; w++ {
			v, err := strconv.Atoi(string(data[(h*width)+w]))
			if err != nil {
				return l, err
			}
			l[h][w] = v
		}
	}
	return l, nil
}

func (l *Layer) countVals() map[int]int {
	count := make(map[int]int)
	for h := range *l {
		for w := range (*l)[h] {
			val := (*l)[h][w]
			if _, ok := count[val]; !ok {
				count[val] = 0
			}
			count[val]++
		}
	}
	return count
}

func NewImage(width, height int, data string) (Image, error) {
	i := Image{}
	step := width * height
	l := 0
	for p := 0; p < len(data); p += step {
		nl, err := NewLayer(width, height, data[p:p+step])
		if err != nil {
			return i, err
		}
		i = append(i, nl)
		l++
	}

	return i, nil
}

func minZeros(ls Image) (int, int) {
	minLayer := -1
	minZeros := math.MaxInt64
	var minCount map[int]int
	for i := range ls {
		count := ls[i].countVals()
		if count[0] < minZeros {
			minZeros = count[0]
			minLayer = i
			minCount = count
		}
	}
	sum := minCount[1] * minCount[2]
	return minLayer, sum
}

func (i Image) Render() Layer {
	height := len(i[0])
	width := len(i[0][0])

	r := make(Layer, height)
	for h := 0; h < height; h++ {
		r[h] = make([]int, width)
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			for l := range i {
				if i[l][h][w] == 2 {
					continue
				}
				r[h][w] = i[l][h][w]
				break
			}
		}
	}

	return r
}

func (l Layer) Print() {
	height := len(l)
	width := len(l[0])

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			switch l[h][w] {
			case 0:
				fmt.Print("⬜️")
			case 1:
				fmt.Print("⬛️")
			default:
				fmt.Print("?")
			}
		}
		fmt.Println()
	}
}
