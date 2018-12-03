package main

import (
	"fmt"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/gogo/protobuf/proto"
	"github.com/peterbourgon/diskv"
)

type World struct {
	DB *diskv.Diskv

	Entities map[int64]*types.Entity
	Chunks   map[types.Point]*types.Chunk
}

func (world *World) Init(saveDir string, cacheSize uint64) error {
	// So that recipes and terrain generator can reference content by name.
	content.PopulateContentNameMaps()

	if db, err := makeDB(saveDir, cacheSize); err == nil {
		world.DB = db
	} else {
		return err
	}

	if !world.DB.Has(KEY_AGE) {
		defaultAge := types.Age{1}
		if blob, err := proto.Marshal(&defaultAge); err == nil {
			world.DB.Write(KEY_AGE, blob)
		} else {
			return err
		}
	}

	return nil
}

//
// Age
//

const KEY_AGE = "Age"

func (world *World) Age() (*types.Age, error) {
	var age types.Age
	if err := gettem(world, KEY_AGE, &age); err == nil {
		return &age, nil
	} else {
		return nil, err
	}
}

func (world *World) IncrementAge() error {
	if age, err := world.Age(); err == nil {
		age.Ticks++
		if err := settem(world, KEY_AGE, age); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

//
// Player
//

func playerKey(name string) string {
	return fmt.Sprintf(`Player.%s`, name)
}

func (world *World) Player(name string) (*types.Player, error) {
	var player types.Player
	if err := gettem(world, playerKey(name), &player); err == nil {
		return &player, nil
	} else if fileIsNotFound(err) {
		return nil, nil
	} else {
		return nil, err
	}
}

func (world *World) CreatePlayer(player *types.Player) error {
	return settem(world, playerKey(player.Name), player)
}

func (world *World) SetPlayerPassword(name, password string) error {
	var player types.Player
	if err := gettem(world, playerKey(name), &player); err != nil {
		return err
	}
	if hashedPassword, err := hashPassword(password); err == nil {
		player.Password = hashedPassword
	}
	return settem(world, playerKey(player.Name), &player)
}

func (world *World) SetPlayerLogin(name string, login bool) error {
	var player types.Player
	if err := gettem(world, playerKey(name), &player); err != nil {
		return err
	}
	player.LoggedIn = login
	return settem(world, playerKey(player.Name), &player)
}

func (world *World) SetPlayerSpawn(name string, spawnPoint *types.AbsolutePoint) error {
	var player types.Player
	if err := gettem(world, playerKey(name), &player); err != nil {
		return err
	}
	player.Spawn = spawnPoint
	return settem(world, playerKey(player.Name), &player)
}
