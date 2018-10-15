package main

import (
	"github.com/felzix/huyilla/types"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	// var ticks uint64
    // setAge := &cobra.Command{
    //     Use:           "set-age",
    //     Short:         "set the world age, as number of ticks since its creation",
    //     SilenceUsage:  true,
    //     SilenceErrors: true,
    //     RunE: func(cmd *cobra.Command, args []string) error {
    //         var result types.Age
    //         params := &types.Age {
    //             Ticks: ticks,
    //         }
    //
    //         if err := CallContract(defaultContract, "SetAge", params, &result); err != nil {
    //             return err
    //         }
    //
    //         return nil
    //     },
    // }
    // setAge.Flags().Uint64VarP(&ticks, "ticks", "t", 0, "non-negative integer")
	// call.AddCommand(setAge)

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
