package types

// A voxel is a 64-bit array stored as a uint64.
// The bits mean the following, from least to most significant:
//
// 0     : 1  : multiblock? : determines the meaning of the bits 2-63
// 1-4   : 4  : pressure or sturdiness : 0-15
// 5-17  : 13 : temperature : 0-8191
// 18-31 : 14 : other (meaning assigned per-block)
// 32-47 : 16 : material 1 : what the voxel is made of
// 48-63 : 16 : form : shape of material(s) in voxel

const (
	MinTemperature  = uint64(0)
	RoomTemperature = uint64(295)
	MaxTemperature  = uint64(0)

	MinPressure = uint64(0)
	MaxPressure = uint64(15)
)

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
		Temperature: uint64(v & 0x000000000003FFE0 >> 5),
		Pressure:    uint64(v & 0x000000000000001E >> 1),
		Multiblock:  uint64(v & 0x0000000000000001 >> 0),
	}
}

// Trusts input to only have valid values.
// If e.g. Pressure were more than 4 bits then Temperature would be corrupted.
func (voxel ExpandedVoxel) Compress() Voxel {
	v := uint64(0)
	v |= voxel.Form << 48
	v |= voxel.Material << 32
	v |= voxel.Other << 18
	v |= voxel.Temperature << 5
	v |= voxel.Pressure << 1
	v |= voxel.Multiblock << 0
	return Voxel(v)
}
