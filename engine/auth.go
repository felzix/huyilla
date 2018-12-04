package main

import (
	"fmt"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (engine *Engine) SignUp(name string, password string) error {
	if player, err := engine.World.Player(name); player != nil {
		return errors.New(fmt.Sprintf(`Player "%s" already exists`, name))
	} else if err != nil {
		return err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	// Create new player
	defaultLocation := newAbsolutePoint(0, 0, 0, 0, 0, 0)

	entity, err := engine.World.CreateEntity(ENTITY["human"], name, defaultLocation)
	if err != nil {
		return err
	}

	engine.World.CreatePlayer(name, hashedPassword, entity.Id, defaultLocation)

	return nil
}

func (engine *Engine) LogIn(name, password string) (*types.PlayerDetails, error) {
	player, err := engine.World.Player(name)
	if player == nil {
		return nil, errors.New(fmt.Sprintf(`No such player "%s"`, name))
	} else if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword(player.Password, []byte(password)) != nil {
		return nil, errors.New("Incorrect password")
	}

	if player.LoggedIn {
		return nil, errors.New("You are already logged in.")
	}
	player.LoggedIn = true

	entity, err := engine.World.Entity(player.EntityId)
	if entity == nil {
		return nil, errors.New(fmt.Sprintf(`Player's entity does not exist: "%d"`, player.EntityId))
	} else if err != nil {
		return nil, err
	}
	if err := engine.World.AddEntityToChunk(entity); err != nil {
		return nil, err
	}

	return &types.PlayerDetails{Player: player, Entity: entity}, nil
}

func (engine *Engine) LogOut(name string) error {
	player, err := engine.World.Player(name)
	if player == nil {
		return errors.New(fmt.Sprintf(`No such player "%s"`, name))
	} else if err != nil {
		return err
	}

	if !player.LoggedIn {
		return errors.New("You are already logged out")
	}
	player.LoggedIn = false

	entity, err := engine.World.Entity(player.EntityId)
	if entity == nil {
		return errors.New(fmt.Sprintf(`Player's entity does not exist: "%d"`, player.EntityId))
	} else if err != nil {
		return err
	}
	if err := engine.World.RemoveEntityFromChunk(entity.Id, entity.Location.Chunk); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) ([]byte, error) {
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14); err == nil {
		return hashedPassword, nil
	} else {
		return nil, errors.Wrap(err, "Failed to hash password")
	}
}
