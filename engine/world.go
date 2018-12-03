package main

import (
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/gogo/protobuf/proto"
	"github.com/peterbourgon/diskv"
	"github.com/satori/go.uuid"
	"strings"
)

type World struct {
	DB *diskv.Diskv

	Players  map[string]*types.Player // name -> player
	Entities map[int64]*types.Entity
	Chunks   map[types.Point]*types.Chunk
}

func (world *World) Init(saveDir string, cacheSize uint64) error {
	// So that recipes and terrain generator can reference content by name.
	content.PopulateContentNameMaps()

	unique, err := uuid.NewV4()
	if err != nil {
		return err
	}

	world.DB = diskv.New(diskv.Options{
		BasePath: saveDir,
		TempDir: "/tmp/tempdir-huyilla-" + unique.String(),
		AdvancedTransform: filesystemTransform,
		InverseTransform: filesystemInverseTransform,
		CacheSizeMax: cacheSize,

	})

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

func filesystemTransform(key string) *diskv.PathKey {
	path := strings.Split(key, ".")
	last := len(path) - 1
	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last],
	}
}

func filesystemInverseTransform(pathKey *diskv.PathKey) (key string) {
	return strings.Join(pathKey.Path, "/") + pathKey.FileName
}

const (
	KEY_AGE = "Age"
)

func (world *World) Age() (*types.Age, error) {
	world.DB.Has(KEY_AGE)
	if blob, err := world.DB.Read(KEY_AGE); err == nil {
		var age types.Age
		if err := proto.Unmarshal(blob, &age); err == nil {
			return &age, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (world *World) IncrementAge() error {
	if age, err := world.Age(); err == nil {
		age.Ticks++
		if blob, err := proto.Marshal(age); err == nil {
			if err := world.DB.Write(KEY_AGE, blob); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}
