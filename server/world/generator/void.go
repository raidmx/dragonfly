package generator

import (
	"github.com/STCraft/dragonfly/server/world"
	"github.com/STCraft/dragonfly/server/world/chunk"
)

// Flat is the flat generator of World. It generates flat worlds (like those in vanilla) with no other
// decoration. It may be constructed by calling NewFlat.
type Void struct {
	// biome is the encoded biome that the generator should use.
	biome uint32
}

// NewFlat creates a new Flat generator. Chunks generated are completely filled with the world.Biome passed. layers is a
// list of block layers placed by the Flat generator. The layers are ordered in a way where the last element in the
// slice is placed as the bottom-most block of the chunk.
func NewVoid(biome world.Biome) Flat {
	f := Flat{
		biome: uint32(biome.EncodeBiome()),
	}
	return f
}

// GenerateChunk ...
func (f Void) GenerateChunk(_ world.ChunkPos, chunk *chunk.Chunk) {
	min, max := int16(chunk.Range().Min()), int16(chunk.Range().Max())

	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			for y := int16(0); y <= max; y++ {
				chunk.SetBiome(x, min+y, z, f.biome)
			}
		}
	}
}
