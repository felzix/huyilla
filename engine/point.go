package engine

import "math"

type Point struct {
    X int
    Y int
    Z int
}

func (p0 Point) Distance(p1 Point) float64 {
    x := float64(p0.X - p1.X)
    y := float64(p0.Y - p1.Y)
    z := float64(p0.Z - p1.Z)

    return math.Sqrt(x*x + y*y + z*z)
}

type AbsolutePoint struct {
    ChunkCoords Point
    VoxelCoords Point
}
