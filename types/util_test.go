package types

import (
	. "github.com/felzix/goblin"
	"math/rand"
	"testing"
)

func TestSpace(t *testing.T) {
	g := Goblin(t)

	g.Describe("Point calculations", func() {
		g.It("point equals another point", func() {
			PointEquals(&Point{X: -1, Y: 0, Z: 3}, &Point{X: -1, Y: 0, Z: 3})
		})

		g.It("creates an absolute point", func() {
			p := NewAbsolutePoint(-1, 1, 2, 3, 4, 5)

			g.Assert(PointEquals(p.Chunk, &Point{X: -1, Y: 1, Z: 2})).IsTrue()
			g.Assert(PointEquals(p.Voxel, &Point{X: 3, Y: 4, Z: 5})).IsTrue()
		})

		g.It("creates a point", func() {
			p := NewPoint(-1, 1, 2)

			g.Assert(PointEquals(p, &Point{X: -1, Y: 1, Z: 2})).IsTrue()
		})

		g.It("creates a random point", func() {
			rand.Seed(42)
			p := RandomPoint(16)

			g.Assert(PointEquals(p, &Point{X: 3, Y: 11, Z: 8})).IsTrue()
		})

		g.It("distance calc", func() {
			d := Distance(&Point{X: -1, Y: 0, Z: 3}, &Point{X: -10, Y: 9, Z: -3})

			g.Assert(d).Equal(14.071247279470288)
		})

		g.It("absolute int64", func() {
			g.Assert(absInt64(3)).Equal(int64(3))
			g.Assert(absInt64(-3)).Equal(int64(3))
		})

		g.It("grid distance", func() {
			g.Assert(GridDistance(NewPoint(0, 0, 0), NewPoint(0, 0, 1))).Equal(int64(1))
			g.Assert(GridDistance(NewPoint(0, 0, 0), NewPoint(0, 1, 1))).Equal(int64(1))
			g.Assert(GridDistance(NewPoint(0, 0, 0), NewPoint(1, 1, 1))).Equal(int64(1))
			g.Assert(GridDistance(NewPoint(0, 1, 0), NewPoint(0, 2, 1))).Equal(int64(1))
			g.Assert(GridDistance(NewPoint(0, 0, 0), NewPoint(0, 1, 2))).Equal(int64(2))
			g.Assert(GridDistance(NewPoint(-1, 0, 0), NewPoint(4, 1, 1))).Equal(int64(5))
		})
	})
}
