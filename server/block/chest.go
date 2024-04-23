package block

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/internal/nbtconv"
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/item/inventory"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/STCraft/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	ChestTypeSingle int = 27
	ChestTypeDouble int = 54
)

// Chest is a container block which may be used to store items. Chests may also be paired to create a bigger
// single container.
type Chest struct {
	chest
	transparent
	bass
	sourceWaterDisplacer

	// CustomName is the custom name of the chest. This name is displayed when the chest is opened, and may
	// include colour codes.
	CustomName string
	// Facing is the direction that the chest is facing.
	Facing cube.Direction

	// Paired is true if the chest is Paired
	Paired bool
	// PairX is the x coordinate of the pair
	PairX int32
	// PairZ is the z coordinate of the pair
	PairZ int32

	inventory *inventory.Inventory
	viewerMu  *sync.RWMutex
	viewers   map[ContainerViewer]struct{}
}

// NewChest creates a new initialised chest of the provided type. The inventory
// is properly initialised. You still need to set the correct direction facing of
// the chest etc.
func NewChest(chestType int) Chest {
	m := &sync.RWMutex{}
	v := map[ContainerViewer]struct{}{}

	return Chest{
		inventory: inventory.New(chestType, func(slot int, _, item item.Stack) {
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

// Inventory returns the inventory of the chest. The size of the inventory will be 27 or 54, depending on
// whether the chest is single or double.
func (c Chest) Inventory() *inventory.Inventory {
	return c.inventory
}

// WithName returns the chest after applying a specific name to the block.
func (c Chest) WithName(a ...any) world.Item {
	c.CustomName = strings.TrimSuffix(fmt.Sprintln(a...), "\n")
	return c
}

// SideClosed ...
func (Chest) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// LeftPair returns whether the chest at the specified position
// is a paired chest on the left side
func (c Chest) LeftPair(pos cube.Pos) bool {
	if pos.X() == int(c.PairX) {
		if c.Facing == cube.East {
			return pos.Z() > int(c.PairZ)
		}

		if c.Facing == cube.West {
			return pos.Z() < int(c.PairZ)
		}
	}

	if pos.Z() == int(c.PairZ) {
		if c.Facing == cube.North {
			return pos.X() > int(c.PairX)
		}

		if c.Facing == cube.South {
			return pos.X() < int(c.PairX)
		}
	}

	panic("this should never happen")
}

// RightPair returns whether the chest at the specified position
// is a paired chest on the right side
func (c Chest) RightPair(pos cube.Pos) bool {
	return !c.LeftPair(pos)
}

// Pair returns the cube.Pos of the chest paired with this chest
func (c Chest) Pair(pos cube.Pos) cube.Pos {
	return cube.Pos{int(c.PairX), pos.Y(), int(c.PairZ)}
}

// open opens the chest, displaying the animation and playing a sound.
func (c Chest) open(w *world.World, pos cube.Pos) {
	for _, v := range w.Viewers(pos.Vec3()) {
		v.ViewBlockAction(pos, OpenAction{})

		if c.Paired {
			v.ViewBlockAction(c.Pair(pos), OpenAction{})
		}
	}
	w.PlaySound(pos.Vec3Centre(), sound.ChestOpen{})
}

// close closes the chest, displaying the animation and playing a sound.
func (c Chest) close(w *world.World, pos cube.Pos) {
	for _, v := range w.Viewers(pos.Vec3()) {
		v.ViewBlockAction(pos, CloseAction{})

		if c.Paired {
			v.ViewBlockAction(c.Pair(pos), CloseAction{})
		}
	}
	w.PlaySound(pos.Vec3Centre(), sound.ChestClose{})
}

// AddViewer adds a viewer to the chest, so that it is updated whenever the inventory of the chest is changed.
func (c Chest) AddViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	c.viewerMu.Lock()
	defer c.viewerMu.Unlock()
	if len(c.viewers) == 0 {
		c.open(w, pos)
	}
	c.viewers[v] = struct{}{}
}

// RemoveViewer removes a viewer from the chest, so that slot updates in the inventory are no longer sent to
// it.
func (c Chest) RemoveViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	c.viewerMu.Lock()
	defer c.viewerMu.Unlock()
	if len(c.viewers) == 0 {
		return
	}
	delete(c.viewers, v)
	if len(c.viewers) == 0 {
		c.close(w, pos)
	}
}

// Activate ...
func (c Chest) Activate(pos cube.Pos, _ cube.Face, w *world.World, u item.User, _ *item.UseContext) bool {
	if opener, ok := u.(ContainerOpener); ok {
		if d, ok := w.Block(pos.Side(cube.FaceUp)).(LightDiffuser); ok && d.LightDiffusionLevel() <= 2 {
			opener.OpenBlockContainer(pos)
		}
		return true
	}
	return false
}

// UseOnBlock ...
func (c Chest) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, c)
	if !used {
		return
	}
	//noinspection GoAssignmentToReceiver
	c = NewChest(ChestTypeSingle)
	c.Facing = user.Rotation().Direction().Opposite()

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (c Chest) NeighbourUpdateTick(pos, neighbour cube.Pos, w *world.World) {
	// Make sure to ignore the neighbour update ticks for same
	// block
	if pos == neighbour {
		return
	}

	b := w.Block(neighbour)

	// If a block was broken and it was the chest pair of this paired chest
	// then we must unpair the chests
	if _, ok := b.(Air); ok && c.Paired {
		n := c.Pair(pos)

		// The paired chest block got broken. We should unpair the chests
		// now.
		if n == neighbour {
			c.Paired = false
			w.SetBlock(pos, c, nil)
		}

		return
	}

	// Check if the block that got placed is a chest
	pair, ok := b.(Chest)

	// It means some other block got placed, we must return
	if !ok {
		return
	}

	// We must ensure that the two chests we are trying to pair must
	// be facing in the same direction
	if c.Facing != pair.Facing {
		return
	}

	// If either of the chests are already paired with each other or some
	// other chest then we must return
	if c.Paired || pair.Paired {
		return
	}

	// Merge the inventory of both the chests into a single large inventory
	// of a double chest
	inv := c.inventory
	facing := c.Facing

	//noinspection GoAssignmentToReceiver
	c = NewChest(ChestTypeDouble)
	c.Facing = facing

	// Add the items from the original chest inventory
	for _, it := range inv.Clear() {
		c.inventory.AddItem(it)
	}

	pair.Paired = true
	pair.PairX = int32(pos.X())
	pair.PairZ = int32(pos.Z())
	pair.inventory = c.inventory

	c.Paired = true
	c.PairX = int32(neighbour.X())
	c.PairZ = int32(neighbour.Z())

	w.SetBlock(pos, c, nil)
	w.SetBlock(neighbour, pair, nil)
}

// BreakInfo ...
func (c Chest) BreakInfo() BreakInfo {
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, oneOf(c))
}

// FuelInfo ...
func (Chest) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

// FlammabilityInfo ...
func (c Chest) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, true)
}

// DecodeNBT ...
func (c Chest) DecodeNBT(data map[string]any) any {
	facing := c.Facing
	pairx := data["pairx"]
	pairz, ok := data["pairz"]

	if ok {
		//noinspection GoAssignmentToReceiver
		c = NewChest(ChestTypeDouble)

		c.Paired = true
		c.PairX = pairx.(int32)
		c.PairZ = pairz.(int32)
	} else {
		//noinspection GoAssignmentToReceiver
		c = NewChest(ChestTypeSingle)
	}

	c.Facing = facing
	c.CustomName = nbtconv.String(data, "CustomName")
	nbtconv.InvFromNBT(c.inventory, nbtconv.Slice(data, "Items"))

	return c
}

// EncodeNBT ...
func (c Chest) EncodeNBT() map[string]any {
	if c.inventory == nil {
		facing, customName := c.Facing, c.CustomName
		if c.Paired {
			//noinspection GoAssignmentToReceiver
			c = NewChest(ChestTypeDouble)
		} else {
			//noinspection GoAssignmentToReceiver
			c = NewChest(ChestTypeSingle)
		}
		c.Facing, c.CustomName = facing, customName
	}
	m := map[string]any{
		"Items": nbtconv.InvToNBT(c.inventory),
		"id":    "Chest",
	}
	if c.CustomName != "" {
		m["CustomName"] = c.CustomName
	}
	if c.Paired {
		m["pairx"] = c.PairX
		m["pairz"] = c.PairZ
	}
	return m
}

// EncodeItem ...
func (Chest) EncodeItem() (name string, meta int16) {
	return "minecraft:chest", 0
}

// EncodeBlock ...
func (c Chest) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:chest", map[string]any{"minecraft:cardinal_direction": c.Facing.String()}
}

// allChests ...
func allChests() (chests []world.Block) {
	for _, direction := range cube.Directions() {
		chests = append(chests, Chest{Facing: direction})
	}
	return
}
