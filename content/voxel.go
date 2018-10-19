package content

import "github.com/felzix/huyilla/types"

var VoxelDefinitions = map[uint64]*types.VoxelDefinition{
    0: {Name: "air",
        State: types.VoxelDefinition_Gas},
    1: {Name: "barren_earth",
        State: types.VoxelDefinition_RigidSolid},
    2: {Name: "barren_grass",
        State: types.VoxelDefinition_RigidSolid},
    3: {Name: "sand",
        State: types.VoxelDefinition_LooseSolid},
    4: {Name: "water",
        State: types.VoxelDefinition_Liquid},
}
