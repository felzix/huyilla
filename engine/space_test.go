package main

import (
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/types"
	"math/rand"
	"testing"
)

func TestSpace(t *testing.T) {
	g := Goblin(t)

	g.Describe("Point calculations", func() {
		g.It("point equals another point", func() {
			pointEquals(&types.Point{X: -1, Y: 0, Z: 3}, &types.Point{X: -1, Y: 0, Z: 3})
		})

		g.It("creates an absolute point", func() {
			p := newAbsolutePoint(-1, 1, 2, 3, 4, 5)

			g.Assert(pointEquals(p.Chunk, &types.Point{X: -1, Y: 1, Z: 2})).IsTrue()
			g.Assert(pointEquals(p.Voxel, &types.Point{X: 3, Y: 4, Z: 5})).IsTrue()
		})

		g.It("creates a point", func() {
			p := newPoint(-1, 1, 2)

			g.Assert(pointEquals(p, &types.Point{X: -1, Y: 1, Z: 2})).IsTrue()
		})

		g.It("creates a random point", func() {
			rand.Seed(42)
			p := randomPoint()

			g.Assert(pointEquals(p, &types.Point{X: 3, Y: 11, Z: 8})).IsTrue()
		})

		g.It("distance calc", func() {
			d := distance(&types.Point{X: -1, Y: 0, Z: 3}, &types.Point{X: -10, Y: 9, Z: -3})

			g.Assert(d).Equal(14.071247279470288)
		})

		g.It("absolute int64", func() {
			g.Assert(absInt64(3)).Equal(int64(3))
			g.Assert(absInt64(-3)).Equal(int64(3))
		})

		g.It("grid distance", func() {
			g.Assert(gridDistance(newPoint(0, 0, 0), newPoint(0, 0, 1))).Equal(int64(1))
			g.Assert(gridDistance(newPoint(0, 0, 0), newPoint(0, 1, 1))).Equal(int64(1))
			g.Assert(gridDistance(newPoint(0, 0, 0), newPoint(1, 1, 1))).Equal(int64(1))
			g.Assert(gridDistance(newPoint(0, 1, 0), newPoint(0, 2, 1))).Equal(int64(1))
			g.Assert(gridDistance(newPoint(0, 0, 0), newPoint(0, 1, 2))).Equal(int64(2))
			g.Assert(gridDistance(newPoint(-1, 0, 0), newPoint(4, 1, 1))).Equal(int64(5))
		})
	})
}
