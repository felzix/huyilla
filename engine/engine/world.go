package engine

import (
	"github.com/felzix/huyilla/types"
	"github.com/gogo/protobuf/proto"
	"github.com/peterbourgon/diskv"
)

type World struct {
	DB             *diskv.Diskv
	Seed           uint64
	WorldGenerator WorldGenerator
}

func (world *World) Init(saveDir string, cacheSize uint64) error {
	if db, err := makeDB(saveDir, cacheSize); err == nil {
		world.DB = db
	} else {
		return err
	}

	if !world.DB.Has(KEY_AGE) {
		defaultAge := types.Age{Ticks: 1}
		if blob, err := proto.Marshal(&defaultAge); err == nil {
			if err := world.DB.Write(KEY_AGE, blob); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	world.WorldGenerator = NewLakeWorldGenerator(3)

	return nil
}

func (world *World) WipeDatabase() error {
	return world.DB.EraseAll()
}

//
// Age
//

const KEY_AGE = "Age"

func (world *World) Age() (*types.Age, error) {
	var age types.Age
	if err := gettum(world, KEY_AGE, &age); err == nil {
		return &age, nil
	} else {
		return nil, err
	}
}

func (world *World) IncrementAge() (*types.Age, error) {
	if age, err := world.Age(); err == nil {
		age.Ticks++
		if err := settum(world, KEY_AGE, age); err == nil {
			return age, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
