package battlefield

import (
	"math"
	"sort"
)

// PlayBattle does battle until no more moves are possible. It returns the final
// score and the number of Elf deaths.
func (bf *Battlefield) PlayBattle() (int, int, int) {
	round := 0
	score := 0
	elfDeaths := 0

	for ; bf.PlayRound() == false; round++ {
	}

	for _, u := range bf.Warriors {
		if u.Alive == true {
			score += u.HitPoints
		}
		if u.MemberOf == Elf && u.Alive == false {
			elfDeaths++
		}
	}

	return score, round, elfDeaths
}

// PlayRound simulates one round of combat, it returns a bool to indicate if
// the battle ended.
func (bf *Battlefield) PlayRound() bool {

	// sort the warriors so that they compete in reading order
	sort.Sort(ByLocation(bf.Warriors))

	for _, w := range bf.Warriors {
		if w.Alive == false {
			continue
		}

		targets := bf.Enemies(w)

		if len(targets) == 0 {
			return true
		}

		// if not in range move
		inRange := bf.InRange(w)
		if len(inRange) == 0 {
			bf.Move(w)
		}

		// if in range attack
		inRange = bf.InRange(w)
		if len(inRange) > 0 {
			bf.Attack(w)
		}
	}
	return false
}

// Move moves a single warrior according to the rules of combat.
func (bf *Battlefield) Move(u *Unit) {

	// if we are already in range, return
	if len(bf.InRange(u)) > 0 {
		return
	}

	// find in-range squares
	attackPoints := map[XY]bool{}
	inRangeUnits := bf.Enemies(u)
	for _, w := range inRangeUnits {
		if w.Loc.Y-1 > 0 && bf.Map[w.Loc.Y-1][w.Loc.X] == "." {
			attackPoints[XY{X: w.Loc.X, Y: w.Loc.Y - 1}] = true
		}

		if w.Loc.Y+1 < len(bf.Map) && bf.Map[w.Loc.Y+1][w.Loc.X] == "." {
			attackPoints[XY{X: w.Loc.X, Y: w.Loc.Y + 1}] = true
		}

		if w.Loc.X-1 > 0 && bf.Map[w.Loc.Y][w.Loc.X-1] == "." {
			attackPoints[XY{X: w.Loc.X - 1, Y: w.Loc.Y}] = true
		}

		if w.Loc.X+1 < len(bf.Map[w.Loc.Y]) && bf.Map[w.Loc.Y][w.Loc.X+1] == "." {
			attackPoints[XY{X: w.Loc.X + 1, Y: w.Loc.Y}] = true
		}
	}

	// compute distance to possible attack squares
	distMap := bf.DistanceTo(u.Loc)
	distMin := int(math.MaxInt64)
	for ap := range attackPoints {
		if distMap[ap.Y][ap.X] != -1 && distMap[ap.Y][ap.X] < distMin {
			distMin = distMap[ap.Y][ap.X]
		}
	}

	closestAPs := XYs{}
	for ap := range attackPoints {
		if distMap[ap.Y][ap.X] == distMin {
			closestAPs = append(closestAPs, ap)
		}
	}

	// if there are no available attack points, end the turn
	if len(closestAPs) == 0 {
		return
	}

	// from those with the lowest range, pick first reading order
	sort.Sort(closestAPs)
	chosenAP := closestAPs[0]

	// compute distance to the chosen attack point
	distMap = bf.DistanceTo(chosenAP)

	// find minimum distance
	startPoints := map[XY]bool{}
	distMin = int(math.MaxInt64)

	if u.Loc.Y-1 > 0 {
		sp := XY{u.Loc.X, u.Loc.Y - 1}
		if distMap[sp.Y][sp.X] != -1 {
			if distMap[sp.Y][sp.X] < distMin {
				distMin = distMap[sp.Y][sp.X]
			}
			startPoints[sp] = true
		}
	}

	if u.Loc.Y+1 < len(bf.Map) {
		sp := XY{u.Loc.X, u.Loc.Y + 1}
		if distMap[sp.Y][sp.X] != -1 {
			if distMap[sp.Y][sp.X] < distMin {
				distMin = distMap[sp.Y][sp.X]
			}
			startPoints[sp] = true
		}
	}

	if u.Loc.X-1 > 0 {
		sp := XY{u.Loc.X - 1, u.Loc.Y}
		if distMap[sp.Y][sp.X] != -1 {
			if distMap[sp.Y][sp.X] < distMin {
				distMin = distMap[sp.Y][sp.X]
			}
			startPoints[sp] = true
		}
	}

	if u.Loc.X+1 < len(bf.Map[u.Loc.Y]) {
		sp := XY{u.Loc.X + 1, u.Loc.Y}
		if distMap[sp.Y][sp.X] != -1 {
			if distMap[sp.Y][sp.X] < distMin {
				distMin = distMap[sp.Y][sp.X]
			}
			startPoints[sp] = true
		}
	}

	nextXY := XYs{}
	for ap := range startPoints {
		if distMap[ap.Y][ap.X] == distMin {
			nextXY = append(nextXY, ap)
		}
	}

	// if there are no moves, end the turn
	if len(nextXY) == 0 {
		return
	}

	sort.Sort(nextXY)
	dest := nextXY[0]

	// Make the move and update the map
	bf.Map[u.Loc.Y][u.Loc.X] = "."
	u.Loc.Y, u.Loc.X = dest.Y, dest.X
	switch u.MemberOf {
	case Elf:
		bf.Map[u.Loc.Y][u.Loc.X] = "E"
	case Goblin:
		bf.Map[u.Loc.Y][u.Loc.X] = "G"
	}
}

// Attack attacks a single warrior according to the rules of combat.
func (bf *Battlefield) Attack(u *Unit) {

	// if we have no units in range, we can't attack
	targets := bf.InRange(u)
	if len(targets) == 0 {
		return
	}

	// select based on least hit points, then reading order
	sort.Sort(ByHitPoints(targets))

	minHitPoints := targets[0].HitPoints
	for t := range targets {
		if targets[t].HitPoints != minHitPoints {
			targets = targets[:t]
			break
		}
	}

	sort.Sort(ByLocation(targets))
	target := targets[0]

	target.HitPoints -= u.AttackPower

	if target.HitPoints <= 0 {
		target.Alive = false
		bf.Map[target.Loc.Y][target.Loc.X] = "."
	}
}
