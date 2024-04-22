package block

import (
	"github.com/STCraft/dragonfly/server/block/model"
	"github.com/STCraft/dragonfly/server/world"
)

// SolidModel represents a block that is fully SolidModel. It always returns a model.SolidModel when Model is called.
type SolidModel struct{}

// Model ...
func (SolidModel) Model() world.BlockModel {
	return model.Solid{}
}

// EmptyModel represents a block that is fully EmptyModel/transparent, such as air or a plant. It always returns a
// model.EmptyModel when Model is called.
type EmptyModel struct{}

// Model ...
func (EmptyModel) Model() world.BlockModel {
	return model.Empty{}
}

// ChestModel represents a block that has a model of a ChestModel.
type ChestModel struct{}

// Model ...
func (ChestModel) Model() world.BlockModel {
	return model.Chest{}
}

// CarpetModel represents a block that has a model of a CarpetModel.
type CarpetModel struct{}

// Model ...
func (CarpetModel) Model() world.BlockModel {
	return model.Carpet{}
}

// TilledGrassModel represents a block that has a model of farmland or dirt paths.
type TilledGrassModel struct{}

// Model ...
func (TilledGrassModel) Model() world.BlockModel {
	return model.TilledGrass{}
}

// LeavesModel represents a block that has a model of LeavesModel. A full block but with no solid faces.
type LeavesModel struct{}

// Model ...
func (LeavesModel) Model() world.BlockModel {
	return model.Leaves{}
}

// ThinModel represents a ThinModel, partial block such as a glass pane or an iron bar, that connects to nearby solid faces.
type ThinModel struct{}

// Model ...
func (ThinModel) Model() world.BlockModel {
	return model.Thin{}
}
