package types

// A voxel is a 64-bit array stored as a uint64.
// The bits mean the following, from least to most significant:
//
// 0     : 1  : multiblock? : determines the meaning of the bits 2-63
// 1-4   : 4  : pressure or sturdiness : 0-15
// 5-17  : 13 : temperature : 0-8191
// 18-31 : 14 : meaning assigned per-block
// 32-47 : 16 : material 1 : what the voxel is made of
// 48-63 : 16 : form : shape of material(s) in voxel

const RoomTemperature = uint64(295)

type ExpandedVoxel struct {
	Form        uint64
	Material    uint64
	Other       uint64
	Temperature uint64
	Pressure    uint64 // TODO consider using a different struct where only this field is named differently
	Multiblock  uint64 // TODO use different structure for multiblocks
}

type Voxel uint64

func (v Voxel) Expand() ExpandedVoxel {
	return ExpandedVoxel{
		Form:        uint64(v & 0xFFFF000000000000 >> 48),
		Material:    uint64(v & 0x0000FFFF00000000 >> 32),
		Other:       uint64(v & 0x00000000FFFC0000 >> 18),
		Temperature: uint64(v & 0x00000000000FFFE0 >> 5),
		Pressure:    uint64(v & 0x000000000000000E >> 1),
		Multiblock:  uint64(v & 0x0000000000000001 >> 0),
	}
}

func (voxel ExpandedVoxel) Compress() Voxel {
	v := uint64(0)
	v |= voxel.Form << 48
	v |= voxel.Material << 32
	v |= RoomTemperature << 5
	return Voxel(v)
}
