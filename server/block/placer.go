package block

import (
	"github.com/stcraft/dragonfly/server/block/cube"
	"github.com/stcraft/dragonfly/server/item"
	"github.com/stcraft/dragonfly/server/world"
)

// Placer represents an entity that is able to place a block at a specific position in the world.
type Placer interface {
	item.User
	PlaceBlock(pos cube.Pos, b world.Block, ctx *item.UseContext)
}
