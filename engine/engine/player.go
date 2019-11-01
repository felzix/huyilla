package engine

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"strings"
)

func playerKey(name string) string {
	return fmt.Sprintf(`Player.%s`, name)
}

func playerNameFromKey(key string) string {
	// The "." (period) is not present because it's used as the filesystem separator.
	return strings.TrimPrefix(key, "Player")
}

func (world *World) Player(name string) (*types.Player, error) {
	var player types.Player
	if err := world.DB.Get(playerKey(name), &player); err == nil {
		return &player, nil
	} else if fileIsNotFound(err) {
		return nil, nil
	} else {
		return nil, err
	}
}

func (world *World) CreatePlayer(name string, password []byte, entityId int64, spawn *types.AbsolutePoint) error {
	player := types.Player{
		Name:     name,
		Password: password,
		EntityId: entityId,
		Spawn:    spawn,
		Token:    "",
	}
	return world.DB.Set(playerKey(player.Name), &player)
}

func (world *World) SetPlayer(player *types.Player) error {
	return world.DB.Set(playerKey(player.Name), player)
}

func (world *World) DeletePlayer(name string) error {
	return world.DB.End(playerKey(name))
}

func (world *World) GetActivePlayers() ([]*types.PlayerDetails, error) {
	var activePlayers []*types.PlayerDetails

	for key := range world.DB.GetByPrefix("Player") {
		name := playerNameFromKey(key)

		if player, err := world.Player(name); player != nil {
			if len(player.Token) > 0 {
				if entity, err := world.Entity(player.EntityId); err == nil {
					activePlayers = append(activePlayers, &types.PlayerDetails{Player: player, Entity: entity})
				} else {
					return nil, err
				}
			}
		} else if err == nil {
			return nil, errors.New(fmt.Sprintf(`Player "%s" should exist but doens't`, name))
		} else {
			return nil, err
		}
	}

	return activePlayers, nil
}
