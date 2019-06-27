package main

import (
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
)

func voxelToRune(v types.Voxel) rune {
	F := content.FORM
	M := content.MATERIAL

	voxel := v.Expand()

	if voxel.Form != F["cube"] {
		return rune(0)
	}

	switch voxel.Material {
	case M["air"]:
		return ' '
	case M["dirt"]:
		return '.'
	case M["grass"]:
		return ','
	case M["water"]:
		return '~'
	default:
		return rune(0)
	}
}

func entityToRune(entity *types.Entity) rune {
	E := content.ENTITY

	switch entity.Type {
	case E["human"]:
		return '@'
	case E["snake"]:
		return '~'
	case E["wisp"]:
		return '*'
	default:
		return rune(0)
	}
}
