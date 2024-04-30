package block

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/block/model"
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
	"github.com/stcraft/dragonfly/server/world/particle"
)

// Cactus is a plant block that generates naturally in dry areas and causes damage.
type Cactus struct {
	transparent

	// Age is the growth state of cactus. Values range from 0 to 15.
	Age int
}

// UseOnBlock handles making sure the neighbouring blocks are air.
func (c Cactus) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, c)
	if !used {
		return false
	}
	if !c.canGrowHere(pos, w, true) {
		return false
	}

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (c Cactus) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !c.canGrowHere(pos, w, true) {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: c})
		dropItem(w, item.NewStack(c, 1), pos.Vec3Centre())
	}
}

// RandomTick ...
func (c Cactus) RandomTick(pos cube.Pos, w *world.World, _ *rand.Rand) {
	if c.Age < 15 {
		c.Age++
	} else if c.Age == 15 {
		c.Age = 0
		if c.canGrowHere(pos.Side(cube.FaceDown), w, false) {
			for y := 1; y < 3; y++ {
				if _, ok := w.Block(pos.Add(cube.Pos{0, y})).(Air); ok {
					w.SetBlock(pos.Add(cube.Pos{0, y}), Cactus{Age: 0}, nil)
					break
				} else if _, ok := w.Block(pos.Add(cube.Pos{0, y})).(Cactus); !ok {
					break
				}
			}
		}
	}
	w.SetBlock(pos, c, nil)
}

// canGrowHere implements logic to check if cactus can live/grow here.
func (c Cactus) canGrowHere(pos cube.Pos, w *world.World, recursive bool) bool {
	for _, face := range cube.HorizontalFaces() {
		if _, ok := w.Block(pos.Side(face)).(Air); !ok {
			return false
		}
	}
	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Cactus); ok && recursive {
		return c.canGrowHere(pos.Side(cube.FaceDown), w, recursive)
	}
	return supportsVegetation(c, w.Block(pos.Sub(cube.Pos{0, 1})))
}

// EntityInside ...
func (c Cactus) EntityInside(_ cube.Pos, _ *world.World, e world.Entity) {
	if l, ok := e.(livingEntity); ok && !l.AttackImmune() {
		l.Hurt(0.5, DamageSource{Block: c})
	}
}

// BreakInfo ...
func (c Cactus) BreakInfo() BreakInfo {
	return newBreakInfo(0.4, alwaysHarvestable, nothingEffective, oneOf(c))
}

// CompostChance ...
func (Cactus) CompostChance() float64 {
	return 0.5
}

// EncodeItem ...
func (c Cactus) EncodeItem() (name string, meta int16) {
	return "minecraft:cactus", 0
}

// EncodeBlock ...
func (c Cactus) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:cactus", map[string]any{"age": int32(c.Age)}
}

// Model ...
func (c Cactus) Model() world.BlockModel {
	return model.Cactus{}
}

// allCactus returns all possible states of a cactus block.
func allCactus() (b []world.Block) {
	for i := 0; i < 16; i++ {
		b = append(b, Cactus{Age: i})
	}
	return
}

// DamageSource is passed as world.DamageSource for damage caused by a block,
// such as a cactus or a falling anvil.
type DamageSource struct {
	// Block is the block that caused the damage.
	Block world.Block
}

func (DamageSource) ReducedByResistance() bool { return true }
func (DamageSource) ReducedByArmour() bool     { return true }
func (DamageSource) Fire() bool                { return false }
