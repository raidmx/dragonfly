package block

import (
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world/sound"
)

// Gold is a precious metal block crafted from 9 gold ingots.
type Gold struct {
	solid
}

// Instrument ...
func (g Gold) Instrument() sound.Instrument {
	return sound.Bell()
}

// BreakInfo ...
func (g Gold) BreakInfo() BreakInfo {
	return newBreakInfo(3, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierIron.HarvestLevel
	}, pickaxeEffective, oneOf(g)).withBlastResistance(30)
}

// PowersBeacon ...
func (Gold) PowersBeacon() bool {
	return true
}

// EncodeItem ...
func (Gold) EncodeItem() (name string, meta int16) {
	return "minecraft:gold_block", 0
}

// EncodeBlock ...
func (Gold) EncodeBlock() (string, map[string]any) {
	return "minecraft:gold_block", nil
}
