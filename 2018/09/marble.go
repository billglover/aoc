package marble

import (
	"container/ring"
	"fmt"
)

func highScore(players int, marbles int) int {

	// create the initial ring
	r := ring.New(1)
	nextMarble, nextPlayer := 1, 0
	r.Value = 0

	// create a slice for scores
	scores := make([]int, players)

	// play the game
	scores = play(r, marbles, nextMarble, nextPlayer, scores)

	// print the scores
	maxScore := 0
	for p := range scores {
		//fmt.Printf("Player %d: %d\n", p+1, scores[p])
		if scores[p] > maxScore {
			maxScore = scores[p]
		}
	}

	return maxScore
}

func play(r *ring.Ring, marbles, marble, player int, scores []int) []int {
	// return if we've reached the last marble
	if marble > marbles {
		return scores
	}

	// if we have a special case
	if marble%23 == 0 {
		scores[player] += marble
		r = r.Move(-7)
		scores[player] += r.Value.(int)
		s := r
		r = s.Link(r)
	} else {
		// create a new ring containing the marble
		nr := ring.New(1)
		nr.Value = marble

		// advance the ring by one
		r = r.Next()

		// link the two rings and rewind
		r = r.Link(nr)
		r = r.Prev()
	}

	// print the ring
	//fmt.Printf("[%d] ", player+1)
	//printRing(r)

	// play the next round
	player = (player + 1) % len(scores)
	scores = play(r, marbles, marble+1, player, scores)

	return scores
}

func printRing(r *ring.Ring) {
	r.Do(func(p interface{}) {
		fmt.Printf("%d ", p.(int))
	})
	fmt.Println()
}
