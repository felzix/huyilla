package engine

import (
    "path"
    "runtime"
    "testing"
)

func TestLoadContent(t *testing.T) {
    directory := getDirectory(t)

    content, err := LoadContent(directory)
    if err != nil {
        t.Fatal(err)
    }

    if content.ET[1].InventoryCapacity != 10 {
        t.Error("Entity inventory capacity was incorrect or not present:", content)
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


func TestLoadEntities(t *testing.T) {
    directory := getDirectory(t)

    entityTypes, err := loadEntities(directory)

    if err != nil {
        t.Fatal(err)
    }

    // 1 -> player
    if entityTypes[1].InventoryCapacity != 10 {
        t.Error("Entity inventory capacity was incorrect or not present:", entityTypes)
    }
}

func getDirectory (t *testing.T) string {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        t.Fatal("Failed to discover current directory")
    }

    return path.Dir(filename) + "/../content"
}