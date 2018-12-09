package marble

import (
	"container/ring"
	"fmt"
)

func highScore(players int, marbles int) int {

	r := ring.New(1)
	nextMarble, nextPlayer := 1, 0
	r.Value = 0

	scores := make([]int, players)

	scores = play(r, marbles, nextMarble, nextPlayer, scores)

	maxScore := 0
	for p := range scores {
		if scores[p] > maxScore {
			maxScore = scores[p]
		}
	}

	return maxScore
}

func play(r *ring.Ring, marbles, initalMarble, player int, scores []int) []int {

	for marble := initalMarble; marble <= marbles; marble++ {

		// handle special (multiples of 23) marbles
		if marble%23 == 0 {
			scores[player] += marble
			r = r.Move(-7)
			scores[player] += r.Value.(int)
			s := r
			r = s.Link(r)

		} else {

			nr := ring.New(1)
			nr.Value = marble

			// advance the ring, insert the new marble, and rewind
			r = r.Next()
			r = r.Link(nr)
			r = r.Prev()
		}

		player = (player + 1) % len(scores)
	}

	return scores
}

func printRing(r *ring.Ring) {
	r.Do(func(p interface{}) {
		fmt.Printf("%d ", p.(int))
	})
	fmt.Println()
}
