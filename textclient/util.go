package main

import (
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
)

func voxelToRune(voxel uint64) rune {
	voxelType := voxel & 0xFFFF

	switch voxelType {
	case content.VOXEL["air"]:
		return ' '
	case content.VOXEL["barren_earth"]:
		return '.'
	case content.VOXEL["barren_grass"]:
		return ','
	case content.VOXEL["water"]:
		return '~'
	default:
		return rune(0)
	}
}

func entityToRune(entity *types.Entity) rune {
	switch entity.Type {
	case content.ENTITY["human"]:
		return '@'
	case content.ENTITY["snake"]:
		return '~'
	case content.ENTITY["wisp"]:
		return '*'
	default:
		return rune(0)
	}
}
