package block

import (
	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Carpet is a colourful block that can be obtained by killing/shearing sheep, or crafted using four string.
type Carpet struct {
	carpet
	transparent
	sourceWaterDisplacer

	// Colour is the colour of the carpet.
	Colour item.Colour
}

// FlammabilityInfo ...
func (c Carpet) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(30, 20, true)
}

// SideClosed ...
func (Carpet) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// BreakInfo ...
func (c Carpet) BreakInfo() BreakInfo {
	return newBreakInfo(0.1, alwaysHarvestable, nothingEffective, oneOf(c))
}

// EncodeItem ...
func (c Carpet) EncodeItem() (name string, meta int16) {
	return "minecraft:" + c.Colour.String() + "_carpet", 0
}

// EncodeBlock ...
func (c Carpet) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + c.Colour.String() + "_carpet", nil
}

// HasLiquidDrops ...
func (Carpet) HasLiquidDrops() bool {
	return true
}

// NeighbourUpdateTick ...
func (c Carpet) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Air); ok {
		w.SetBlock(pos, nil, nil)
		dropItem(w, item.NewStack(c, 1), pos.Vec3Centre())
	}
}

// UseOnBlock handles not placing carpets on top of air blocks.
func (c Carpet) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, c)
	if !used {
		return
	}

	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Air); ok {
		return
	}

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// allCarpet ...
func allCarpet() (carpets []world.Block) {
	for _, c := range item.Colours() {
		carpets = append(carpets, Carpet{Colour: c})
	}
	return
}
