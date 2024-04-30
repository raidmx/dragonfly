package item

import (
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/world"
	"github.com/stcraft/dragonfly/server/world/sound"
)

// FireCharge is an item that can be used to place fire when used on a block, or shot from a dispenser to create a small
// fireball.
type FireCharge struct{}

// EncodeItem ...
func (f FireCharge) EncodeItem() (name string, meta int16) {
	return "minecraft:fire_charge", 0
}

// UseOnBlock ...
func (f FireCharge) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, u User, ctx *UseContext) bool {
	if l, ok := w.Block(pos).(ignitable); ok && l.Ignite(pos, w, u) {
		ctx.SubtractFromCount(1)
		w.PlaySound(pos.Vec3Centre(), sound.FireCharge{})
		return true
	} else if s := pos.Side(face); w.Block(s) == air() {
		ctx.SubtractFromCount(1)
		w.PlaySound(s.Vec3Centre(), sound.FireCharge{})
		w.SetBlock(s, fire(), nil)
		w.ScheduleBlockUpdate(s, time.Duration(30+rand.Intn(10))*time.Second/20)
		return true
	}
	return false
}
