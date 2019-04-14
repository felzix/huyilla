package main

import "github.com/felzix/huyilla/content"

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
