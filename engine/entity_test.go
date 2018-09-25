package engine

import (
    "testing"
)

func TestTypeProperty(t *testing.T) {
    content := getContent(t)
    entity := MakeEntity(content.E["player"])

    invCap := entity.InventoryCapacity(content)

    if invCap != 10 {
        t.Error("Player inventory capacity should be 10 but is", invCap)
    }
}
