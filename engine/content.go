package engine

import (
    "encoding/json"
    "io/ioutil"
)

type Content struct {
    ET EntityTypes
    EP EntityProperties
}
type EntityTypes map[string]EntityType
type EntityProperties map[EntityType]EntityProperty

type EntityProperty struct {
    Tags              []string `json:"tags"`
    InventoryCapacity uint     `json:"inventory_capacity"`
}



func LoadContent (directory string) (*Content, error) {
    content := Content{}

    ET, EP, err := loadEntities(directory)
    if err == nil {
        content.ET = ET
        content.EP = EP
    } else {
        return &Content{}, err
    }


    return &content, nil
}

func loadEntities (directory string) (EntityTypes, EntityProperties, error) {
    // manifest := { entity type name -> EntityType }
    rawManifest, err := ioutil.ReadFile(directory + "/entity/manifest.json")
    if err != nil {
        return nil, nil, err
    }

    var manifest map[string]EntityType
    json.Unmarshal(rawManifest, &manifest)

    entityTypes := EntityTypes{}
    entityProperties := EntityProperties{}
    for name, eType := range manifest {
        entityTypes[name] = eType
        entityProperties[eType], err = loadEntity(directory, name)
        if err != nil {
            return nil, nil, err
        }
    }

    return entityTypes, entityProperties, nil
}

func loadEntity(directory string, name string) (EntityProperty, error) {
    rawEntityFile, err := ioutil.ReadFile(directory + "/entity/" + name + ".json")
    if err != nil {
        return EntityProperty{}, err
    }

    var entityProperty EntityProperty
    if err := json.Unmarshal(rawEntityFile, &entityProperty); err != nil {
        return EntityProperty{}, err
    }

    return entityProperty, nil
}