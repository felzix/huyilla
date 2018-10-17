package types

import (
    "github.com/gogo/protobuf/proto"
    "testing"
)


func TestConfigProtobuf (t *testing.T) {
    optionsMap := map[string]*Primitive{
        "PlayerCap": &Primitive{Value: &Primitive_Int{Int: 10}},
    }
    config := &Config{Options: &PrimitiveMap{Map: optionsMap}}

    data, err := proto.Marshal(config)
    if err != nil {
        t.Fatal("marshaling error: ", err)
    }
    newConfig := &Config{}
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
    primitiveInt := &Primitive{Value: &Primitive_Int{Int: 12}}

    data, err := proto.Marshal(primitiveInt)

    if err != nil {
        t.Fatal("marshaling error: ", err)
    }
    newPrimitiveInt := &Primitive{}
    err = proto.Unmarshal(data, newPrimitiveInt)
    if err != nil {
        t.Fatal("unmarshaling error: ", err)
    }

    if primitiveInt.GetInt() != newPrimitiveInt.GetInt() {
        t.Fatalf("data mismatch %v != %v", primitiveInt.GetInt(), newPrimitiveInt.GetInt())
    }
}

func TestAgeProtobuf (t *testing.T) {
    age := &Age{Ticks: 3}

    data, err := proto.Marshal(age)

    if err != nil {
        t.Fatal("marshaling error: ", err)
    }
    newAge := &Age{}
    err = proto.Unmarshal(data, newAge)
    if err != nil {
        t.Fatal("unmarshaling error: ", err)
    }

    if age.GetTicks() != newAge.GetTicks() {
        t.Fatalf("data mismatch %v != %v", age.GetTicks(), newAge.GetTicks())
    }
}
