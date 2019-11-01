package engine

import (
	C "github.com/felzix/huyilla/constants"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPlayer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Player Suite")
}

var _ = Describe("Player", func () {
	NAME := "felzix"
	PASS := "murakami"
	var h *Engine

	BeforeEach(func() {
		engine, err := NewEngine(C.SEED, NewLakeWorldGenerator(3), NewMemoryDatabase())
		h = engine
		Expect(err).To(BeNil())
	})

	It("loads human type", func() {
		err := h.SignUp(NAME, PASS)
		Expect(err).To(BeNil())
		_, err = h.LogIn(NAME, PASS)
		Expect(err).To(BeNil())

		player, err := h.World.Player(NAME)
		Expect(err).To(BeNil())
		Expect(player).ToNot(BeNil())
		Expect(player.Name).To(Equal(NAME))
	})
})
