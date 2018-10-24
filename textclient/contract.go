package main

import (
    "encoding/base64"
    "github.com/felzix/huyilla/types"
    "github.com/pkg/errors"
    "golang.org/x/crypto/ed25519"
    "io/ioutil"
    "log"
)


func generateKey(privFile string) error {
    _, priv, err := ed25519.GenerateKey(nil)
    if err != nil {
        return errors.Wrapf(err, "Error generating key pair")
    }
    data := base64.StdEncoding.EncodeToString(priv)
    if err := ioutil.WriteFile(privFile, []byte(data), 0664); err != nil {
        return errors.Wrapf(err, "Unable to write private key")
    }
    return nil
}

func getConfig () (map[string]interface{}, error) {
    var config types.Config

    if err := StaticCallContract("GetConfig", &types.Nothing{}, &config); err != nil {
        return nil, err
    }

    native := make(map[string]interface{})
    for k, v := range config.Options.Map {
        switch value := v.Value.(type) {
        case *types.Primitive_Int: native[k] = value.Int
        case *types.Primitive_Bool: native[k] = value.Bool
        case *types.Primitive_String_: native[k] = value.String_
        case *types.Primitive_Float: native[k] = value.Float
        default: native[k] = nil
        }
    }

    return native, nil
}

func signUp (name string) error {
    playerName := types.PlayerName{Name: name}
    if err := CallContract("SignUp", &playerName, &types.Nothing{}); err != nil {
        return err
    }
    return nil
}

func logIn (name string) (*types.PlayerDetails, error) {
    playerName := types.PlayerName{Name: name}
    var player types.PlayerDetails
    if err := CallContract("LogIn", &playerName, &player); err != nil {
        return nil, err
    }
    return &player, nil
}

func getAge () (uint64, error) {
    var age types.Age
    if err := StaticCallContract("GetAge", &types.Nothing{}, &age); err != nil {
        return 0, err
    }

    return age.Ticks, nil
}


func getChunk (point *types.Point) (*types.Chunk, error) {
    var chunk types.Chunk

    if err := StaticCallContract("GetChunk", point, &chunk); err != nil {
        return nil, err
    }

    log.Print(chunk.Voxels)

    return &chunk, nil
}
