package block

import (
	"time"

	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/internal/nbtconv"
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
	"github.com/stcraft/dragonfly/server/world/particle"
	"github.com/stcraft/dragonfly/server/world/sound"
)

// Note is a musical block that emits sounds when powered with redstone.
type Note struct {
	solid
	bass

	// Pitch is the current pitch the note block is set to. Value ranges from 0-24.
	Pitch int
}

// playNote ...
func (n Note) playNote(pos cube.Pos, w *world.World) {
	w.PlaySound(pos.Vec3(), sound.Note{Instrument: n.instrument(pos, w), Pitch: n.Pitch})
	w.AddParticle(pos.Vec3(), particle.Note{Instrument: n.Instrument(), Pitch: n.Pitch})
}

// updateInstrument ...
func (n Note) instrument(pos cube.Pos, w *world.World) sound.Instrument {
	if instrumentBlock, ok := w.Block(pos.Side(cube.FaceDown)).(interface {
		Instrument() sound.Instrument
	}); ok {
		return instrumentBlock.Instrument()
	}
	return sound.Piano()
}

// DecodeNBT ...
func (n Note) DecodeNBT(data map[string]any) any {
	n.Pitch = int(nbtconv.Uint8(data, "note"))
	return n
}

// EncodeNBT ...
func (n Note) EncodeNBT() map[string]any {
	return map[string]any{"note": byte(n.Pitch)}
}

// Activate ...
func (n Note) Activate(pos cube.Pos, _ cube.Face, w *world.World, _ item.User, _ *item.UseContext) bool {
	if _, ok := w.Block(pos.Side(cube.FaceUp)).(Air); !ok {
		return false
	}
	n.Pitch = (n.Pitch + 1) % 25
	n.playNote(pos, w)
	w.SetBlock(pos, n, &world.SetOpts{DisableBlockUpdates: true, DisableLiquidDisplacement: true})
	return true
}

// BreakInfo ...
func (n Note) BreakInfo() BreakInfo {
	return newBreakInfo(0.8, alwaysHarvestable, axeEffective, oneOf(n))
}

// FuelInfo ...
func (Note) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

// EncodeItem ...
func (n Note) EncodeItem() (name string, meta int16) {
	return "minecraft:noteblock", 0
}

// EncodeBlock ...
func (n Note) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:noteblock", nil
}
