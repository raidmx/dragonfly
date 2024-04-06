package block

import (
	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/world"
)

// Placer represents an entity that is able to place a block at a specific position in the world.
type Placer interface {
	item.User
	PlaceBlock(pos cube.Pos, b world.Block, ctx *item.UseContext)
}
