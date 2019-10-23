package types

import (
	C "github.com/felzix/huyilla/constants"
)

func (m Chunk) GetVoxel(x, y, z uint64) Voxel {
	index := m.GetVoxelIndex(x, y, z)
	return Voxel(m.Voxels[index])
}

func (m Chunk) GetVoxelIndex(x, y, z uint64) int {
	return int((x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + z)
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
