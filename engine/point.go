package engine

import "math"

type Point struct {
    x int
    y int
    z int
}

func (p0 Point) Distance(p1 Point) float64 {
    x := float64(p0.x - p1.x)
    y := float64(p0.y - p1.y)
    z := float64(p0.z - p1.z)

    return math.Sqrt(x*x + y*y + z*z)
}

type AbsolutePoint struct {
    chunkCoords Point
    voxelCoords Point
}
