package enchantment

import (
	"time"

	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
)

// Flame turns your arrows into flaming arrows allowing you to set your targets on fire.
type Flame struct{}

// Name ...
func (Flame) Name() string {
	return "Flame"
}

// MaxLevel ...
func (Flame) MaxLevel() int {
	return 1
}

// Cost ...
func (Flame) Cost(int) (int, int) {
	return 20, 50
}

// Rarity ...
func (Flame) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}

// BurnDuration always returns a hundred seconds, no matter the level.
func (Flame) BurnDuration() time.Duration {
	return time.Second * 100
}

// CompatibleWithEnchantment ...
func (Flame) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (Flame) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Bow)
	return ok
}
