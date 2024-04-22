package block

// StructureBlock ...
type StructureBlock struct {
	solid
}

// EncodeItem ...
func (s StructureBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:structure_block", 0
}

// EncodeBlock ...
func (s StructureBlock) EncodeBlock() (string, map[string]any) {
	return "minecraft:structure_block", map[string]any{
		"structure_block_type": "corner",
	}
}
