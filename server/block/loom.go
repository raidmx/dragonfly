package block

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
)

// Loom is a block used to apply patterns on banners. It is also used as a shepherd's job site block that is found in
// villages.
type Loom struct {
	solid
	bass

	// Facing is the direction the loom is facing.
	Facing cube.Direction
}

// FuelInfo ...
func (Loom) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

// BreakInfo ...
func (l Loom) BreakInfo() BreakInfo {
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, oneOf(l))
}

// Activate ...
func (Loom) Activate(pos cube.Pos, _ cube.Face, _ *world.World, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos)
		return true
	}
	return false
}

// UseOnBlock ...
func (l Loom) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, l)
	if !used {
		return
	}
	l.Facing = user.Rotation().Direction().Opposite()
	place(w, pos, l, user, ctx)
	return placed(ctx)
}

// EncodeItem ...
func (Loom) EncodeItem() (name string, meta int16) {
	return "minecraft:loom", 0
}

// EncodeBlock ...
func (l Loom) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:loom", map[string]interface{}{"direction": int32(horizontalDirection(l.Facing))}
}

// allLooms ...
func allLooms() (looms []world.Block) {
	for _, d := range cube.Directions() {
		looms = append(looms, Loom{Facing: d})
	}
	return
}
