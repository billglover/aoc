package marble

import "testing"

var testCases = []struct {
	players int
	marbles int
	score   int
}{
	{players: 9, marbles: 25, score: 32},
	{players: 10, marbles: 1618, score: 8317},
	{players: 13, marbles: 7999, score: 146373},
	{players: 17, marbles: 1104, score: 2764},
	{players: 21, marbles: 6111, score: 54718},
	{players: 30, marbles: 5807, score: 37305},
}

func TestHighScore(t *testing.T) {
	for _, tc := range testCases {
		score := highScore(tc.players, tc.marbles)
		if score != tc.score {
			t.Errorf("players: %d, marbles: %d - got %d, want %d", tc.players, tc.marbles, score, tc.score)
		}
	}
}

func TestPuzzleAnswerPartOne(t *testing.T) {
	score := highScore(452, 71250)
	t.Logf("Part One: %d", score)
}

func TestPuzzleAnswerPartTwo(t *testing.T) {
	score := highScore(452, 7125000)
	t.Logf("Part Two: %d", score)
}
