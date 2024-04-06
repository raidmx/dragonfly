package item

import (
	"math/rand"
	"time"

	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/STCraft/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// Firework is an item (and entity) used for creating decorative explosions, boosting when flying with elytra, and
// loading into a crossbow as ammunition.
type Firework struct {
	// Duration is the flight duration of the firework.
	Duration time.Duration
	// Explosions is the list of explosions the firework should create when launched.
	Explosions []FireworkExplosion
}

// Use ...
func (f Firework) Use(w *world.World, user User, ctx *UseContext) bool {
	if g, ok := user.(interface {
		Gliding() bool
	}); !ok || !g.Gliding() {
		return false
	}

	pos := user.Position()

	w.PlaySound(pos, sound.FireworkLaunch{})
	create := w.EntityRegistry().Config().Firework
	w.AddEntity(create(pos, user.Rotation(), true, f, user))

	ctx.SubtractFromCount(1)
	return true
}

// UseOnBlock ...
func (f Firework) UseOnBlock(blockPos cube.Pos, _ cube.Face, clickPos mgl64.Vec3, w *world.World, user User, ctx *UseContext) bool {
	pos := blockPos.Vec3().Add(clickPos)
	create := w.EntityRegistry().Config().Firework
	w.AddEntity(create(pos, cube.Rotation{rand.Float64() * 360, 90}, false, f, user))
	w.PlaySound(pos, sound.FireworkLaunch{})

	ctx.SubtractFromCount(1)
	return true
}

// EncodeNBT ...
func (f Firework) EncodeNBT() map[string]any {
	explosions := make([]any, 0, len(f.Explosions))
	for _, explosion := range f.Explosions {
		explosions = append(explosions, explosion.EncodeNBT())
	}
	return map[string]any{"Fireworks": map[string]any{
		"Explosions": explosions,
		"Flight":     uint8((f.Duration/10 - time.Millisecond*50).Milliseconds() / 50),
	}}
}

// DecodeNBT ...
func (f Firework) DecodeNBT(data map[string]any) any {
	if fireworks, ok := data["Fireworks"].(map[string]any); ok {
		if explosions, ok := fireworks["Explosions"].([]any); ok {
			f.Explosions = make([]FireworkExplosion, len(explosions))
			for i, explosion := range f.Explosions {
				f.Explosions[i] = explosion.DecodeNBT(explosions[i].(map[string]any)).(FireworkExplosion)
			}
		}
		if durationTicks, ok := fireworks["Flight"].(uint8); ok {
			f.Duration = (time.Duration(durationTicks)*time.Millisecond*50 + time.Millisecond*50) * 10
		}
	}
	return f
}

// RandomisedDuration returns the randomised flight duration of the firework.
func (f Firework) RandomisedDuration() time.Duration {
	return f.Duration + time.Duration(rand.Intn(int(time.Millisecond*600)))
}

// EncodeItem ...
func (Firework) EncodeItem() (name string, meta int16) {
	return "minecraft:firework_rocket", 0
}
