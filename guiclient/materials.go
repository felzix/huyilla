package main

import (
	"github.com/felzix/huyilla/content"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

var materials = make(map[uint64]material.IMaterial)

func makeMaterials() {
	materials[content.MATERIAL["dirt"]] = material.NewStandard(math32.NewColor("SaddleBrown"))
	materials[content.MATERIAL["grass"]] = material.NewStandard(math32.NewColor("SpringGreen"))
	materials[content.MATERIAL["water"]] = material.NewStandard(math32.NewColor("DarkBlue"))

	materials[content.MATERIAL["human"]] = material.NewStandard(math32.NewColor("DarkRed"))
}
