package engine

import (
    "testing"
)

func TestLoadContent(t *testing.T) {
    directory := getDirectory(t)

    content, err := LoadContent(directory)
    if err != nil {
        t.Fatal(err)
    }

    if content.E["player"] != 1 {
        t.Error("Player entity has wrong type or is not present:", content)
    }

    if content.EP[1].InventoryCapacity != 10 {
        t.Error("Entity inventory capacity was incorrect or not present:", content)
    }
}


func TestLoadEntities(t *testing.T) {
    directory := getDirectory(t)

    entityTypes, entityProperties, err := loadEntities(directory)

    if err != nil {
        t.Fatal(err)
    }


    if entityTypes["player"] != 1 {
        t.Error("Player entity has wrong type or is not present:", entityTypes)
    }

    if entityProperties[1].InventoryCapacity != 10 {
        t.Error("Entity inventory capacity was incorrect or not present:", entityProperties)
    }
}


func TestLoadEntity(t *testing.T) {
    directory := getDirectory(t)

    entityProperties, err := loadEntity(directory, "player")
    if err != nil {
        t.Fatal(err)
    }

    if entityProperties.InventoryCapacity != 10 {
        t.Error("Entity inventory capacity was incorrect or not present:", entityProperties)
    }
}
