package effect

import (
	"image/color"

	"github.com/stcraft/dragonfly/server/world"
)

// Resistance is a lasting effect that reduces the damage taken from any sources except for void damage or
// custom damage.
type Resistance struct {
	nopLasting
}

// Multiplier returns a damage multiplier for the damage source passed.
func (Resistance) Multiplier(e world.DamageSource, lvl int) float64 {
	if !e.ReducedByResistance() {
		return 1
	}
	if v := 1 - 0.2*float64(lvl); v >= 0 {
		return v
	}
	return 0
}

// RGBA ...
func (Resistance) RGBA() color.RGBA {
	return color.RGBA{R: 0x99, G: 0x45, B: 0x3a, A: 0xff}
}
