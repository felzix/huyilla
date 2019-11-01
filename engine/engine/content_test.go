package engine

import (
	C "github.com/felzix/huyilla/constants"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestContent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Content Suite")
}

var _ = Describe("Content", func() {
	var h *Engine

	BeforeEach(func() {
		engine, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
		h = engine
		Expect(err).To(BeNil())
	})

	It("loads human type", func() {
		content := h.GetContent()
		Expect(content.E[0].Name).To(Equal("human"))
	})
})
