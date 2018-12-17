package battlefield

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"testing"
)

var map1 = `#######
#.G.E.#
#E.G.E#
#.G.E.#
#######`

var map2 = `#######
#.E...#
#.....#
#...G.#
#######`

var map3 = `#######
#.GE..#
#.EG.E#
#.G.E.#
#######`

var map5 = `#######
#.G...#
#...EG#
#.#.#G#
#..G#E#
#.....#
#######`

var map6 = `#######
#G..#E#
#E#E.E#
#G.##.#
#...#E#
#...E.#
#######`

var map7 = `#######
#E..EG#
#.#G.E#
#E.##E#
#G..#.#
#..E#.#
#######`

var map8 = `#######
#E.G#.#
#.#G..#
#G.#.G#
#G..#.#
#...E.#
#######`

var map9 = `#######
#.E...#
#.#..G#
#.###.#
#E#G#G#
#...#G#
#######`

var map10 = `#########
#G......#
#.E.#...#
#..##..G#
#...##..#
#...#...#
#.G...G.#
#.....G.#
#########`

func TestNewFromReader(t *testing.T) {
	got, err := NewFromReader(bytes.NewBufferString(map1))

	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}

	want := [][]string{
		[]string{"#", "#", "#", "#", "#", "#", "#"},
		[]string{"#", ".", "G", ".", "E", ".", "#"},
		[]string{"#", "E", ".", "G", ".", "E", "#"},
		[]string{"#", ".", "G", ".", "E", ".", "#"},
		[]string{"#", "#", "#", "#", "#", "#", "#"},
	}

	if reflect.DeepEqual(got.Map, want) == false {
		t.Errorf("unexpected battlefield:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestUnits(t *testing.T) {
	bf, _ := NewFromReader(bytes.NewBufferString(map1))

	want := Units{
		&Unit{Loc: XY{2, 1}, MemberOf: Goblin},
		&Unit{Loc: XY{4, 1}, MemberOf: Elf},
		&Unit{Loc: XY{1, 2}, MemberOf: Elf},
		&Unit{Loc: XY{3, 2}, MemberOf: Goblin},
		&Unit{Loc: XY{5, 2}, MemberOf: Elf},
		&Unit{Loc: XY{2, 3}, MemberOf: Goblin},
		&Unit{Loc: XY{4, 3}, MemberOf: Elf},
	}

	got := bf.Warriors

	for i := range got {
		if got[i].Loc.X != want[i].Loc.X && got[i].Loc.Y != want[i].Loc.Y {
			t.Errorf("unexpected units:\ngot:  %v\nwant: %v", got, want)
		}

		if got[i].MemberOf != want[i].MemberOf {
			t.Errorf("unexpected clan: got %s, want %s", got[i].MemberOf, want[i].MemberOf)
		}

		if got[i].Alive != true {
			t.Errorf("unexpected dead unit: (%d, %d)", got[i].Loc.X, got[i].Loc.Y)
		}
	}
}

func TestEnemies(t *testing.T) {
	bf, _ := NewFromReader(bytes.NewBufferString(map1))
	bf.Warriors[5].Alive = false

	want := Units{
		&Unit{Loc: XY{2, 1}, MemberOf: Goblin},
		&Unit{Loc: XY{3, 2}, MemberOf: Goblin},
	}

	got := bf.Enemies(bf.Warriors[1])

	if len(got) != len(want) {
		t.Fatalf("unexpected number of enemies: got: %d, want: %d", len(got), len(want))
	}

	for i := range got {
		if got[i].Loc.X != want[i].Loc.X && got[i].Loc.Y != want[i].Loc.Y {
			t.Errorf("unexpected units:\ngot:  %v\nwant: %v", got, want)
		}

		if got[i].MemberOf != Goblin {
			t.Errorf("unexpected clan: got %s, want %s", got[i].MemberOf, want[i].MemberOf)
		}
	}
}

var inRangeCases = []struct {
	data string
	want Units
}{
	{data: map1, want: Units{}},
	{
		data: map3,
		want: Units{
			&Unit{Loc: XY{3, 1}},
			&Unit{Loc: XY{2, 2}},
		},
	},
	{data: map2, want: Units{}},
}

func TestInRange(t *testing.T) {
	for _, tc := range inRangeCases {
		bf, _ := NewFromReader(bytes.NewBufferString(tc.data))
		got := bf.InRange(bf.Warriors[0])

		if len(got) != len(tc.want) {
			t.Fatalf("unexpected number of enemies in range: got: %d, want: %d", len(got), len(tc.want))
		}

		for i := range got {
			if got[i].Loc.X != tc.want[i].Loc.X && got[i].Loc.Y != tc.want[i].Loc.Y {
				t.Errorf("unexpected units in range:\ngot:  %v\nwant: %v", got[i], tc.want[i])
			}
		}
	}

}

func TestDistanceTo(t *testing.T) {
	bf, _ := NewFromReader(bytes.NewBufferString(map2))

	want := [][]int{
		[]int{-1, -1, -1, -1, -1, -1, -1},
		[]int{-1, 4, -1, 2, 1, 2, -1},
		[]int{-1, 3, 2, 1, 0, 1, -1},
		[]int{-1, 4, 3, 2, -1, 2, -1},
		[]int{-1, -1, -1, -1, -1, -1, -1},
	}

	got := bf.DistanceTo(XY{4, 2})

	if reflect.DeepEqual(got, want) == false {
		fmt.Println("got:")
		printGrid(got)

		fmt.Println()

		fmt.Println("want:")
		printGrid(want)

		t.Error("unexpected distance map returned")
	}
}

func TestRemaining(t *testing.T) {
	bf, _ := NewFromReader(bytes.NewBufferString(map1))
	got := bf.Remaining()
	want := 7

	if got != want {
		t.Errorf("unexpected number of remaining warriors: got %d, want %d", got, want)
	}

	bf.Warriors[0].Alive = false

	got = bf.Remaining()
	want = 6

	if got != want {
		t.Errorf("unexpected number of remaining warriors: got %d, want %d", got, want)
	}
}

var moveCases = []struct {
	data    string
	warrior int
	newLoc  XY
}{
	{data: map2, warrior: 0, newLoc: XY{2, 1}},
}

func TestMove(t *testing.T) {
	for _, tc := range moveCases {
		bf, _ := NewFromReader(bytes.NewBufferString(tc.data))

		want := [][]string{
			[]string{"#", "#", "#", "#", "#", "#", "#"},
			[]string{"#", ".", ".", "E", ".", ".", "#"},
			[]string{"#", ".", ".", ".", ".", ".", "#"},
			[]string{"#", ".", ".", ".", "G", ".", "#"},
			[]string{"#", "#", "#", "#", "#", "#", "#"},
		}

		bf.Move(bf.Warriors[tc.warrior])
		if reflect.DeepEqual(bf.Map, want) == false {
			fmt.Println("got:")
			printBattlefield(bf.Map)
			fmt.Println()
			fmt.Println("want:")
			printBattlefield(want)
			t.Error("unexpected move made")
		}
	}
}

var battleCases = []struct {
	data  string
	score int
}{
	{data: map5, score: 27730},
	{data: map6, score: 36334},
	{data: map7, score: 39514},
	{data: map8, score: 27755},
	{data: map9, score: 28944},
	{data: map10, score: 18740},
}

func TestPlayBattle(t *testing.T) {
	for _, tc := range battleCases {
		bf, _ := NewFromReader(bytes.NewBufferString(tc.data))
		score, rounds, _ := bf.PlayBattle()

		if score*rounds != tc.score {
			t.Errorf("unexpected score: got %d, want %d", score*rounds, tc.score)
		}
	}
}

func TestPlayPartOne(t *testing.T) {
	bf, _ := NewFromReader(bytes.NewBufferString(input))
	score, rounds, _ := bf.PlayBattle()
	fmt.Println("Part One:", score*rounds)
}

func TestPlayPartTwo(t *testing.T) {
	elfDeaths := 1
	power := 4
	score := 0
	rounds := 0

	for ; elfDeaths > 0; power++ {
		bf, _ := NewFromReader(bytes.NewBufferString(input))
		bf.SetElfPower(power)
		score, rounds, elfDeaths = bf.PlayBattle()
	}
	fmt.Println("Part Two:", rounds*score)
	fmt.Println("Power:", power-1)
}

func printGrid(grid [][]int) {
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == int(math.MaxInt64) {
				fmt.Print(" . ")
				continue
			}
			fmt.Printf("%2d ", grid[y][x])
		}
		fmt.Println()
	}
}

func printBattlefield(grid [][]string) {
	for y := range grid {
		for x := range grid[y] {
			fmt.Printf("%s ", grid[y][x])
		}
		fmt.Println()
	}
}

var input = `################################
#############..#################
#############..#.###############
############G..G.###############
#############....###############
##############.#...#############
################..##############
#############G.##..#..##########
#############.##.......#..######
#######.####.G##.......##.######
######..####.G.......#.##.######
#####.....#..GG....G......######
####..###.....#####.......######
####.........#######..E.G..#####
####.G..G...#########....E.#####
#####....G.G#########.#...######
###........G#########....#######
##..#.......#########....##.E.##
##.#........#########.#####...##
#............#######..#.......##
#.G...........#####........E..##
#....G........G..G.............#
#..................E#...E...E..#
#....#...##...G...E..........###
#..###...####..........G###E.###
#.###########..E.......#########
#.###############.......########
#################.......########
##################....#..#######
##################..####.#######
#################..#####.#######
################################`
