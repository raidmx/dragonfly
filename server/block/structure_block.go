package block

import "github.com/stcraft/dragonfly/server/internal/nbtconv"

// StructureBlock is used to generate structures manually. They can also
// be used to save and load structures, alongside structure void blocks.
type StructureBlock struct {
	solid

	X int32
	Y int32
	Z int32
}

// EncodeItem ...
func (s StructureBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:structure_block", 0
}

// EncodeBlock ...
func (s StructureBlock) EncodeBlock() (string, map[string]any) {
	return "minecraft:structure_block", map[string]any{
		"structure_block_type": "data",
	}
}

// DecodeNBT ...
func (s StructureBlock) DecodeNBT(data map[string]any) any {
	s.X = nbtconv.Int32(data, "x")
	s.Y = nbtconv.Int32(data, "y")
	s.Z = nbtconv.Int32(data, "z")

	return s
}

// EncodeNBT ...
func (s StructureBlock) EncodeNBT() map[string]any {
	return map[string]any{
		"id":               "StructureBlock",
		"x":                s.X,
		"y":                s.Y,
		"z":                s.Z,
		"xStructureOffset": int32(0),
		"yStructureOffset": int32(2),
		"zStructureOffset": int32(0),
		"xStructureSize":   int32(16),
		"yStructureSize":   int32(256),
		"zStructureSize":   int32(16),
		"structureName":    "Chunk Border",
		"showBoundingBox":  byte(1),
		"seed":             int64(0),
		"rotation":         byte(0),
		"removeBlocks":     byte(0),
		"mirror":           byte(0),
		"isPowered":        byte(1),
		"isMovable":        byte(0),
		"integrity":        float32(1.0),
		"includePlayers":   byte(1),
		"includeEntities":  byte(0),
		"dataField":        "",
		"data":             int32(5),
	}
}
