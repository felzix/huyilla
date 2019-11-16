package types

import (
	"fmt"
	C "github.com/felzix/huyilla/constants"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
)

type World struct {
	DB             Database
	Seed           uint64
	WorldGenerator WorldGenerator
}

func NewWorld(seed uint64, generator WorldGenerator, db Database) (*World, error) {
	world := &World{
		Seed: seed,
		DB: db,
		WorldGenerator: generator,
	}

	if !world.DB.Has(KeyAge) {
		defaultAge := NewAge(1)
		if err := world.DB.Set(KeyAge, &defaultAge); err != nil {
			return nil, err
		}
	}

	return world, nil
}

func (world *World) WipeDatabase() error {
	return world.DB.EndAll()
}

//
// Age
//

const KeyAge = "Age"

func (world *World) Age() (Age, error) {
	var age Age
	if err := world.DB.Get(KeyAge, &age); err == nil {
		return age, nil
	} else {
		return Age(0), err
	}
}

func (world *World) IncrementAge() (Age, error) {
	if age, err := world.Age(); err == nil {
		age.Increment()
		if err := world.DB.Set(KeyAge, &age); err == nil {
			return age, nil
		} else {
			return Age(0), err
		}
	} else {
		return Age(0), err
	}
}


//
// Chunk
//

func chunkKey(p Point) string {
	return fmt.Sprintf(`Chunk.%d.%d.%d`, p.X, p.Y, p.Z)
}

func (world *World) Chunk(p Point) (*Chunk, error) {
	if chunk, err := world.OnlyGetChunk(p); chunk != nil {
		return chunk, nil
	} else if err == nil {
		if chunk, err := world.GenerateChunk(p); err == nil {
			return chunk, nil
		} else {
			return nil, errors.Wrap(err, "Chunk generation failure")
		}
	} else {
		return nil, err
	}
}

func (world *World) OnlyGetChunk(p Point) (*Chunk, error) {
	var chunk Chunk
	err := world.DB.Get(chunkKey(p), &chunk)
	switch err.(type) {
	case nil: // chunk found
		return &chunk, nil
	case ThingNotFoundError: // chunk not found
		return nil, nil
	default: // something went wrong
		return nil, err
	}
}

func (world *World) CreateChunk(p Point, chunk *Chunk) error {
	return world.SetChunk(p, chunk)
}

func (world *World) SetChunk(p Point, chunk *Chunk) error {
	age, err := world.Age()
	if err != nil {
		return err
	}

	chunk.Tick = age

	return world.DB.Set(chunkKey(p), chunk)
}

func (world *World) DeleteChunk(p Point) error {
	return world.DB.End(chunkKey(p))
}

func (world *World) AddEntityToChunk(entity *Entity) error {
	if chunk, err := world.Chunk(entity.Location.Chunk); err == nil {
		chunk.Entities = append(chunk.Entities, entity.Id)
		return world.SetChunk(entity.Location.Chunk, chunk)
	} else {
		return err
	}
}

func (world *World) RemoveEntityFromChunk(entityId EntityId, p Point) error {
	chunk, err := world.Chunk(p)
	if err != nil {
		return err
	}

	entities := chunk.Entities
	for i := 0; i < len(entities); i++ {
		id := entities[i]
		if entityId == id {
			// idiomatic way of removing a list element in Go
			entities[i] = entities[len(entities)-1]
			chunk.Entities = entities[:len(entities)-1]
			break
		}
	}

	return world.SetChunk(p, chunk)
}

func (world *World) GenerateChunk(p Point) (*Chunk, error) {
	if chunkSeed, err := hashstructure.Hash(p, nil); err == nil {
		rand.Seed(int64(world.Seed * chunkSeed))
	} else {
		return nil, err
	}

	world.WorldGenerator.SetupForChunk(p)

	chunk := NewChunk(0, C.CHUNK_LENGTH)
	var x, y, z int64
	for x = 0; x < C.CHUNK_SIZE; x++ {
		for y = 0; y < C.CHUNK_SIZE; y++ {
			for z = 0; z < C.CHUNK_SIZE; z++ {
				index := (x * C.CHUNK_SIZE * C.CHUNK_SIZE) + (y * C.CHUNK_SIZE) + z
				location := NewAbsolutePoint(p.X, p.Y, p.Z, x, y, z)
				chunk.Voxels[index] = uint64(world.WorldGenerator.GenVoxel(location))
			}
		}
	}

	if err := world.SetChunk(p, chunk); err != nil {
		return nil, err
	}
	return chunk, nil
}

//
// Entity
//

func entityKey(id EntityId) string {
	return fmt.Sprintf(`Entity.%d`, id)
}

func (world *World) Entity(id EntityId) (*Entity, error) {
	var entity Entity
	err := world.DB.Get(entityKey(id), &entity)
	switch err.(type) {
	case nil: // entity found
		return &entity, nil
	case ThingNotFoundError: // entity not found
		return nil, nil
	default: // something went wrong
		return nil, err
	}
}

func (world *World) CreateEntity(type_ EntityType, playerName string, location AbsolutePoint) (*Entity, error) {
	id := world.genUniqueEntityId()
	var entity *Entity
	if playerName == "" {
		entity = NewEntity(id, type_, location)
	} else {
		entity = NewPlayerEntity(id, type_, location, playerName)
	}

	if err := world.DB.Set(entityKey(entity.Id), entity); err == nil {
		return entity, nil
	} else {
		return nil, err
	}
}

func (world *World) SetEntity(id EntityId, entity *Entity) error {
	return world.DB.Set(entityKey(id), entity)
}

func (world *World) DeleteEntity(id EntityId) error {
	return world.DB.End(entityKey(id))
}

func (world *World) EntityExists(id EntityId) bool {
	return world.DB.Has(entityKey(id))
}

func (world *World) genUniqueEntityId() EntityId {
	var id EntityId
	for {
		id = EntityId(rand.Int63())
		if !world.EntityExists(id) {
			break
		}
	}
	return id
}

//
// Item
//

func itemKey(id ItemId) string {
	return fmt.Sprintf(`Item.%d`, id)
}

func (world *World) Item(id ItemId) (*Item, error) {
	var item Item
	err := world.DB.Get(itemKey(id), &item)
	switch err.(type) {
	case nil: // entity found
		return &item, nil
	case ThingNotFoundError: // entity not found
		return nil, nil
	default: // something went wrong
		return nil, err
	}
}

func (world *World) SetItem(item *Item) error {
	return world.DB.Set(itemKey(item.Id), item)
}

func (world *World) DeleteItem(id ItemId) error {
	return world.DB.End(itemKey(id))
}

func (world *World) ItemExists(id ItemId) bool {
	return world.DB.Has(itemKey(id))
}

func (world *World) genUniqueItemId() ItemId {
	var id ItemId
	for {
		id = ItemId(rand.Int63())
		if !world.ItemExists(id) {
			break
		}
	}
	return id
}

//
// Player
//

func playerKey(name string) string {
	return fmt.Sprintf(`Player.%s`, name)
}

func playerNameFromKey(key string) string {
	// The "." (period) is not always present because it's used as the filesystem separator.
	s := strings.TrimPrefix(key, "Player")
	s = strings.TrimPrefix(s, ".")
	return s
}

func (world *World) Player(name string) (*Player, error) {
	var player Player

	err := world.DB.Get(playerKey(name), &player)
	switch err.(type) {
	case nil: // player found
		return &player, nil
	case ThingNotFoundError: // player not found
		return nil, nil
	default: // something went wrong
		return nil, err
	}
}

func (world *World) CreatePlayer(name string, password []byte, entityId EntityId, spawn AbsolutePoint) error {
	player := NewPlayer(name, password, entityId, spawn)
	return world.DB.Set(playerKey(player.Name), player)
}

func (world *World) SetPlayer(player *Player) error {
	return world.DB.Set(playerKey(player.Name), player)
}

func (world *World) DeletePlayer(name string) error {
	return world.DB.End(playerKey(name))
}

func (world *World) GetActivePlayers() ([]*PlayerDetails, error) {
	var activePlayers []*PlayerDetails

	for key := range world.DB.GetByPrefix("Player") {
		name := playerNameFromKey(key)

		if player, err := world.Player(name); player != nil {
			if len(player.Token) > 0 {
				if entity, err := world.Entity(player.EntityId); err == nil {
					activePlayers = append(activePlayers, NewPlayerDetails(player, entity))
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
