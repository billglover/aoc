package password

import (
	"fmt"
	"testing"
)

func TestPasswordCount(t *testing.T) {
	tcs := []struct {
		pw    int
		partB bool
		meets bool
	}{
		{pw: 111111, partB: false, meets: true},
		{pw: 223450, partB: false, meets: false},
		{pw: 123789, partB: false, meets: false},
		{pw: 112233, partB: true, meets: true},
		{pw: 123444, partB: true, meets: false},
		{pw: 111122, partB: true, meets: true},
		{pw: 111111, partB: true, meets: false},
		{pw: 223333, partB: true, meets: true},
	}

	for _, tc := range tcs {
		got := Check(tc.pw, tc.partB)
		if got != tc.meets {
			t.Errorf("%d: %v != %v", tc.pw, got, tc.meets)
		}
	}
}

func TestPartOne(t *testing.T) {
	b1 := 168630
	b2 := 718098

	count := 0
	for i := b1; i <= b2; i++ {
		if Check(i, false) {
			count++
		}
	}
	fmt.Println("Part One:", count)
}

func TestPartTwo(t *testing.T) {
	b1 := 168630
	b2 := 718098

	count := 0
	for i := b1; i <= b2; i++ {
		if Check(i, true) {
			count++
		}
	}
	fmt.Println("Part Two:", count)
}
