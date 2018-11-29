package main

import (
	"fmt"
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Engine struct {
	Config   types.Config
	Age      uint64
	Players  map[string]*types.Player // name -> player
	Entities map[int64]*types.Entity
	Chunks   map[types.Point]*types.Chunk
	Actions  []*types.Action // TODO locking
}

func (engine *Engine) Init(config *types.Config) error {
	// So that recipes and terrain generator can reference content by name.
	content.PopulateContentNameMaps()

	engine.Config = *config
	engine.Age = 1
	engine.Players = make(map[string]*types.Player)
	engine.Entities = make(map[int64]*types.Entity)
	engine.Chunks = make(map[types.Point]*types.Chunk)
	engine.Actions = make([]*types.Action, 0)

	return nil
}

func (engine *Engine) SignUp(name, password string) error {
	if _, ok := engine.Players[name]; ok {
		return errors.New(fmt.Sprintf(`Player "%s" already exists`, name))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return errors.Wrap(err, "Failed to hash password")
	}

	// Create new player
	defaultLocation := newAbsolutePoint(0, 0, 0, 0, 0, 0)

	entity := engine.newEntity(ENTITY["human"], name, defaultLocation)

	player := &types.Player{
		EntityId: entity.Id,
		Name:     name,
		Password: hashedPassword,
		Spawn:    defaultLocation,
		LoggedIn: false,
	}
	engine.Players[player.Name] = player

	return nil
}

func (engine *Engine) LogIn(name, password string) (*types.PlayerDetails, error) {
	player, ok := engine.Players[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf(`No such player "%s"`, name))
	}

	if bcrypt.CompareHashAndPassword(player.Password, []byte(password)) != nil {
		return nil, errors.New("Incorrect password")
	}

	if player.LoggedIn {
		return nil, errors.New("You are already logged in.")
	}
	player.LoggedIn = true

	entity := engine.Entities[player.EntityId]
	if err := engine.addEntityToChunk(entity); err != nil {
		return nil, err
	}

	return &types.PlayerDetails{Player: player, Entity: entity}, nil
}

func (engine *Engine) LogOut(name string) error {
	player, ok := engine.Players[name]
	if !ok {
		return errors.New(fmt.Sprintf(`No such player "%s"`, name))
	}

	if !player.LoggedIn {
		return errors.New("You are already logged out.")
	}
	player.LoggedIn = false

	entity := engine.Entities[player.EntityId]
	if err := engine.removeEntityFromChunk(entity); err != nil {
		return err
	}

	return nil
}
