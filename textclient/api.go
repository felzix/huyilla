package main

import (
	"bytes"
	"fmt"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type API struct {
	Base string
	Username string
	password []byte // made private for the meager security that offers
}

func NewAPI(base, username, password string) *API {
	return &API{Base: base, Username: username, password: []byte(password)}
}

func (api *API) Auth() ([]byte, error) {
	return (&types.Auth{Name: api.Username, Password: api.password}).Marshal()
}

func (api *API) Signup() error {
	auth, err := api.Auth()
	if err != nil {
		return err
	}
	res, err := http.Post(api.Base + "/auth/signup", "application/protobuf", bytes.NewReader(auth))
	if err != nil {
		return errors.New(fmt.Sprintf("Signup failure: %v", err))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Signup failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf(`Signup failure: Expected status 200 but got %d. %s`, res.StatusCode, body))
	}

	return nil
}

func (api *API) Login() ([]byte, error) {
	auth, err := api.Auth()
	if err != nil {
		return nil, err
	}
	res, err := http.Post(api.Base + "/auth/login", "application/protobuf", bytes.NewReader(auth))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Login failure: %v", err))
	}

	token, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Login failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf(`Login failure: Expected status 200 but got %d. %s`, res.StatusCode, token))
	}

	return token, nil
}

func (api *API) UserExists() (bool, error) {
	res, err := http.Post(api.Base + "/auth/exists", "application/protobuf", bytes.NewReader([]byte(api.Username)))
	if err != nil {
		return false, errors.New(fmt.Sprintf("UserExists failure: %v", err))
	}

	rawExists, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, errors.New(fmt.Sprintf("UserExists failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return false, errors.New(fmt.Sprintf(`UserExists failure: Expected status 200 but got %d. %s`, res.StatusCode, rawExists))
	}

	if string(rawExists) == "true" {
		return true, nil
	} else if string(rawExists) == "false" {
		return false, nil
	} else {
		return false, errors.New(fmt.Sprintf(`UserExists failure: Expected true or false but got: %v`, rawExists))
	}
}

func (api *API) GetPlayer(name string) (*types.Player, error) {
	res, err := http.Get(api.Base + "/player/" + name)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("GetPlayer failure: %v", err))
	}

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("GetPlayer failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf(`GetPlayer failure: Expected status 200 but got %d. %s`, res.StatusCode, blob))
	}

	var player types.Player
	if err := player.Unmarshal(blob); err != nil {
		return nil, errors.New(fmt.Sprintf(`GetPlayer failure: Malformed protobuf blob: %v`, err))
	}

	return &player, nil
}
