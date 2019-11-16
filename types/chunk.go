package types

import (
	C "github.com/felzix/huyilla/constants"
)

type Chunk struct {
	Tick Age
	Voxels []uint64
	Entities []EntityId
	Items []ItemId
}

func NewChunk(tick Age, chunkLength uint64) *Chunk {
	return &Chunk{
		Tick:   tick,
		Voxels: make([]uint64, chunkLength),
		Entities: make([]EntityId, 0),
		Items: make([]ItemId, 0),
	}
}

func (c Chunk) GetVoxel(x, y, z uint64) Voxel {
	index := c.GetVoxelIndex(x, y, z)
	return Voxel(c.Voxels[index])
}

func (c Chunk) GetVoxelIndex(x, y, z uint64) int {
	return CalculateVoxelIndex(x, y, z)
}

func (c Chunk) Marshal() ([]byte, error) {
	return ToBytes(c)
}

func (c *Chunk) Unmarshal(blob []byte) error {
	return FromBytes(blob, &c)
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
