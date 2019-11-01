package engine

import (
	"github.com/felzix/huyilla/types"
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

	if !world.DB.Has(KEY_AGE) {
		defaultAge := types.Age{Ticks: 1}
		if err := world.DB.Set(KEY_AGE, &defaultAge); err != nil {
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

const KEY_AGE = "Age"

func (world *World) Age() (*types.Age, error) {
	var age types.Age
	if err := world.DB.Get(KEY_AGE, &age); err == nil {
		return &age, nil
	} else {
		return nil, err
	}
}

func (world *World) IncrementAge() (*types.Age, error) {
	if age, err := world.Age(); err == nil {
		age.Ticks++
		if err := world.DB.Set(KEY_AGE, age); err == nil {
			return age, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
