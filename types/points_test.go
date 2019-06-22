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
			d := p.Derive(40, 40, 10, CHUNK_SIZE)
			e := NewAbsolutePoint(3, 4, 4, 12, 13, 0)
			g.Assert(d.Equals(e)).IsTrue(fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
		})

		g.It("derives a negative delta", func() {
			const CHUNK_SIZE = 16
			p := NewAbsolutePoint(1, 2, 3, 4, 5, 6)
			d := p.Derive(-40, -40, -40, CHUNK_SIZE)
			e := NewAbsolutePoint(-2, -1, 0, 12, 13, 14)
			g.Assert(d.Equals(e)).IsTrue(fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
		})

		g.It("derives a neighbor", func() {
			const CHUNK_SIZE = 16
			p := NewAbsolutePoint(1, 2, 3, 0, 2, 15)
			d := p.Derive(-1, -1, -1, CHUNK_SIZE)
			e := NewAbsolutePoint(0, 2, 3, 15, 1, 14)
			g.Assert(d.Equals(e)).IsTrue(fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
		})

		g.It("derives chunkSize", func() {
			const CHUNK_SIZE = 16
			p := NewAbsolutePoint(1, 1, 1, 1, 0, 1)
			d := p.Derive(-CHUNK_SIZE, -CHUNK_SIZE*2, +CHUNK_SIZE, CHUNK_SIZE)
			e := NewAbsolutePoint(0, -1, 2, 1, 0, 1)
			g.Assert(d.Equals(e)).IsTrue(fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
		})

		g.It("knows its neighbors", func() {
			const CHUNK_SIZE = 16
			p := NewAbsolutePoint(1, 2, 3, 0, 2, 15)
			n := p.Neighbors(CHUNK_SIZE)

			test := func(i int, e *AbsolutePoint) {
				g.Assert(n[i].Equals(e)).IsTrue(fmt.Sprintf("neighbor %d: %s != %s", i, n[i].ToString(), e.ToString()))
			}

			test(0, NewAbsolutePoint(0, 2, 3, 15, 1, 14))
			test(1, NewAbsolutePoint(0, 2, 3, 15, 1, 15))
			test(2, NewAbsolutePoint(0, 2, 4, 15, 1, 0))
			test(3, NewAbsolutePoint(0, 2, 3, 15, 2, 14))
			test(4, NewAbsolutePoint(0, 2, 3, 15, 2, 15))
			test(5, NewAbsolutePoint(0, 2, 4, 15, 2, 0))
			test(6, NewAbsolutePoint(0, 2, 3, 15, 3, 14))
			test(7, NewAbsolutePoint(0, 2, 3, 15, 3, 15))
			test(8, NewAbsolutePoint(0, 2, 4, 15, 3, 0))
			test(9, NewAbsolutePoint(1, 2, 3, 0, 1, 14))
			test(10, NewAbsolutePoint(1, 2, 3, 0, 1, 15))
			test(11, NewAbsolutePoint(1, 2, 4, 0, 1, 0))
			test(12, NewAbsolutePoint(1, 2, 3, 0, 2, 14))
			test(13, NewAbsolutePoint(1, 2, 4, 0, 2, 0))
			test(14, NewAbsolutePoint(1, 2, 3, 0, 3, 14))
			test(15, NewAbsolutePoint(1, 2, 3, 0, 3, 15))
			test(16, NewAbsolutePoint(1, 2, 4, 0, 3, 0))
			test(17, NewAbsolutePoint(1, 2, 3, 1, 1, 14))
			test(18, NewAbsolutePoint(1, 2, 3, 1, 1, 15))
			test(19, NewAbsolutePoint(1, 2, 4, 1, 1, 0))
			test(20, NewAbsolutePoint(1, 2, 3, 1, 2, 14))
			test(21, NewAbsolutePoint(1, 2, 3, 1, 2, 15))
			test(22, NewAbsolutePoint(1, 2, 4, 1, 2, 0))
			test(23, NewAbsolutePoint(1, 2, 3, 1, 3, 14))
			test(24, NewAbsolutePoint(1, 2, 3, 1, 3, 15))
			test(25, NewAbsolutePoint(1, 2, 4, 1, 3, 0))
		})
	})
}
