package enchantment

import (
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
)

// FireProtection is an armour enchantment that decreases fire damage.
type FireProtection struct{}

// Name ...
func (FireProtection) Name() string {
	return "Fire Protection"
}

// MaxLevel ...
func (FireProtection) MaxLevel() int {
	return 4
}

// Cost ...
func (FireProtection) Cost(level int) (int, int) {
	min := 10 + (level-1)*8
	return min, min + 8
}

// Rarity ...
func (FireProtection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityUncommon
}

// Modifier returns the base protection modifier for the enchantment.
func (FireProtection) Modifier() float64 {
	return 0.08
}

// CompatibleWithEnchantment ...
func (FireProtection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	_, blastProtection := t.(BlastProtection)
	_, projectileProtection := t.(ProjectileProtection)
	_, protection := t.(Protection)
	return !blastProtection && !projectileProtection && !protection
}

// CompatibleWithItem ...
func (FireProtection) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}
