package client

import (
	"bytes"
	"fmt"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

type API struct {
	Base     string
	Username string
	password []byte // made private for the meager security that offers
	token    []byte // same
}

func NewAPI(base, username, password string) *API {
	return &API{Base: base, Username: username, password: []byte(password)}
}

func (api *API) MakeAuth() ([]byte, error) {
	return (&types.Auth{Name: api.Username, Password: api.password}).Marshal()
}

func (api *API) Request(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	fullUrl := api.Base + url
	req, _ := http.NewRequest(method, fullUrl, body)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (api *API) Ping() (string, error) {
	res, err := api.Request("GET", "/ping", nil, nil)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ping failure: %v", err))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Ping failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf(`Ping failure: Expected status 200 but got %d. %s`, res.StatusCode, body))
	}

	return string(body), nil
}

func (api *API) Signup() error {
	auth, err := api.MakeAuth()
	if err != nil {
		return err
	}

	res, err := api.Request("POST", "/auth/signup", bytes.NewReader(auth), map[string]string{
		"Content-Type": "application/protobuf",
	})
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

func (api *API) Login() error {
	auth, err := api.MakeAuth()
	if err != nil {
		return err
	}

	res, err := api.Request("POST", "/auth/login", bytes.NewReader(auth), map[string]string{
		"Content-Type": "application/protobuf",
	})
	if err != nil {
		return errors.New(fmt.Sprintf("Login failure: %v", err))
	}

	token, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Login failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf(`Login failure: Expected status 200 but got %d. %s`, res.StatusCode, token))
	}

	api.token = token
	return nil
}

func (api *API) Logout() error {
	if _, err := api.Request("POST", "/auth/logout", nil, map[string]string{
		"token": string(api.token),
	}); err != nil {
		return errors.New(fmt.Sprintf("Logout failure: %v", err))
	}

	api.token = nil
	return nil
}

func (api *API) UserExists() (bool, error) {
	res, err := api.Request("GET", "/auth/exists/"+api.Username, nil, nil)
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

func (api *API) GetPlayer(name string) (*types.Entity, error) {
	res, err := api.Request("GET", "/world/player/"+name, nil, map[string]string{
		"Content-Type": "application/protobuf",
		"token":        string(api.token),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("GetPlayer failure: %v", err))
	}

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("GetPlayer failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf(`GetPlayer failure: Expected status 200 but got %d. %s`, res.StatusCode, blob))
	}

	var entity types.Entity
	if err := entity.Unmarshal(blob); err != nil {
		return nil, errors.New(fmt.Sprintf(`GetPlayer failure: Malformed protobuf blob: %v`, err))
	}

	return &entity, nil
}

func (api *API) GetWorldAge() (types.Age, error) {
	res, err := api.Request("GET", "/world/age", nil, nil)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("GetWorldAge failure: %v", err))
	}

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("GetWorldAge failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return 0, errors.New(fmt.Sprintf(`GetWorldAge failure: Expected status 200 but got %d. %s`, res.StatusCode, blob))
	}

	var age types.Age
	if err := age.Unmarshal(blob); err != nil {
		return types.Age(0), errors.New(fmt.Sprintf("GetWorldAge failure: %v", err))
	}
	return age, nil
}

func (api *API) GetChunks(point types.Point, radius uint64) (*types.Chunks, error) {
	res, err := api.Request("GET", fmt.Sprintf("/world/chunk/%d/%d/%d?radius=%d", point.X, point.Y, point.Z, radius), nil, map[string]string{
		"Content-Type": "application/protobuf",
		"token":        string(api.token),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("GetChunk failure: %v", err))
	}

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("GetChunk failure: %v", err))
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf(`GetChunk failure: Expected status 200 but got %d. %s`, res.StatusCode, blob))
	}

	var chunks types.Chunks
	if err := chunks.Unmarshal(blob); err != nil {
		return nil, errors.New(fmt.Sprintf(`GetChunk failure: Malformed protobuf blob: %v`, err))
	}

	return &chunks, nil
}

func (api *API) IssueAction(action *types.Action) error {
	blob, err := action.Marshal()
	if err != nil {
		return errors.New(fmt.Sprintf("GetChunk failure: %v", err))
	}

	_, err = api.Request("POST", fmt.Sprintf("/world/act"), bytes.NewReader(blob), map[string]string{
		"token":        string(api.token),
	})
	if err != nil {
		return errors.New(fmt.Sprintf("GetChunk failure: %v", err))
	}

	return nil
}

func (api *API) IssueMoveAction(player *types.Player, whereTo types.AbsolutePoint) error {
	return api.IssueAction(types.NewMoveAction(player.EntityId, player.EntityId, whereTo))
}
