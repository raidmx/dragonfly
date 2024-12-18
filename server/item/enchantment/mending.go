package enchantment

import (
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
)

// Mending is an enchantment that repairs the item when experience orbs are collected.
type Mending struct{}

// Name ...
func (Mending) Name() string {
	return "Mending"
}

// MaxLevel ...
func (Mending) MaxLevel() int {
	return 1
}

// Cost ...
func (Mending) Cost(level int) (int, int) {
	min := level * 25
	return min, min + 50
}

// Rarity ...
func (Mending) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

// Treasure ...
func (Mending) Treasure() bool {
	return true
}

// CompatibleWithEnchantment ...
func (Mending) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	_, infinity := t.(Infinity)
	return !infinity
}

// CompatibleWithItem ...
func (Mending) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Durable)
	return ok
}
