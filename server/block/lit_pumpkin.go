package block

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
)

// LitPumpkin is a decorative light emitting block crafted with a Carved Pumpkin & Torch
type LitPumpkin struct {
	solid

	// Facing is the direction the pumpkin is facing.
	Facing cube.Direction
}

// LightEmissionLevel ...
func (l LitPumpkin) LightEmissionLevel() uint8 {
	return 15
}

// UseOnBlock ...
func (l LitPumpkin) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, l)
	if !used {
		return
	}
	l.Facing = user.Rotation().Direction().Opposite()

	place(w, pos, l, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (l LitPumpkin) BreakInfo() BreakInfo {
	return newBreakInfo(1, alwaysHarvestable, axeEffective, oneOf(l))
}

// EncodeItem ...
func (l LitPumpkin) EncodeItem() (name string, meta int16) {
	return "minecraft:lit_pumpkin", 0
}

// EncodeBlock ...
func (l LitPumpkin) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:lit_pumpkin", map[string]any{"minecraft:cardinal_direction": l.Facing.String()}
}

func allLitPumpkins() (pumpkins []world.Block) {
	for i := cube.Direction(0); i <= 3; i++ {
		pumpkins = append(pumpkins, LitPumpkin{Facing: i})
	}
	return
}
