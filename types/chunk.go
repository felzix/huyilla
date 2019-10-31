package types

import (
	C "github.com/felzix/huyilla/constants"
)

func (m Chunk) GetVoxel(x, y, z uint64) Voxel {
	index := m.GetVoxelIndex(x, y, z)
	return Voxel(m.Voxels[index])
}

func (m Chunk) GetVoxelIndex(x, y, z uint64) int {
	return CalculateVoxelIndex(x, y, z)
}

func (m DetailedChunk) GetVoxel(x, y, z uint64) Voxel {
	index := m.GetVoxelIndex(x, y, z)
	return Voxel(m.Voxels[index])
}

func (m DetailedChunk) GetVoxelIndex(x, y, z uint64) int {
	return CalculateVoxelIndex(x, y, z)
}

func CalculateVoxelIndex(x, y, z uint64) int {
	return int((x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + z)
}

func EachVoxel(fn func(uint64, uint64, uint64)) {
	for x := 0; x < C.CHUNK_SIZE; x++ {
		for y := 0; y < C.CHUNK_SIZE; y++ {
			for z := 0; z < C.CHUNK_SIZE; z++ {
				fn(uint64(x), uint64(y), uint64(z))
			}
		}
	}
}
