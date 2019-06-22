package types

import (
	"fmt"
	"math/rand"
	"strconv"
)

func NewChunk(tick, chunkLength uint64) *Chunk {
	return &Chunk{
		Tick:   tick,
		Voxels: make([]uint64, chunkLength),
	}
}

func NewChunks(radius uint64) *Chunks {
	diameter := 1 + radius*2
	size := diameter*diameter*diameter
	return &Chunks{
		Chunks: make([]*DetailedChunk, size),
		Points: make([]*Point, size),
	}
}

func NewMove(player string, whereTo *AbsolutePoint) *Action {
	return &Action{
		PlayerName: player,
		Action: &Action_Move{
			Move: &Action_MoveAction{
				WhereTo: whereTo,
			},
		},
	}
}

func NewAbsolutePoint(cX, cY, cZ, vX, vY, vZ int64) *AbsolutePoint {
	return &AbsolutePoint{
		Chunk: NewPoint(cX, cY, cZ),
		Voxel: NewPoint(vX, vY, vZ),
	}
}

func NewPoint(x, y, z int64) *Point {
	return &Point{X: x, Y: y, Z: z}
}

func (m *Point) Clone() *Point {
	return NewPoint(m.X, m.Y, m.Z)
}

func (m *AbsolutePoint) Clone() *AbsolutePoint {
	return NewAbsolutePoint(m.Chunk.X, m.Chunk.Y, m.Chunk.Z, m.Voxel.X, m.Voxel.Y, m.Voxel.Z)
}

func RandomPoint(chunkSize uint64) *Point {
	size := int64(chunkSize)
	x := rand.Int63n(size)
	y := rand.Int63n(size)
	z := rand.Int63n(size)
	return NewPoint(x, y, z)
}

func (m *Point) ToString() string {
	return fmt.Sprintf("(%d,%d,%d)", m.X, m.Y, m.Z)
}

func StringToPoint(x, y, z string) (*Point, error) {
	X, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return nil, err
	}

	Y, err := strconv.ParseInt(y, 10, 64)
	if err != nil {
		return nil, err
	}

	Z, err := strconv.ParseInt(z, 10, 64)
	if err != nil {
		return nil, err
	}

	return NewPoint(X, Y, Z), nil
}

func (m *AbsolutePoint) ToString() string {
	return fmt.Sprintf("(%s,%s)", m.Chunk.ToString(), m.Voxel.ToString())
}
