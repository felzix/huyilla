package engine

import (
    "testing"
)

func TestSimpleItem(t *testing.T) {
    content := getContent(t)

    item := MakeSimpleItem(content.F["pile"], content.M["dirt"])

    if form := item.Form(); form != 2 {
        t.Errorf(`Item form is "%v" instead of "%d"`, form, 2)
    }

    if material := item.Material(); material != 2 {
        t.Errorf(`Item material is "%v" instead of "%d"`, material, 2)
    }
}


func TestComplexItem(t *testing.T) {
    content := getContent(t)

    shaft := MakeSimpleItem(content.F["pole"], content.M["oakwood"])
    head := MakeSimpleItem(content.F["spearhead"], content.M["copper"])

    spear := MakeComplexItem(content.F["spear"], head, shaft)

    if form := spear.Form(); form != 12 {
        t.Errorf(`Item form is "%v" instead of "%d"`, form, 12)
    }

    if components := spear.components; len(components) == 2 {
        if components[0] != head {
            t.Errorf(`Spear head is actually "%v"`, components[0])
        }
        if components[1] != shaft {
            t.Errorf(`Spear shaft is actually "%v"`, components[1])
        }
    } else {
        t.Errorf(`Item has "%d" components instead of "%d"`, len(components), 2)
    }
}
