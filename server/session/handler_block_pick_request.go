package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/stcraft/dragonfly/server/block/cube"
)

// BlockPickRequestHandler handles the BlockPickRequest packet.
type BlockPickRequestHandler struct{}

// Handle ...
func (b BlockPickRequestHandler) Handle(p packet.Packet, s *Session) error {
	pk := p.(*packet.BlockPickRequest)
	s.c.PickBlock(cube.Pos{int(pk.Position.X()), int(pk.Position.Y()), int(pk.Position.Z())})
	return nil
}
