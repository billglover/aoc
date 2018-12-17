package battlefield

// Clan is the clan a warrior belongs to.
type Clan int

func (c Clan) String() string {
	names := map[Clan]string{
		Elf:    "Elf",
		Goblin: "Goblin",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "Unknown"
}

// Clans are distinct and are never created or destroyed. They are constant.
const (
	Elf Clan = iota
	Goblin
)

// XY represents 2D coordinates
type XY struct {
	X int
	Y int
}

type XYs []XY

// Len returns the number of XYs. It partially implements the sort.Interface.
func (xys XYs) Len() int {
	return len(xys)
}

// Less indicates if one XY is located before another. It partially
// implements the sort.Interface.
func (xys XYs) Less(a, b int) bool {
	if xys[a].Y == xys[b].Y {
		return xys[a].X < xys[b].X
	}
	return xys[a].Y < xys[b].Y
}

// Swap changes the position of two Units. It partially implements the
// sort.Interface.
func (xys XYs) Swap(a, b int) {
	xys[a], xys[b] = xys[b], xys[a]
}

// Unit represents a character in the game.
type Unit struct {
	Loc         XY
	MemberOf    Clan
	Alive       bool
	AttackPower int
	HitPoints   int
}

// Units represents a collection of individual Units
type Units []*Unit

// ByLocation implements sort.Interface for Units based on the Loc field.
type ByLocation Units

// Len returns the number of units. It partially implements the sort.Interface.
func (u ByLocation) Len() int {
	return len(u)
}

// Less indicates if one unit is located before another. It partially
// implements the sort.Interface.
func (u ByLocation) Less(a, b int) bool {
	if u[a].Loc.Y == u[b].Loc.Y {
		return u[a].Loc.X < u[b].Loc.X
	}
	return u[a].Loc.Y < u[b].Loc.Y
}

// Swap changes the position of two Units. It partially implements the
// sort.Interface.
func (u ByLocation) Swap(a, b int) {
	u[a], u[b] = u[b], u[a]
}

// ByHitPoints implements sort.Interface for Units based on the HitPoints field.
type ByHitPoints Units

// Len returns the number of units. It partially implements the sort.Interface.
func (u ByHitPoints) Len() int {
	return len(u)
}

// Less indicates if one unit is located before another. It partially
// implements the sort.Interface.
func (u ByHitPoints) Less(a, b int) bool {
	return u[a].HitPoints < u[b].HitPoints
}

// Swap changes the position of two Units. It partially implements the
// sort.Interface.
func (u ByHitPoints) Swap(a, b int) {
	u[a], u[b] = u[b], u[a]
}
