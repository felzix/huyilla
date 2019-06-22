package types

import (
	"math"
)

func (m *AbsolutePoint) Derive(deltaX, deltaY, deltaZ int64, chunkSize uint64) *AbsolutePoint {
	size := int64(chunkSize)
	derived := m.Clone()

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

func (m *AbsolutePoint) Neighbors(chunkSize uint64) []*AbsolutePoint {
	voxels := make([]*AbsolutePoint, 3*3*3-1)

	edge := []int64{-1, 0, +1}
	i := 0
	for _, x := range edge {
		for _, y := range edge {
			for _, z := range edge {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				voxels[i] = m.Derive(x, y, z, chunkSize)
				i++
			}
		}
	}

	return voxels
}

func (m *Point) Equals(other *Point) bool {
	return m.X == other.X &&
		m.Y == other.Y &&
		m.Z == other.Z
}

func (m *AbsolutePoint) Equals(other *AbsolutePoint) bool {
	return m.Chunk.Equals(other.Chunk) && m.Voxel.Equals(other.Voxel)
}

func (m *Point) Distance(other *Point) float64 {
	x := float64(m.X - other.X)
	y := float64(m.Y - other.Y)
	z := float64(m.Z - other.Z)

	return math.Sqrt(x*x + y*y + z*z)
}

func (m *Point) GridDistance(other *Point) int64 {
	x := absInt64(m.X - other.X)
	y := absInt64(m.Y - other.Y)
	z := absInt64(m.Z - other.Z)

	most := x
	if y > most {
		most = y
	}
	if z > most {
		most = z
	}

	return most
}

// If you don't need to know the specific distance,
// just if it's more or less than another distance,
// there's no need to take the square root.
func (m *Point) RelativeDistance(other *Point) float64 {
	x := float64(m.X - other.X)
	y := float64(m.Y - other.Y)
	z := float64(m.Z - other.Z)

	return x*x + y*y + z*z
}
