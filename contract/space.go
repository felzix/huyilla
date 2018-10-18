package main

import (
    "github.com/felzix/huyilla/types"
)


func newAbsolutePoint(cX, cY, cZ, vX, vY, vZ int64) *types.AbsolutePoint {
    return &types.AbsolutePoint{
        Chunk: newPoint(cX, cY, cZ),
        Voxel: newPoint(vX, vY, vZ),
    }
}

func newPoint(x, y, z int64) *types.Point {
    return &types.Point{x, y, z}
}

func pointEquals (p0, p1 *types.Point) bool {
    return p0.X == p1.X &&
           p0.Y == p1.Y &&
           p0.Z == p1.Z
}