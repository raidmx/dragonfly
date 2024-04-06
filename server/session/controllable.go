package session

import (
	"github.com/STCraft/dragonfly/server/block/cube"
	"github.com/STCraft/dragonfly/server/cmd"
	"github.com/STCraft/dragonfly/server/entity/effect"
	"github.com/STCraft/dragonfly/server/item"
	"github.com/STCraft/dragonfly/server/item/inventory"
	"github.com/STCraft/dragonfly/server/player/chat"
	"github.com/STCraft/dragonfly/server/player/form"
	"github.com/STCraft/dragonfly/server/player/skin"
	"github.com/STCraft/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

// Controllable represents an entity that may be controlled by a Session. Generally, Controllable is
// implemented in the form of a Player.
// Methods in Controllable will be added as Session needs them in order to handle packets.
type Controllable interface {
	Name() string
	world.Entity
	item.User
	form.Submitter
	cmd.Source
	chat.Subscriber

	Locale() language.Tag

	SetHeldItems(right, left item.Stack)

	Move(deltaPos mgl64.Vec3, deltaYaw, deltaPitch float64)
	Speed() float64

	Chat(msg ...any)
	ExecuteCommand(commandLine string)
	GameMode() world.GameMode
	SetGameMode(mode world.GameMode)
	Effects() []effect.Effect

	UseItem()
	ReleaseItem()
	UseItemOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)
	UseItemOnEntity(e world.Entity) bool
	BreakBlock(pos cube.Pos)
	PickBlock(pos cube.Pos)
	AttackEntity(e world.Entity) bool
	Drop(s item.Stack) (n int)
	SwingArm()
	PunchAir()

	ExperienceLevel() int
	SetExperienceLevel(level int)

	EnchantmentSeed() int64
	ResetEnchantmentSeed()

	Respawn()
	Dead() bool

	StartSneaking()
	Sneaking() bool
	StopSneaking()
	StartSprinting()
	Sprinting() bool
	StopSprinting()
	StartSwimming()
	Swimming() bool
	StopSwimming()
	StartFlying()
	Flying() bool
	StopFlying()
	StartGliding()
	Gliding() bool
	StopGliding()
	Jump()

	StartBreaking(pos cube.Pos, face cube.Face)
	ContinueBreaking(face cube.Face)
	FinishBreaking()
	AbortBreaking()

	Exhaust(points float64)

	OpenSign(pos cube.Pos, frontSide bool)
	EditSign(pos cube.Pos, frontText, backText string) error
	TurnLecternPage(pos cube.Pos, page int) error

	EnderChestInventory() *inventory.Inventory

	// UUID returns the UUID of the controllable. It must be unique for all controllable entities present in
	// the server.
	UUID() uuid.UUID
	// XUID returns the XBOX Live User ID of the controllable. Every controllable must have one of these if
	// they are authenticated via XBOX Live, as they must be connected to an XBOX Live account.
	XUID() string
	// Skin returns the skin of the controllable. Each controllable must have a skin, as it defines how the
	// entity looks in the world.
	Skin() skin.Skin
	SetSkin(skin.Skin)
}
