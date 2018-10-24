package main

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
)


const (
    ADDR = "Huyilla"
    CHAIN_ID = "default"
    WRITE_URI = "http://localhost:46658/rpc"
    READ_URI = "http://localhost:46658/query"
    PRIV_FILE = "key"
)


func contract() (*client.Contract, error) {
	contractAddr, err := ResolveAddress(ADDR)
	if err != nil {
		return nil, err
	}

	// create rpc client
	rpcClient := client.NewDAppChainRPCClient(
		CHAIN_ID,
		WRITE_URI,
		READ_URI)
	// create contract
	contract := client.NewContract(rpcClient, contractAddr.Local)
	return contract, nil
}

func CallContract(method string, params proto.Message, result interface{}) error {
	privKeyB64, err := ioutil.ReadFile(PRIV_FILE)
	if err != nil {
		return err
	}

	privKey, err := base64.StdEncoding.DecodeString(string(privKeyB64))
	if err != nil {
		return err
	}

	signer := auth.NewEd25519Signer(privKey)

	contract, err := contract()
	if err != nil {
		return err
	}
	_, err = contract.Call(method, params, signer, result)
	return err
}

func StaticCallContract(method string, params proto.Message, result interface{}) error {
	contract, err := contract()
	if err != nil {
		return err
	}

	_, err = contract.StaticCall(method, params, loom.RootAddress(CHAIN_ID), result)
	return err
}
