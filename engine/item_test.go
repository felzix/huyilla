package engine

import (
    "testing"
)

func TestPhysicalFeatures(t *testing.T) {
    content := getContent(t)

    item := MakeSimpleItem(content.F["pile"], content.M["dirt"])

    if form := item.Material(); form != 2 {
        t.Errorf(`Item form is "%v" instead of "%v"`, form, 2)
    }

    if material := item.Material(); material != 2 {
        t.Errorf(`Item material is "%v" instead of "%v"`, material, 2)
    }
}
