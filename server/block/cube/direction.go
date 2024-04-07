package cube

// Direction represents a direction towards one of the horizontal axes of the world.
type Direction int

const (
	// North represents the north direction, towards the negative Z.
	North Direction = iota
	// South represents the south direction, towards the positive Z.
	South
	// West represents the west direction, towards the negative X.
	West
	// East represents the east direction, towards the positive X.
	East
	// NorthEast represents the north east direction
	NorthEast
	// NorthWest represents the north west direction
	NorthWest
	// SouthEast represents the south east direction
	SouthEast
	// SouthWest represents the south west direction
	SouthWest
)

// Face converts the direction to a Face and returns it.
func (d Direction) Face() Face {
	return Face(d + 2)
}

// Opposite returns Direction opposite to the current one.
func (d Direction) Opposite() Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case West:
		return East
	case East:
		return West
	case NorthEast:
		return SouthWest
	case NorthWest:
		return SouthEast
	case SouthEast:
		return NorthWest
	case SouthWest:
		return NorthEast
	}
	panic("invalid direction")
}

// RotateRight rotates the direction 90 degrees to the right horizontally and returns the new direction.
func (d Direction) RotateRight() Direction {
	switch d {
	case North:
		return East
	case East:
		return South
	case South:
		return West
	case West:
		return North
	case NorthEast:
		return SouthEast
	case NorthWest:
		return NorthEast
	case SouthEast:
		return SouthWest
	case SouthWest:
		return NorthWest
	}
	panic("invalid direction")
}

// RotateLeft rotates the direction 90 degrees to the left horizontally and returns the new direction.
func (d Direction) RotateLeft() Direction {
	switch d {
	case North:
		return West
	case East:
		return North
	case South:
		return East
	case West:
		return South
	case NorthEast:
		return NorthWest
	case NorthWest:
		return SouthWest
	case SouthEast:
		return NorthEast
	case SouthWest:
		return SouthEast
	}
	panic("invalid direction")
}

// String returns the Direction as a string.
func (d Direction) String() string {
	switch d {
	case North:
		return "north"
	case East:
		return "east"
	case South:
		return "south"
	case West:
		return "west"
	case NorthEast:
		return "northeast"
	case NorthWest:
		return "northwest"
	case SouthEast:
		return "southeast"
	case SouthWest:
		return "southwest"
	}
	panic("invalid direction")
}

var directions = [...]Direction{North, East, South, West}

// Directions returns a list of all directions, going from North to West.
func Directions() []Direction {
	return directions[:]
}

func AllDirections() []Direction {
	return []Direction{North, East, South, West, NorthEast, NorthWest, SouthEast, SouthWest}
}
