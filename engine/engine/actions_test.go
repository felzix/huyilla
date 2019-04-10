package engine

import (
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

				if len(h.Actions) != 0 {
					t.Errorf(`Expected 0 action but found %d`, len(h.Actions))
				}

				h.RegisterAction(&types.Action{
					PlayerName: NAME,
					Action: &types.Action_Move{
						Move: &types.Action_MoveAction{
							WhereTo: whereTo,
						}}})

				if len(h.Actions) != 1 {
					t.Errorf(`Expected 1 action but found %d`, len(h.Actions))
				}

				if err := h.Tick(); err != nil {
					t.Fatal(err)
				}

				if len(h.Actions) != 0 {
					t.Errorf(`Expected 0 action but found %d`, len(h.Actions))
				}

				player, err := h.World.Player(NAME)
				if err != nil {
					t.Fatal("Error:", err)
				} else if player == nil {
					t.Fatal("The player doesn't exist")
				}
				if player.EntityId == 0 {
					t.Fatal("The player doesn't have an entity id")
				}

				entity, err := h.World.Entity(player.EntityId)
				if entity == nil {
					t.Fatal("Entity should exist but doesn't")
				} else if err != nil {
					t.Fatal(err)
				}

				if !(absolutePointEquals(entity.Location, whereTo)) {
					t.Errorf(`Player should be at "%s" but is at "%s"`,
						absolutePointToString(whereTo),
						absolutePointToString(entity.Location))
				}

				chunk, err := h.World.Chunk(entity.Location.Chunk)
				if err != nil {
					t.Fatal(err)
				}

				entityPresent := false
				for _, e := range chunk.Entities {
					if e == entity.Id {
						entityPresent = true
					}
				}
				if !entityPresent {
					t.Error("Entity is not actually present in the chunk")
				}
			})
		})
	})
}
