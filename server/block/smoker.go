package block

import (
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/internal/nbtconv"
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
	"github.com/stcraft/dragonfly/server/world/sound"
)

// Smoker is a type of furnace that cooks food items, similar to a furnace, but twice as fast. It also serves as a
// butcher's job site block.
// The empty value of Smoker is not valid. It must be created using block.NewSmoker(cube.Face).
type Smoker struct {
	solid
	bassDrum
	*smelter

	// Facing is the direction the smoker is facing.
	Facing cube.Direction
	// Lit is true if the smoker is lit.
	Lit bool
}

// NewSmoker creates a new initialised smoker. The smelter is properly initialised.
func NewSmoker(face cube.Direction) Smoker {
	return Smoker{
		Facing:  face,
		smelter: newSmelter(),
	}
}

// Tick is called to check if the smoker should update and start or stop smelting.
func (s Smoker) Tick(_ int64, pos cube.Pos, w *world.World) {
	if s.Lit && rand.Float64() <= 0.016 { // Every three or so seconds.
		w.PlaySound(pos.Vec3Centre(), sound.SmokerCrackle{})
	}
	if lit := s.smelter.tickSmelting(time.Second*5, time.Millisecond*200, s.Lit, func(i item.SmeltInfo) bool {
		return i.Food
	}); s.Lit != lit {
		s.Lit = lit
		w.SetBlock(pos, s, nil)
	}
}

// EncodeItem ...
func (s Smoker) EncodeItem() (name string, meta int16) {
	return "minecraft:smoker", 0
}

// EncodeBlock ...
func (s Smoker) EncodeBlock() (name string, properties map[string]interface{}) {
	if s.Lit {
		return "minecraft:lit_smoker", map[string]interface{}{"minecraft:cardinal_direction": s.Facing.String()}
	}
	return "minecraft:smoker", map[string]interface{}{"minecraft:cardinal_direction": s.Facing.String()}
}

// UseOnBlock ...
func (s Smoker) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, s)
	if !used {
		return false
	}

	place(w, pos, NewSmoker(user.Rotation().Direction().Opposite()), user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (s Smoker) BreakInfo() BreakInfo {
	xp := s.Experience()
	return newBreakInfo(3.5, alwaysHarvestable, pickaxeEffective, oneOf(s)).withXPDropRange(xp, xp)
}

// Activate ...
func (s Smoker) Activate(pos cube.Pos, _ cube.Face, _ *world.World, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos)
		return true
	}
	return false
}

// EncodeNBT ...
func (s Smoker) EncodeNBT() map[string]interface{} {
	if s.smelter == nil {
		//noinspection GoAssignmentToReceiver
		s = NewSmoker(s.Facing)
	}
	remaining, maximum, cook := s.Durations()
	return map[string]interface{}{
		"BurnTime":     int16(remaining.Milliseconds() / 50),
		"CookTime":     int16(cook.Milliseconds() / 50),
		"BurnDuration": int16(maximum.Milliseconds() / 50),
		"StoredXPInt":  int16(s.Experience()),
		"Items":        nbtconv.InvToNBT(s.Inventory()),
		"id":           "Smoker",
	}
}

// DecodeNBT ...
func (s Smoker) DecodeNBT(data map[string]interface{}) interface{} {
	remaining := nbtconv.TickDuration[int16](data, "BurnTime")
	maximum := nbtconv.TickDuration[int16](data, "BurnDuration")
	cook := nbtconv.TickDuration[int16](data, "CookTime")

	xp := int(nbtconv.Int16(data, "StoredXPInt"))
	lit := s.Lit

	//noinspection GoAssignmentToReceiver
	s = NewSmoker(s.Facing)
	s.Lit = lit
	s.setExperience(xp)
	s.setDurations(remaining, maximum, cook)
	nbtconv.InvFromNBT(s.Inventory(), nbtconv.Slice(data, "Items"))
	return s
}

// allSmokers ...
func allSmokers() (smokers []world.Block) {
	for _, face := range cube.Directions() {
		smokers = append(smokers, Smoker{Facing: face})
		smokers = append(smokers, Smoker{Facing: face, Lit: true})
	}
	return
}
