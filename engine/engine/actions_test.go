package engine

import (
	"fmt"
	. "github.com/felzix/goblin"
	"github.com/felzix/huyilla/types"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestAction(t *testing.T) {
	g := Goblin(t)

	g.Describe("Action testing", func() {
		NAME := "FAKE"
		PASS := "PASS"
		var h *Engine

		g.BeforeEach(func() {
			unique, err := uuid.NewV4()
			if err != nil {
				t.Fatal(err)
			}

			h = &Engine{}
			if err := h.Init("/tmp/savedir-huyilla-" + unique.String()); err != nil {
				t.Fatal(err)
			}

			if err := h.SignUp(NAME, PASS); err != nil {
				t.Fatal(err)
			}

			g.Poll(5, 200, func() bool {
				_, err := h.LogIn(NAME, PASS)
				return err == nil
			})

		})

		g.AfterEach(func() {
			if h == nil || h.World == nil {
				return
			}
			if err := h.World.WipeDatabase(); err != nil {
				t.Fatal(err)
			}
		})

		g.Describe("queue actions and tick", func() {
			g.It("moves correctly", func() {
				whereTo := NewAbsolutePoint(0, 0, 0, 2, 4, 8)

				g.Assert(len(h.Actions)).Equal(0)

				h.RegisterAction(&types.Action{
					PlayerName: NAME,
					Action: &types.Action_Move{
						Move: &types.Action_MoveAction{
							WhereTo: whereTo,
						}}})

				g.Assert(len(h.Actions)).Equal(1)

				err := h.Tick()
				g.Assert(err).IsNil()

				g.Assert(len(h.Actions)).Equal(0)

				player, err := h.World.Player(NAME)
				g.Assert(err).IsNil()
				g.Assert(player).IsNotNil()
				g.Assert(player.EntityId).NotEqual(0)

				entity, err := h.World.Entity(player.EntityId)
				g.Assert(err).IsNil()
				g.Assert(entity).IsNotNil()
				g.Assert(absolutePointEquals(entity.Location, whereTo)).IsTrue(fmt.Sprintf(
					`Player should be at "%s" but is at "%s"`,
						absolutePointToString(whereTo),
						absolutePointToString(entity.Location),
				))

				chunk, err := h.World.Chunk(entity.Location.Chunk)
				g.Assert(err).IsNil()

				entityPresent := false
				for _, e := range chunk.Entities {
					if e == entity.Id {
						entityPresent = true
					}
				}
				g.Assert(entityPresent).IsTrue("Entity is not actually present in the chunk")
			})
		})
	})
}
