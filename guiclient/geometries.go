package main

import (
	"github.com/felzix/huyilla/content"
	"github.com/felzix/huyilla/types"
	"github.com/g3n/engine/geometry"
)

var geometries = make(map[types.Form]geometry.IGeometry)

func makeGeometries() {
	geometries[content.FORM["cube"]] = geometry.NewCube(1)
	geometries[content.FORM["human"]] = geometry.NewCylinder(0.3, 1.8, 16, 1, true, true)
}
