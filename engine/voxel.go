package engine


type Voxel struct {
    vType      VoxelType
    properties uint32
    physics    uint16
}

type VoxelType     uint16
type VoxelProperty uint


func NewVoxel(vType VoxelType) Voxel {
    return Voxel{vType,uint32(0), uint16(0)}
}
