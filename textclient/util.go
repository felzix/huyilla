package main

import (
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
)

func voxelToRune(v types.Voxel) rune {
	voxel := v.Expand()

	if voxel.Form != content.FORM["cube"] {
		return rune(0)
	}

	switch voxel.Material {
	case content.MATERIAL["air"]:
		return ' '
	case content.MATERIAL["dirt"]:
		return '.'
	case content.MATERIAL["grass"]:
		return ','
	case content.MATERIAL["water"]:
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
