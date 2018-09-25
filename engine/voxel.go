package engine


type Voxel struct {
    Type      VoxelType
    Properties uint32
    Physics    uint16
}

type VoxelType     uint16
type VoxelProperty struct {
    RigidSolid bool  // like stone
    LooseSolid bool  // like sand
    Liquid     bool  // like water
    Gas        bool  // like air
}


func NewVoxel(vType VoxelType) Voxel {
    return Voxel{vType,uint32(0), uint16(0)}
}
