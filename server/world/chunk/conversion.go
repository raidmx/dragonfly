package chunk

import (
	"bytes"
	_ "embed"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

// legacyBlockEntry represents a block entry used in versions prior to 1.13.
type legacyBlockEntry struct {
	Name string `nbt:"name"`
	Meta int16  `nbt:"meta"`
}

var (
	//go:embed legacy_states.nbt
	legacyMappingsData []byte
	// legacyMappings allows simple conversion from a legacy block entry to a new one.
	legacyMappings = make(map[legacyBlockEntry]blockEntry)
)

// upgradeLegacyEntry upgrades a legacy block entry to a new one.
func upgradeLegacyEntry(name string, meta int16) (blockEntry, bool) {
	entry, ok := legacyMappings[legacyBlockEntry{Name: name, Meta: meta}]
	if !ok {
		// Also try cases where the meta should be disregarded.
		entry, ok = legacyMappings[legacyBlockEntry{Name: name}]
	}
	return entry, ok
}

// init creates conversions for each legacy and alias entry.
func init() {
	var entry struct {
		Legacy  legacyBlockEntry `nbt:"legacy"`
		Updated blockEntry       `nbt:"updated"`
	}
	dec := nbt.NewDecoder(bytes.NewBuffer(legacyMappingsData))
	for {
		if err := dec.Decode(&entry); err != nil {
			break
		}
		legacyMappings[entry.Legacy] = entry.Updated
	}
}
