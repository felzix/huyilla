package types

import (
	"fmt"
	. "github.com/felzix/goblin"
	"math/rand"
	"testing"
)

func TestSpace(t *testing.T) {
	g := Goblin(t)

	g.Describe("Point calculations", func() {
		g.It("point equals another point", func() {
			g.Assert(NewPoint(-1, 0, 3).Equals(NewPoint(-1, 0, 3))).IsTrue()
		})

		g.It("creates an absolute point", func() {
			p := NewAbsolutePoint(-1, 1, 2, 3, 4, 5)

			g.Assert(p.Chunk.Equals(&Point{X: -1, Y: 1, Z: 2})).IsTrue()
			g.Assert(p.Voxel.Equals(&Point{X: 3, Y: 4, Z: 5})).IsTrue()
		})

		g.It("creates a point", func() {
			p := NewPoint(-1, 1, 2)

			g.Assert(p.Equals(&Point{X: -1, Y: 1, Z: 2})).IsTrue()
		})

		g.It("creates a random point", func() {
			rand.Seed(42)
			p := RandomPoint(16)

			g.Assert(p.Equals(&Point{X: 3, Y: 11, Z: 8})).IsTrue()
		})

		g.It("distance calc", func() {
			d := NewPoint(-1, 0, 3).Distance(NewPoint(-10, 9, -3))

			g.Assert(d).Equal(14.071247279470288)
		})

		g.It("absolute int64", func() {
			g.Assert(absInt64(3)).Equal(int64(3))
			g.Assert(absInt64(-3)).Equal(int64(3))
		})

		g.It("grid distance", func() {
			g.Assert(NewPoint(0, 0, 0).GridDistance(NewPoint(0, 0, 1))).Equal(int64(1))
			g.Assert(NewPoint(0, 0, 0).GridDistance(NewPoint(0, 1, 1))).Equal(int64(1))
			g.Assert(NewPoint(0, 0, 0).GridDistance(NewPoint(1, 1, 1))).Equal(int64(1))
			g.Assert(NewPoint(0, 1, 0).GridDistance(NewPoint(0, 2, 1))).Equal(int64(1))
			g.Assert(NewPoint(0, 0, 0).GridDistance(NewPoint(0, 1, 2))).Equal(int64(2))
			g.Assert(NewPoint(-1, 0, 0).GridDistance(NewPoint(4, 1, 1))).Equal(int64(5))
		})

		g.It("derives a positive delta", func() {
			const CHUNK_SIZE = 16
			p := NewAbsolutePoint(1, 2, 3, 4, 5, 6)
			d := p.Derive(40, 40, 40, CHUNK_SIZE)
			e := NewAbsolutePoint(3, 4, 5, 12, 13, 14)
			g.Assert(d.Equals(e)).IsTrue(fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
		})

		g.It("derives a negative delta", func() {
			const CHUNK_SIZE = 16
			p := NewAbsolutePoint(1, 2, 3, 4, 5, 6)
			d := p.Derive(-40, -40, -40, CHUNK_SIZE)
			e := NewAbsolutePoint(-1, 0, 1, 12, 13, 14)
			g.Assert(d.Equals(e)).IsTrue(fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
		})

		g.It("derives chunkSize", func() {
			const CHUNK_SIZE = 16
			p := NewAbsolutePoint(1, 1, 1, 1, 0, 1)
			d := p.Derive(-CHUNK_SIZE, -CHUNK_SIZE*2, +CHUNK_SIZE, CHUNK_SIZE)
			e := NewAbsolutePoint(0, -1, 2, 1, 0, 1)
			g.Assert(d.Equals(e)).IsTrue(fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
		})
	})
}
