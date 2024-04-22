package block

import (
	"math/rand"
	"time"

	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/STCraft/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// Coral is a non-solid block that comes in 5 variants.
type Coral struct {
	EmptyModel
	transparent
	bassDrum
	sourceWaterDisplacer

	// Type is the type of coral of the block.
	Type CoralType
	// Dead is whether the coral is dead.
	Dead bool
}

// UseOnBlock ...
func (c Coral) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, c)
	if !used {
		return false
	}
	if !w.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, w) {
		return false
	}
	if liquid, ok := w.Liquid(pos); ok {
		if water, ok := liquid.(Water); ok {
			if water.Depth != 8 {
				return false
			}
		}
	}

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// HasLiquidDrops ...
func (c Coral) HasLiquidDrops() bool {
	return false
}

// SideClosed ...
func (c Coral) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// NeighbourUpdateTick ...
func (c Coral) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !w.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, w) {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: c})
		return
	}
	if c.Dead {
		return
	}
	w.ScheduleBlockUpdate(pos, time.Second*5/2)
}

// ScheduledTick ...
func (c Coral) ScheduledTick(pos cube.Pos, w *world.World, _ *rand.Rand) {
	if c.Dead {
		return
	}

	adjacentWater := false
	pos.Neighbours(func(neighbour cube.Pos) {
		if liquid, ok := w.Liquid(neighbour); ok {
			if _, ok := liquid.(Water); ok {
				adjacentWater = true
			}
		}
	}, w.Range())
	if !adjacentWater {
		c.Dead = true
		w.SetBlock(pos, c, nil)
	}
}

// BreakInfo ...
func (c Coral) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(c))
}

// EncodeBlock ...
func (c Coral) EncodeBlock() (name string, properties map[string]any) {
	if c.Dead {
		return "minecraft:dead_" + c.Type.String() + "_coral", nil
	}
	return "minecraft:" + c.Type.String() + "_coral", nil
}

// EncodeItem ...
func (c Coral) EncodeItem() (name string, meta int16) {
	if c.Dead {
		return "minecraft:dead_" + c.Type.String() + "_coral", 0
	}
	return "minecraft:" + c.Type.String() + "_coral", 0
}

// allCoral returns a list of all coral block variants
func allCoral() (c []world.Block) {
	f := func(dead bool) {
		for _, t := range CoralTypes() {
			c = append(c, Coral{Type: t, Dead: dead})
		}
	}
	f(true)
	f(false)
	return
}
