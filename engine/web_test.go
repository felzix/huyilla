package main

import (
	"bytes"
	"fmt"
	"github.com/felzix/huyilla/types"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHuyilla_Web_Ping(t *testing.T) {
	defer func() { http.DefaultServeMux = new(http.ServeMux) }()
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()
	web := httptest.NewServer(pingHandler(h))
	defer web.Close()

	res, err := http.Get(web.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatal(fmt.Sprintf(`Expected status 200 but got %d`, res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "pong" {
		t.Errorf(`Expected "pong" but got "%s"`, string(body))
	}
}

func TestHuyilla_Web_Signup_flow(t *testing.T) {
	defer func() { http.DefaultServeMux = new(http.ServeMux) }()
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()

	NAME := "felzix"
	PASS := "murakami"

	auth, err := (&types.Auth{Name: NAME, Password: []byte(PASS)}).Marshal()
	if err != nil {
		t.Fatal(err)
	}

	// Signup

	web_signup := httptest.NewServer(signupHandler(h))
	defer web_signup.Close()


	res, err := http.Post(web_signup.URL+"/auth/signup", "application/protobuf", bytes.NewReader(auth))
	if err != nil {
		t.Fatal(err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatal(fmt.Sprintf(`Expected status 200 but got %d`, res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "Signup successful!" {
		t.Errorf(`Expected "Signup successful!" but got "%s"`, string(body))
	}

	// Verify database

	player, err := h.World.Player(NAME)
	if player == nil {
		t.Fatalf("Player does not exist")
	} else if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if player.LoggedIn {
		t.Error("Player should not be logged-in just because they signed up")
	}

	if player.Name != "felzix" {
		t.Errorf(`Player name was "%v" instead of "felzix"`, player.Name)
	}

	// Login

	web_login := httptest.NewServer(loginHandler(h))
	defer web_login.Close()

	res, err = http.Post(web_login.URL+"/auth/login", "application/protobuf", bytes.NewReader(auth))
	if err != nil {
		t.Fatal(err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatal(fmt.Sprintf(`Expected status 200 but got %d`, res.StatusCode))
	}

	token, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(token) < 100 || token[0] != 'e' {
		t.Errorf(`Bad token. Token="%s"`, token)
	}

	// Verify Database

	player, err = h.World.Player(NAME)
	if player == nil {
		t.Fatalf("Player does not exist")
	} else if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if !player.LoggedIn {
		t.Error("Player should be logged-in")
	}

	// Logout

	web_logout := httptest.NewServer(logoutHandler(h))
	defer web_logout.Close()

	req, err := http.NewRequest("POST", web_logout.URL+"/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("token", string(token))
	req.Header.Add("content-type", "application/protobuf")
	client := http.Client{}
	res, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	} else if res.StatusCode != http.StatusOK {
		t.Fatal(fmt.Sprintf(`Expected status 200 but got %d`, res.StatusCode))
	}

	// Verify Database

	player, err = h.World.Player(NAME)
	if player == nil {
		t.Fatalf("Player does not exist")
	} else if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if player.LoggedIn {
		t.Error("Player should be logged-out")
	}
}
