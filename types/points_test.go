package types


import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"
	"testing"
)

func TestPoints(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Points Suite")
}

var _ = Describe("Point Calculations", func() {
	It("point equals another point", func() {
		Expect(NewPoint(-1, 0, 3).Equals(NewPoint(-1, 0, 3))).To(Equal(true))
	})

	It("creates an absolute point", func() {
		p := NewAbsolutePoint(-1, 1, 2, 3, 4, 5)

		Expect(p.Chunk.Equals(&Point{X: -1, Y: 1, Z: 2})).To(Equal(true))
		Expect(p.Voxel.Equals(&Point{X: 3, Y: 4, Z: 5})).To(Equal(true))
	})

	It("creates a point", func() {
		p := NewPoint(-1, 1, 2)

		Expect(p.Equals(&Point{X: -1, Y: 1, Z: 2})).To(Equal(true))
	})

	It("creates a random point", func() {
		rand.Seed(42)
		p := RandomPoint(16)

		Expect(p.Equals(&Point{X: 3, Y: 11, Z: 8})).To(Equal(true))
	})

	It("distance calc", func() {
		d := NewPoint(-1, 0, 3).Distance(NewPoint(-10, 9, -3))

		Expect(d).To(Equal(14.071247279470288))
	})

	It("absolute int64", func() {
		Expect(absInt64(3)).To(Equal(int64(3)))
		Expect(absInt64(-3)).To(Equal(int64(3)))
	})

	It("grid distance", func() {
		Expect(NewPoint(0, 0, 0).GridDistance(NewPoint(0, 0, 1))).To(Equal(int64(1)))
		Expect(NewPoint(0, 0, 0).GridDistance(NewPoint(0, 1, 1))).To(Equal(int64(1)))
		Expect(NewPoint(0, 0, 0).GridDistance(NewPoint(1, 1, 1))).To(Equal(int64(1)))
		Expect(NewPoint(0, 1, 0).GridDistance(NewPoint(0, 2, 1))).To(Equal(int64(1)))
		Expect(NewPoint(0, 0, 0).GridDistance(NewPoint(0, 1, 2))).To(Equal(int64(2)))
		Expect(NewPoint(-1, 0, 0).GridDistance(NewPoint(4, 1, 1))).To(Equal(int64(5)))
	})

	It("derives a positive delta", func() {
		const CHUNK_SIZE = 16
		p := NewAbsolutePoint(1, 2, 3, 4, 5, 6)
		d := p.Derive(40, 40, 10, CHUNK_SIZE)
		e := NewAbsolutePoint(3, 4, 4, 12, 13, 0)
		Expect(d.Equals(e)).To(Equal(true), fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
	})

	It("derives a negative delta", func() {
		const CHUNK_SIZE = 16
		p := NewAbsolutePoint(1, 2, 3, 4, 5, 6)
		d := p.Derive(-40, -40, -40, CHUNK_SIZE)
		e := NewAbsolutePoint(-2, -1, 0, 12, 13, 14)
		Expect(d.Equals(e)).To(Equal(true), fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
	})

	It("derives a neighbor", func() {
		const CHUNK_SIZE = 16
		p := NewAbsolutePoint(1, 2, 3, 0, 2, 15)
		d := p.Derive(-1, -1, -1, CHUNK_SIZE)
		e := NewAbsolutePoint(0, 2, 3, 15, 1, 14)
		Expect(d.Equals(e)).To(Equal(true), fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
	})

	It("derives chunkSize", func() {
		const CHUNK_SIZE = 16
		p := NewAbsolutePoint(1, 1, 1, 1, 0, 1)
		d := p.Derive(-CHUNK_SIZE, -CHUNK_SIZE*2, +CHUNK_SIZE, CHUNK_SIZE)
		e := NewAbsolutePoint(0, -1, 2, 1, 0, 1)
		Expect(d.Equals(e)).To(Equal(true), fmt.Sprintf("%s != %s", d.ToString(), e.ToString()))
	})

	It("knows its neighbors", func() {
		const CHUNK_SIZE = 16
		p := NewAbsolutePoint(1, 2, 3, 0, 2, 15)
		n := p.Neighbors(CHUNK_SIZE)

		test := func(i int, e *AbsolutePoint) {
			Expect(n[i].Equals(e)).To(Equal(true), fmt.Sprintf("neighbor %d: %s != %s", i, n[i].ToString(), e.ToString()))
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
