package charge

import "testing"

func TestCellPower(t *testing.T) {
	power := cellPower(xy{3, 5}, 8)
	if power != 4 {
		t.Errorf("incorrect cell power: got %d, want %d", power, 4)
	}
}

var sqPowerCases = []struct {
	serial int
	power  int
	pos    xy
}{
	{serial: 18, power: 29, pos: xy{33, 45}},
	{serial: 42, power: 30, pos: xy{21, 61}},
}

func TestLocateMaxPower_Example(t *testing.T) {
	for _, tc := range sqPowerCases {
		power, pos := maxPower(3, tc.serial)
		if power != tc.power {
			t.Errorf("incorrect power: got %d, want %d", power, tc.power)
		}

		if pos != tc.pos {
			t.Errorf("incorrect coordinates: got (%d, %d), want (%d, %d)", pos.x, pos.y, tc.pos.x, tc.pos.y)
		}
	}
}

func TestLocateMaxPower(t *testing.T) {
	power, pos := maxPower(3, 2866)
	t.Logf("Power: %d", power)
	t.Logf("Location: (%d, %d)", pos.x, pos.y)
}

var sqPowerSizeCases = []struct {
	serial int
	power  int
	pos    xy
	size   int
}{
	{serial: 18, power: 113, pos: xy{90, 269}, size: 16},
	{serial: 42, power: 119, pos: xy{232, 251}, size: 12},
}

func TestLocateMaxPowerSize_Example(t *testing.T) {
	for _, tc := range sqPowerSizeCases {
		power, pos, size := maxPowerSize(tc.serial)
		if power != tc.power {
			t.Errorf("incorrect power: got %d, want %d", power, tc.power)
		}

		if pos != tc.pos {
			t.Errorf("incorrect coordinates: got (%d, %d), want (%d, %d)", pos.x, pos.y, tc.pos.x, tc.pos.y)
		}

		if size != tc.size {
			t.Errorf("incorrect size: got %d, want %d", size, tc.size)
		}
	}
}

func TestLocateMaxPowerSize(t *testing.T) {
	power, pos, size := maxPowerSize(2866)
	t.Logf("Power: %d", power)
	t.Logf("Location: (%d, %d)", pos.x, pos.y)
	t.Logf("Size: %d", size)
}
