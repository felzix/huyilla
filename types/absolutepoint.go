package types

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
)

type AbsolutePoint struct {
	Chunk Point
	Voxel Point
}

func NewAbsolutePoint(cX, cY, cZ, vX, vY, vZ int64) AbsolutePoint {
	return AbsolutePoint{
		Chunk: NewPoint(cX, cY, cZ),
		Voxel: NewPoint(vX, vY, vZ),
	}
}

func (p *AbsolutePoint) Clone() AbsolutePoint {
	return NewAbsolutePoint(p.Chunk.X, p.Chunk.Y, p.Chunk.Z, p.Voxel.X, p.Voxel.Y, p.Voxel.Z)
}

func (p AbsolutePoint) X() int64 {
	return (p.Chunk.X * C.CHUNK_SIZE) + p.Voxel.X
}
func (p AbsolutePoint) Y() int64 {
	return (p.Chunk.Y * C.CHUNK_SIZE) + p.Voxel.Y
}
func (p AbsolutePoint) Z() int64 {
	return (p.Chunk.Z * C.CHUNK_SIZE) + p.Voxel.Z
}

func (p AbsolutePoint) Derive(deltaX, deltaY, deltaZ int64) AbsolutePoint {
	size := int64(C.CHUNK_SIZE)
	derived := p.Clone()

	derived.Chunk.X += deltaX / size
	derived.Chunk.Z += deltaZ / size
	derived.Chunk.Y += deltaY / size

	derived.Voxel.X += deltaX % size
	derived.Voxel.Y += deltaY % size
	derived.Voxel.Z += deltaZ % size

	if derived.Voxel.X >= size {
		derived.Chunk.X += 1
		derived.Voxel.X -= size
	} else if derived.Voxel.X < 0 {
		derived.Chunk.X -= 1
		derived.Voxel.X += size
	}

	if derived.Voxel.Y >= size {
		derived.Chunk.Y += 1
		derived.Voxel.Y -= size
	} else if derived.Voxel.Y < 0 {
		derived.Chunk.Y -= 1
		derived.Voxel.Y += size
	}

	if derived.Voxel.Z >= size {
		derived.Chunk.Z += 1
		derived.Voxel.Z -= size
	} else if derived.Voxel.Z < 0 {
		derived.Chunk.Z -= 1
		derived.Voxel.Z += size
	}

	return derived
}

func (p AbsolutePoint) Neighbors(chunkSize uint64) []AbsolutePoint {
	voxels := make([]AbsolutePoint, 3*3*3-1)

	edge := []int64{-1, 0, +1}
	i := 0
	for _, x := range edge {
		for _, y := range edge {
			for _, z := range edge {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				voxels[i] = p.Derive(x, y, z)
				i++
			}
		}
	}

	return voxels
}

func (p AbsolutePoint) Equals(other AbsolutePoint) bool {
	return p.Chunk.Equals(other.Chunk) && p.Voxel.Equals(other.Voxel)
}

func (p AbsolutePoint) ToString() string {
	return fmt.Sprintf("(%s,%s)", p.Chunk.ToString(), p.Voxel.ToString())
}

func (p AbsolutePoint) Marshal() ([]byte, error) {
	return ToBytes(&p)
}

func (p *AbsolutePoint) Unmarshal(input []byte) error {
	return FromBytes(input, &p)
}