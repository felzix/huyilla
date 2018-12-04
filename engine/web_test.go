package main

import (
	"bytes"
	"github.com/felzix/huyilla/types"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHuyilla_Web_Ping(t *testing.T) {
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()
	web := httptest.NewServer(pingHandler(h))
	defer web.Close()
	defer func() { http.DefaultServeMux = new(http.ServeMux) }()

	res, err := http.Get(web.URL + "/ping")
	if err != nil {
		t.Fatal(err)
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
	h := &Engine{}
	h.Init()
	defer h.World.WipeDatabase()
	web := httptest.NewServer(signupHandler(h))
	defer web.Close()
	defer func() { http.DefaultServeMux = new(http.ServeMux) }()

	NAME := "felzix"
	PASS := "murakami"

	auth := types.Auth{Name: NAME, Password: []byte(PASS)}
	blob, err := auth.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.Post(web.URL+"/auth/signup", "application/protobuf", bytes.NewReader(blob))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)

	var body []byte
	if _, err := res.Body.Read(body); err != nil {
		t.Fatal(err)
	}
	t.Log(body)
	if string(body) != "Signup successful!" {
		t.Errorf(`Expected "Signup successful!" but got "%s"`, string(body))
	}

	player, err := h.World.Player(NAME)
	if player == nil {
		t.Fatalf("Player does not exist")
	} else if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if player.Name != "felzix" {
		t.Errorf(`Player name was "%v" instead of "felzix"`, player.Name)
	}
}
