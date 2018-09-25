package engine

import (
    "encoding/json"
    "io/ioutil"
)

type Content struct {
    E  EntityTypes
    EP EntityProperties
    M  Materials
    MP MaterialProperties
    F  Forms
    FP FormProperties
}
type EntityTypes map[string]EntityType
type EntityProperties map[EntityType]EntityProperty
type Materials map[string]Material
type MaterialProperties map[Material]MaterialProperty
type Forms map[string]Form
type FormProperties map[Form]FormProperty

type EntityProperty struct {
    Tags              []string `json:"tags"`
    InventoryCapacity uint     `json:"inventory_capacity"`
}
type MaterialProperty struct {
    Metallic bool
    Wooden bool
    Fleshy bool
    Salty bool

    PH uint
    SolidAt Kelvin
    GasAt Kelvin
}
type FormProperty struct {
    Wieldable bool
    Wearable bool
}

type Kelvin uint

func Celsius (k uint) Kelvin {
    return Kelvin(k + 273)
}


func LoadContent (directory string) (*Content, error) {
    content := Content{}

    E, EP, err := loadEntities(directory)
    if err == nil {
        content.E = E
        content.EP = EP
    } else {
        return &Content{}, err
    }

    M, MP, err := loadMaterials(directory)
    if err == nil {
        content.M = M
        content.MP = MP
    } else {
        return &Content{}, err
    }

    F, FP, err := loadForms(directory)
    if err == nil {
        content.F = F
        content.FP = FP
    } else {
        return &Content{}, err
    }

    return &content, nil
}


func loadEntities (directory string) (EntityTypes, EntityProperties, error) {
    manifest, err := getManifest(directory, "entity")
    if err != nil {
        return nil, nil, err
    }

    entityTypes := EntityTypes{}
    entityProperties := EntityProperties{}
    for t, name := range manifest {
        eType := EntityType(t + 1)
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


func loadMaterials (directory string) (Materials, MaterialProperties, error) {
    manifest, err := getManifest(directory, "material")
    if err != nil {
        return nil, nil, err
    }

    materials := Materials{}
    materialProperties := MaterialProperties{}
    for m, name := range manifest {
        material := Material(m + 1)
        materials[name] = material
        materialProperties[material], err = loadMaterial(directory, name)
        if err != nil {
            return nil, nil, err
        }
    }

    return materials, materialProperties, nil
}

func loadMaterial(directory string, name string) (MaterialProperty, error) {
    rawMaterialFile, err := ioutil.ReadFile(directory + "/material/" + name + ".json")
    if err != nil {
        return MaterialProperty{}, err
    }

    var materialProperty MaterialProperty
    if err := json.Unmarshal(rawMaterialFile, &materialProperty); err != nil {
        return MaterialProperty{}, err
    }

    return materialProperty, nil
}



func loadForms (directory string) (Forms, FormProperties, error) {
    manifest, err := getManifest(directory, "form")
    if err != nil {
        return nil, nil, err
    }

    forms := Forms{}
    formProperties := FormProperties{}
    for f, name := range manifest {
        form := Form(f + 1)
        forms[name] = form
        formProperties[form], err = loadForm(directory, name)
        if err != nil {
            return nil, nil, err
        }
    }

    return forms, formProperties, nil
}

func loadForm(directory string, name string) (FormProperty, error) {
    rawMaterialFile, err := ioutil.ReadFile(directory + "/form/" + name + ".json")
    if err != nil {
        return FormProperty{}, err
    }

    var formProperty FormProperty
    if err := json.Unmarshal(rawMaterialFile, &formProperty); err != nil {
        return FormProperty{}, err
    }

    return formProperty, nil
}




func getManifest (directory string, subdirectory string) ([]string, error) {
    rawManifest, err := ioutil.ReadFile(directory + "/" + subdirectory + "/manifest.json")
    if err != nil {
        return nil, err
    }

    var manifestWrapper map[string][]string
    json.Unmarshal(rawManifest, &manifestWrapper)

    return manifestWrapper["elements"], nil
}