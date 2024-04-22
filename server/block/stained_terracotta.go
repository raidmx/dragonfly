package block

import (
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
)

// StainedTerracotta is a block formed from clay, with a hardness and blast resistance comparable to stone. In contrast
// to Terracotta, t can be coloured in the same 16 colours that wool can be dyed, but more dulled and earthen.
type StainedTerracotta struct {
	SolidModel
	bassDrum

	// Colour specifies the colour of the block.
	Colour item.Colour
}

// SoilFor ...
func (t StainedTerracotta) SoilFor(block world.Block) bool {
	_, ok := block.(DeadBush)
	return ok
}

// BreakInfo ...
func (t StainedTerracotta) BreakInfo() BreakInfo {
	return newBreakInfo(1.25, pickaxeHarvestable, pickaxeEffective, oneOf(t)).withBlastResistance(21)
}

// SmeltInfo ...
func (t StainedTerracotta) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(GlazedTerracotta{Colour: t.Colour}, 1), 0.1)
}

// EncodeItem ...
func (t StainedTerracotta) EncodeItem() (name string, meta int16) {
	return "minecraft:" + t.Colour.String() + "_terracotta", 0
}

// EncodeBlock ...
func (t StainedTerracotta) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + t.Colour.String() + "_terracotta", nil
}

// allStainedTerracotta returns stained terracotta blocks with all possible colours.
func allStainedTerracotta() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, StainedTerracotta{Colour: c})
	}
	return b
}
