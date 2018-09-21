package engine

import (
    "encoding/json"
    "io/ioutil"
)

type Content struct {
    ET EntityTypes
}
type EntityTypes map[EntityType]EntityProperties
type EntityProperties struct {
    Tags              []string `json:"tags"`
    InventoryCapacity uint     `json:"inventory_capacity"`
}

func LoadContent (directory string) (*Content, error) {
    content := Content{}

    ET, err := loadEntities(directory)

    if err == nil {
        content.ET = ET
    } else {
        return &Content{}, err
    }


    return &content, nil
}

func loadEntities (directory string) (EntityTypes, error) {
    // manifest := { entity type name -> EntityType }
    rawManifest, err := ioutil.ReadFile(directory + "/entity/manifest.json")
    if err != nil {
        return nil, err
    }

    var manifest map[string]EntityType
    json.Unmarshal(rawManifest, &manifest)

    entityTypes := EntityTypes{}
    for name, eType := range manifest {
        entityTypes[eType], err = loadEntity(directory, name)
        if err != nil {
            return nil, err
        }
    }

    return entityTypes, nil
}

func loadEntity(directory string, name string) (EntityProperties, error) {
    rawEntityFile, err := ioutil.ReadFile(directory + "/entity/" + name + ".json")
    if err != nil {
        return EntityProperties{}, err
    }

    var entityProperties EntityProperties
    if err := json.Unmarshal(rawEntityFile, &entityProperties); err != nil {
        return EntityProperties{}, err
    }

    return entityProperties, nil
}