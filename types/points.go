package types

import (
	"math"
)

/*

   def find_new_coords(self, x_delta, y_delta, z_delta, chunk_size):
       newcoords = self.copy()

       def incr(point, dim, by):
           original_value = getattr(point, dim)
           new_value = original_value + by
           setattr(point, dim, new_value)

       def once(dim, delta):
           # int modulo and division are a little weird w/ negatives
           mod_fix = -1 if delta < 0 else +1

           fixed_chunk_size = chunk_size * mod_fix

           incr(newcoords.chunk_coords, dim, delta // fixed_chunk_size)
           incr(newcoords.voxel_coords, dim, delta % fixed_chunk_size)

           if getattr(newcoords.voxel_coords, dim) >= chunk_size:
               incr(newcoords.chunk_coords, dim, +1)
               incr(newcoords.voxel_coords, dim, -chunk_size)
           elif getattr(newcoords.voxel_coords, dim) < 0:
               incr(newcoords.chunk_coords, dim, -1)
               incr(newcoords.voxel_coords, dim, +chunk_size)

       once("x", x_delta)
       once("y", y_delta)
       once("z", z_delta)

       return newcoords

*/

func (m *AbsolutePoint) Derive(deltaX, deltaY, deltaZ int64, chunkSize uint64) *AbsolutePoint {
	size := int64(chunkSize)
	derived := m.Clone()

	var modFix int64
	if deltaX < 0 {
		modFix = -1
	} else {
		modFix = +1
	}
	fixedSize := size * modFix

	derived.Chunk.X += deltaX / fixedSize * modFix
	derived.Chunk.Y += deltaY / fixedSize * modFix
	derived.Chunk.Z += deltaZ / fixedSize * modFix

	derived.Voxel.X += deltaX % fixedSize * modFix
	derived.Voxel.Y += deltaY % fixedSize * modFix
	derived.Voxel.Z += deltaZ % fixedSize * modFix

	return derived
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
