package main

import (
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/client"
)

func ParseBytes(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		return hex.DecodeString(s[2:])
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		b, err = base64.StdEncoding.DecodeString(s)
	}

	return b, err
}

func ParseAddress(s string) (loom.Address, error) {
	addr, err := loom.ParseAddress(s)
	if err == nil {
		return addr, nil
	}

	b, err := ParseBytes(s)
	if len(b) != 20 {
		return loom.Address{}, loom.ErrInvalidAddress
	}

	return loom.Address{ChainID: CHAIN_ID, Local: loom.LocalAddress(b)}, nil
}

func ResolveAddress(s string) (loom.Address, error) {
	rpcClient := client.NewDAppChainRPCClient(CHAIN_ID, WRITE_URI, READ_URI)
	contractAddr, err := ParseAddress(s)
	if err != nil {
		// if address invalid, try to resolve it using registry
		contractAddr, err = rpcClient.Resolve(s)
		if err != nil {
			return loom.Address{}, err
		}
	}

	return contractAddr, nil
}
