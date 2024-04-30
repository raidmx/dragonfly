package enchantment

import (
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
)

// Protection is an armour enchantment which increases the damage reduction.
type Protection struct{}

// Name ...
func (Protection) Name() string {
	return "Protection"
}

// MaxLevel ...
func (Protection) MaxLevel() int {
	return 4
}

// Cost ...
func (Protection) Cost(level int) (int, int) {
	min := 1 + (level-1)*11
	return min, min + 11
}

// Rarity ...
func (Protection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

// Modifier returns the base protection modifier for the enchantment.
func (Protection) Modifier() float64 {
	return 0.04
}

// CompatibleWithEnchantment ...
func (Protection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	_, blastProtection := t.(BlastProtection)
	_, fireProtection := t.(FireProtection)
	_, projectileProtection := t.(ProjectileProtection)
	return !blastProtection && !fireProtection && !projectileProtection
}

// CompatibleWithItem ...
func (Protection) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}
