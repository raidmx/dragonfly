package block

import (
	"sync"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/internal/nbtconv"
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/item/inventory"
	"github.com/stcraft/dragonfly/server/world"
)

// Dispenser is a low-capacity storage block that can fire projectiles, use
// certain items or tools or place certain blocks, fluids or entities when
// given a redstone signal. Items that do not have unique dispenser functions
// are instead ejected into the world.
type Dispenser struct {
	solid

	CustomName string
	Facing     cube.Face
	Triggered  bool

	inventory *inventory.Inventory
	viewerMu  *sync.RWMutex
	viewers   map[ContainerViewer]struct{}
}

// NewDispenser creates a new initialised dispenser. The inventory of the dispenser
// is properly initialised. You still need to set the dispenser's facing direction
// etc.
func NewDispenser() Dispenser {
	m := &sync.RWMutex{}
	v := map[ContainerViewer]struct{}{}

	return Dispenser{
		inventory: inventory.New(9, func(slot int, _, item item.Stack) {
			m.RLock()
			defer m.RUnlock()
			for viewer := range v {
				viewer.ViewSlotChange(slot, item)
			}
		}),
		viewerMu: m,
		viewers:  v,
	}
}

// EncodeItem ...
func (d Dispenser) EncodeItem() (name string, meta int16) {
	return "minecraft:dispenser", 0
}

// EncodeBlock ...
func (d Dispenser) EncodeBlock() (string, map[string]any) {
	var facing int32

	switch d.Facing {
	case cube.FaceDown:
		facing = 0
	case cube.FaceUp:
		facing = 1
	case cube.FaceNorth:
		facing = 2
	case cube.FaceSouth:
		facing = 3
	case cube.FaceWest:
		facing = 4
	case cube.FaceEast:
		facing = 5
	}

	return "minecraft:dispenser", map[string]any{
		"facing_direction": facing,
		"triggered_bit":    d.Triggered,
	}
}

// Inventory returns the inventory of the dispenser. The size of the inventory will
// always be 5
func (d Dispenser) Inventory() *inventory.Inventory {
	return d.inventory
}

// AddViewer adds a viewer to the dispenser, so that it is updated whenever the inventory of the dispenser is changed.
func (d Dispenser) AddViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	d.viewerMu.Lock()
	defer d.viewerMu.Unlock()
	d.viewers[v] = struct{}{}
}

// RemoveViewer removes a viewer from the dispenser, so that slot updates in the inventory are no longer sent to
// it.
func (d Dispenser) RemoveViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	d.viewerMu.Lock()
	defer d.viewerMu.Unlock()
	if len(d.viewers) == 0 {
		return
	}
	delete(d.viewers, v)
}

// Activate ...
func (d Dispenser) Activate(pos cube.Pos, _ cube.Face, w *world.World, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		if d, ok := w.Block(pos.Side(cube.FaceUp)).(LightDiffuser); ok && d.LightDiffusionLevel() <= 2 {
			opener.OpenBlockContainer(pos)
		}
		return true
	}
	return false
}

// UseOnBlock ...
func (d Dispenser) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, d)
	if !used {
		return
	}

	//noinspection GoAssignmentToReceiver
	d = NewDispenser()
	d.Facing = user.Rotation().Opposite().Direction().Face()

	place(w, pos, d, user, ctx)
	return placed(ctx)
}

// DecodeNBT ...
func (d Dispenser) DecodeNBT(data map[string]any) any {
	facing := d.Facing
	//noinspection GoAssignmentToReceiver
	d = NewDispenser()

	d.Facing = facing
	d.CustomName = nbtconv.String(data, "CustomName")
	nbtconv.InvFromNBT(d.inventory, nbtconv.Slice(data, "Items"))

	return d
}

// EncodeNBT ...
func (d Dispenser) EncodeNBT() map[string]any {
	if d.inventory == nil {
		facing, customName := d.Facing, d.CustomName
		//noinspection GoAssignmentToReceiver
		d = NewDispenser()
		d.Facing, d.CustomName = facing, customName
	}
	m := map[string]any{
		"Items": nbtconv.InvToNBT(d.inventory),
		"id":    "Dispenser",
	}
	if d.CustomName != "" {
		m["CustomName"] = d.CustomName
	}
	return m
}

// allDispensers ...
func allDispensers() (Dispensers []world.Block) {
	for _, face := range cube.Faces() {
		Dispensers = append(Dispensers, Dispenser{Facing: face, Triggered: false})
		Dispensers = append(Dispensers, Dispenser{Facing: face, Triggered: true})
	}
	return
}
