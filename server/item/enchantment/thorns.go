package enchantment

import (
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
)

// Thorns is an enchantment that inflicts damage on attackers.
type Thorns struct{}

// Name ...
func (Thorns) Name() string {
	return "Thorns"
}

// MaxLevel ...
func (Thorns) MaxLevel() int {
	return 3
}

// Cost ...
func (Thorns) Cost(level int) (int, int) {
	min := 10 + 20*(level-1)
	return min, min + 50
}

// Rarity ...
func (Thorns) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

// CompatibleWithEnchantment ...
func (Thorns) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (Thorns) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}

// ThornsDamageSource is used for damage caused by thorns.
type ThornsDamageSource struct {
	// Owner is the owner of the armour with the Thorns enchantment.
	Owner world.Entity
}

func (ThornsDamageSource) ReducedByResistance() bool { return true }
func (ThornsDamageSource) ReducedByArmour() bool     { return false }
func (ThornsDamageSource) Fire() bool                { return false }
