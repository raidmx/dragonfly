package block

import (
	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/world"
)

// Hopper is a low-capacity storage block that can be used to collect item
// entities directly above it, as well as to transfer items into and out of
// other containers. A hopper can be locked with redstone power to stop it
// from moving items into or out of itself.
type Hopper struct {
	solid

	Facing  cube.Face
	Toggled bool
}

// EncodeItem ...
func (h Hopper) EncodeItem() (name string, meta int16) {
	return "minecraft:hopper", 0
}

// EncodeBlock ...
func (h Hopper) EncodeBlock() (string, map[string]any) {
	var facing int32

	switch h.Facing {
	case cube.FaceDown:
		facing = 0
	case cube.FaceUp:
		facing = 1
	case cube.FaceNorth:
		facing = 2
	case cube.FaceSouth:
		facing = 3
	case cube.FaceWest:
		facing = 4
	case cube.FaceEast:
		facing = 5
	}

	return "minecraft:hopper", map[string]any{
		"facing_direction": facing,
		"toggle_bit":       h.Toggled,
	}
}

// allHoppers ...
func allHoppers() (hoppers []world.Block) {
	for _, face := range cube.Faces() {
		hoppers = append(hoppers, Hopper{Facing: face, Toggled: false})
		hoppers = append(hoppers, Hopper{Facing: face, Toggled: true})
	}
	return
}
