package block

import (
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
)

// StainedGlass is a decorative, fully transparent solid block that is dyed into a different colour.
type StainedGlass struct {
	transparent
	solid
	clicksAndSticks

	// Colour specifies the colour of the block.
	Colour item.Colour
}

// BreakInfo ...
func (g StainedGlass) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(g))
}

// EncodeItem ...
func (g StainedGlass) EncodeItem() (name string, meta int16) {
	return "minecraft:" + g.Colour.String() + "_stained_glass", 0
}

// EncodeBlock ...
func (g StainedGlass) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + g.Colour.String() + "_stained_glass", nil
}

// allStainedGlass returns stained-glass blocks with all possible colours.
func allStainedGlass() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, StainedGlass{Colour: c})
	}
	return b
}
