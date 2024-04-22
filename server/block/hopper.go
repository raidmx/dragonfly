package block

import (
	"sync"

	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/internal/nbtconv"
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/item/inventory"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Hopper is a low-capacity storage block that can be used to collect item
// entities directly above it, as well as to transfer items into and out of
// other containers. A hopper can be locked with redstone power to stop it
// from moving items into or out of itself.
type Hopper struct {
	solid

	CustomName string
	Facing     cube.Face
	Toggled    bool

	inventory *inventory.Inventory
	viewerMu  *sync.RWMutex
	viewers   map[ContainerViewer]struct{}
}

// NewHopper creates a new initialised hopper. The inventory of the hopper
// is properly initialised. You still need to set the hopper's facing direction
// etc.
func NewHopper() Hopper {
	m := &sync.RWMutex{}
	v := map[ContainerViewer]struct{}{}

	return Hopper{
		inventory: inventory.New(5, func(slot int, _, item item.Stack) {
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
func (h Hopper) EncodeItem() (name string, meta int16) {
	return "minecraft:hopper", 0
}

// EncodeBlock ...
func (h Hopper) EncodeBlock() (string, map[string]any) {
	var facing int32

	switch h.Facing {
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

	return "minecraft:hopper", map[string]any{
		"facing_direction": facing,
		"toggle_bit":       h.Toggled,
	}
}

// Inventory returns the inventory of the hopper. The size of the inventory will
// always be 5
func (h Hopper) Inventory() *inventory.Inventory {
	return h.inventory
}

// AddViewer adds a viewer to the hopper, so that it is updated whenever the inventory of the hopper is changed.
func (h Hopper) AddViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	h.viewerMu.Lock()
	defer h.viewerMu.Unlock()
	h.viewers[v] = struct{}{}
}

// RemoveViewer removes a viewer from the hopper, so that slot updates in the inventory are no longer sent to
// it.
func (h Hopper) RemoveViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	h.viewerMu.Lock()
	defer h.viewerMu.Unlock()
	if len(h.viewers) == 0 {
		return
	}
	delete(h.viewers, v)
}

// Activate ...
func (h Hopper) Activate(pos cube.Pos, _ cube.Face, w *world.World, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		if d, ok := w.Block(pos.Side(cube.FaceUp)).(LightDiffuser); ok && d.LightDiffusionLevel() <= 2 {
			opener.OpenBlockContainer(pos)
		}
		return true
	}
	return false
}

// UseOnBlock ...
func (h Hopper) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, h)
	if !used {
		return
	}

	//noinspection GoAssignmentToReceiver
	h = NewHopper()
	h.Facing = user.Rotation().Direction().Face()

	place(w, pos, h, user, ctx)
	return placed(ctx)
}

// DecodeNBT ...
func (h Hopper) DecodeNBT(data map[string]any) any {
	facing := h.Facing
	//noinspection GoAssignmentToReceiver
	h = NewHopper()

	h.Facing = facing
	h.CustomName = nbtconv.String(data, "CustomName")
	nbtconv.InvFromNBT(h.inventory, nbtconv.Slice(data, "Items"))

	return h
}

// EncodeNBT ...
func (h Hopper) EncodeNBT() map[string]any {
	if h.inventory == nil {
		facing, customName := h.Facing, h.CustomName
		//noinspection GoAssignmentToReceiver
		h = NewHopper()
		h.Facing, h.CustomName = facing, customName
	}
	m := map[string]any{
		"Items": nbtconv.InvToNBT(h.inventory),
		"id":    "Hopper",
	}
	if h.CustomName != "" {
		m["CustomName"] = h.CustomName
	}
	return m
}

// allHoppers ...
func allHoppers() (hoppers []world.Block) {
	for _, face := range cube.Faces() {
		hoppers = append(hoppers, Hopper{Facing: face, Toggled: false})
		hoppers = append(hoppers, Hopper{Facing: face, Toggled: true})
	}
	return
}
