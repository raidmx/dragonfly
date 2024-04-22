package block

import (
	"math/rand"

	"github.com/STCraft/dragonfly/server/item"
)

// NetherGoldOre is a variant of gold ore found exclusively in The Nether.
type NetherGoldOre struct {
	SolidModel
}

// BreakInfo ...
func (n NetherGoldOre) BreakInfo() BreakInfo {
	return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, silkTouchDrop(item.NewStack(item.GoldNugget{}, rand.Intn(4)+2), item.NewStack(n, 1))).withXPDropRange(0, 1)
}

// SmeltInfo ...
func (NetherGoldOre) SmeltInfo() item.SmeltInfo {
	return newOreSmeltInfo(item.NewStack(item.GoldIngot{}, 1), 1)
}

// EncodeItem ...
func (NetherGoldOre) EncodeItem() (name string, meta int16) {
	return "minecraft:nether_gold_ore", 0
}

// EncodeBlock ...
func (NetherGoldOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:nether_gold_ore", nil
}
