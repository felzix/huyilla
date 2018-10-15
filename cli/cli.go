package main

import (
	"github.com/felzix/huyilla/types"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
    "reflect"

    "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"
)


func getKeygenCmd() (*cobra.Command) {
	var privFile string
	keygenCmd := &cobra.Command{
		Use:           "genkey",
		Short:         "generate a public and private key pair",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, priv, err := ed25519.GenerateKey(nil)
			if err != nil {
				return errors.Wrapf(err, "Error generating key pair")
			}
			data := base64.StdEncoding.EncodeToString(priv)
			if err := ioutil.WriteFile(privFile, []byte(data), 0664); err != nil {
				return errors.Wrapf(err, "Unable to write private key")
			}
			fmt.Printf("written private key file '%s'\n", privFile)
			return nil
		},
	}
	keygenCmd.Flags().StringVarP(&privFile, "key", "k", "key", "private key file")
	return keygenCmd
}

func main() {
	defaultContract := "Huyilla"

	root := &cobra.Command{
		Use:   "huyilla",
		Short: "Huyilla",
	}

	keygen := getKeygenCmd()
	root.AddCommand(keygen)

	call := ContractCallCommand()
	root.AddCommand(call)

	getAge := &cobra.Command{
		Use:           "get-age",
		Short:         "get number of ticks since world was created",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
		    var result types.Age

			if err := StaticCallContract(defaultContract, "GetAge", &types.Nothing{}, &result); err != nil {
				return err
			}

			log.Printf(`world age is "%d"`, result.Ticks)
			return nil
		},
	}
    call.AddCommand(getAge)

    getConfig := &cobra.Command{
        Use:           "get-config",
        Short:         "get config options",
        SilenceUsage:  true,
        SilenceErrors: true,
        RunE: func(cmd *cobra.Command, args []string) error {
            var result types.Config

            if err := StaticCallContract(defaultContract, "GetConfig", &types.Nothing{}, &result); err != nil {
                return err
            }

            for k, v := range result.Options.Map {
                log.Println(v)
                log.Println(v.Value)
                log.Println(v.GetValue())
                log.Println(v.GetInt())
                log.Println(v.GetString_())
                log.Println(reflect.TypeOf(v))
                log.Println(reflect.TypeOf(v.Value))
                log.Println(reflect.TypeOf(v.GetValue()))
                log.Println(reflect.TypeOf(v.GetInt()))
                log.Println(reflect.TypeOf(v.GetString_()))
                log.Printf(`%s -> %v`, k, v.GetValue())
            }

            return nil
        },
    }
    call.AddCommand(getConfig)

	var playerCapParam int64
    setPlayerCap := &cobra.Command{
        Use:           "config-player-cap",
        Short:         "set PlayerCap config option",
        SilenceUsage:  true,
        SilenceErrors: true,
        RunE: func(cmd *cobra.Command, args []string) error {
            playerCap := &types.Primitive{Value: &types.Primitive_Int{Int: playerCapParam}}
            optionsMap := map[string]*types.Primitive{
                "PlayerCap": playerCap,
            }
            params := types.PrimitiveMap{Map: optionsMap}

            if err := CallContract(defaultContract, "SetConfigOptions", &params, &types.Nothing{}); err != nil {
                return err
            }

            return nil
        },
    }
    setPlayerCap.Flags().Int64Var(&playerCapParam, "cap", 10, "1+")
	call.AddCommand(setPlayerCap)

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
