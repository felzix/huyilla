package main

import (
    "encoding/base64"
    "fmt"
    "github.com/felzix/huyilla/types"
    "github.com/pkg/errors"
    "github.com/spf13/cobra"
    "golang.org/x/crypto/ed25519"
    "io/ioutil"
    "log"
    "os"
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
                out := fmt.Sprintf(`%s -> `, k)
            	switch value := v.Value.(type) {
                case *types.Primitive_Int: out += fmt.Sprint(value.Int)
                case *types.Primitive_Bool: out += fmt.Sprint(value.Bool)
                case *types.Primitive_String_: out += value.String_
                case *types.Primitive_Float: out += fmt.Sprint(value.Float)
                default: out = fmt.Sprintf(`ERROR: unrecognized type for "%v"="%v"`, k, v)
                }
                log.Print(out)
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

    var x, y, z int64
    getChunk := &cobra.Command{
        Use:           "get-chunk",
        Short:         "get chunk",
        SilenceUsage:  true,
        SilenceErrors: true,
        RunE: func(cmd *cobra.Command, args []string) error {
            var chunk types.Chunk
            point := types.Point{x, y, z}

            if err := StaticCallContract(defaultContract, "GetChunk", &point, &chunk); err != nil {
                return err
            }

            log.Print(chunk.Voxels)

            return nil
        },
    }
    getChunk.Flags().Int64VarP(&x, "x", "x", 0, "int64")
    getChunk.Flags().Int64VarP(&y, "y", "y", 0, "int64")
    getChunk.Flags().Int64VarP(&z, "z", "z", 0, "int64")
    call.AddCommand(getChunk)

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}
