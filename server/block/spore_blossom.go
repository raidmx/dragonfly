package block

import (
	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/STCraft/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// SporeBlossom is a decorative block.
type SporeBlossom struct {
	EmptyModel
	transparent
}

// HasLiquidDrops ...
func (s SporeBlossom) HasLiquidDrops() bool {
	return true
}

// NeighbourUpdateTick ...
func (s SporeBlossom) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !w.Block(pos.Side(cube.FaceUp)).Model().FaceSolid(pos.Side(cube.FaceUp), cube.FaceDown, w) {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: s})
	}
}

// UseOnBlock ...
func (s SporeBlossom) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, s)
	if !used {
		return
	}
	if !w.Block(pos.Side(cube.FaceUp)).Model().FaceSolid(pos.Side(cube.FaceUp), cube.FaceDown, w) {
		return
	}

	place(w, pos, s, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (s SporeBlossom) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(s))
}

// FlammabilityInfo ...
func (SporeBlossom) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(15, 100, true)
}

// CompostChance ...
func (SporeBlossom) CompostChance() float64 {
	return 0.65
}

// EncodeItem ...
func (s SporeBlossom) EncodeItem() (name string, meta int16) {
	return "minecraft:spore_blossom", 0
}

// EncodeBlock ...
func (s SporeBlossom) EncodeBlock() (string, map[string]any) {
	return "minecraft:spore_blossom", nil
}
