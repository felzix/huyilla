package main

import (
    "github.com/felzix/huyilla/types"
    "github.com/gogo/protobuf/proto"
    "testing"
)

func TestHuyilla_GetConfig(t *testing.T) {
    meta, err := Contract.Meta()
    if err != nil {
        t.Errorf(`Error: %v`, err)
    }

    if meta.Name != "Huyilla" {
        t.Errorf(`Contract name is "%v"; should be "Huyilla"`, meta.Name)
    }
    if meta.Version != "0.0.1" {
        t.Errorf(`Contract version is "%v"; should be "0.0.1"`, meta.Version)
    }
}

func TestConfigProtobuf (t *testing.T) {
    optionsMap := map[string]*types.Primitive{
        "PlayerCap": &types.Primitive{Value: &types.Primitive_Int{Int: 10}},
    }
    config := &types.Config{Options: &types.PrimitiveMap{Map: optionsMap}}

    data, err := proto.Marshal(config)
    if err != nil {
        t.Fatal("marshaling error: ", err)
    }
    newConfig := &types.Config{}
    err = proto.Unmarshal(data, newConfig)
    if err != nil {
        t.Fatal("unmarshaling error: ", err)
    }
    playerCap := config.GetOptions().GetMap()["PlayerCap"].GetInt()
    newPlayerCap := newConfig.GetOptions().GetMap()["PlayerCap"].GetInt()
    if playerCap != newPlayerCap {
        t.Fatalf("data mismatch %v != %v", playerCap, newPlayerCap)
    }
}

func TestPrimitiveProtobuf (t *testing.T) {
    primitiveInt := &types.Primitive{Value: &types.Primitive_Int{Int: 12}}

    data, err := proto.Marshal(primitiveInt)

    if err != nil {
        t.Fatal("marshaling error: ", err)
    }
    newPrimitiveInt := &types.Primitive{}
    err = proto.Unmarshal(data, newPrimitiveInt)
    if err != nil {
        t.Fatal("unmarshaling error: ", err)
    }

    if primitiveInt.GetInt() != newPrimitiveInt.GetInt() {
        t.Fatalf("data mismatch %v != %v", primitiveInt.GetInt(), newPrimitiveInt.GetInt())
    }
}

func TestAgeProtobuf (t *testing.T) {
    age := &types.Age{Ticks: 3}

    data, err := proto.Marshal(age)

    if err != nil {
        t.Fatal("marshaling error: ", err)
    }
    newAge := &types.Age{}
    err = proto.Unmarshal(data, newAge)
    if err != nil {
        t.Fatal("unmarshaling error: ", err)
    }

    if age.GetTicks() != newAge.GetTicks() {
        t.Fatalf("data mismatch %v != %v", age.GetTicks(), newAge.GetTicks())
    }
}
