package content

import "github.com/felzix/huyilla/types"

var VoxelDefinitions = map[uint64]*types.VoxelDefinition{
	uint64(0): {Name: "air",
		State: types.VoxelDefinition_Gas},
	uint64(1): {Name: "barren_earth",
		State: types.VoxelDefinition_RigidSolid},
	uint64(2): {Name: "barren_grass",
		State: types.VoxelDefinition_RigidSolid},
	uint64(3): {Name: "sand",
		State: types.VoxelDefinition_LooseSolid},
	uint64(4): {Name: "water",
		State: types.VoxelDefinition_Liquid},
}
