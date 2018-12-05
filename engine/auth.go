package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (engine *Engine) SignUp(name string, password string) error {
	if player, err := engine.World.Player(name); player != nil {
		return errors.New(fmt.Sprintf(`Player "%s" already exists`, name))
	} else if err != nil {
		return err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	// Create new player
	defaultLocation := newAbsolutePoint(0, 0, 0, 0, 0, 0)

	entity, err := engine.World.CreateEntity(ENTITY["human"], name, defaultLocation)
	if err != nil {
		return err
	}

	engine.World.CreatePlayer(name, hashedPassword, entity.Id, defaultLocation)

	return nil
}

func (engine *Engine) LogIn(name, password string) (string, error) {
	player, err := engine.World.Player(name)
	if player == nil {
		return "", errors.New(fmt.Sprintf(`No such player "%s"`, name))
	} else if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword(player.Password, []byte(password)) != nil {
		return "", errors.New("Incorrect password")
	}

	if player.LoggedIn {
		return "", errors.New("You are already logged in.")
	}
	player.LoggedIn = true
	if err := engine.World.SetPlayer(player); err != nil {
		return "", err
	}

	entity, err := engine.World.Entity(player.EntityId)
	if entity == nil {
		return "", errors.New(fmt.Sprintf(`Player's entity does not exist: "%d"`, player.EntityId))
	} else if err != nil {
		return "", err
	}
	if err := engine.World.AddEntityToChunk(entity); err != nil {
		return "", err
	}

	return makeToken(engine.Secret, name, time.Now().Add(time.Hour * 24).Unix())
}

func (engine *Engine) LogOut(name string) error {
	player, err := engine.World.Player(name)
	if player == nil {
		return errors.New(fmt.Sprintf(`No such player "%s"`, name))
	} else if err != nil {
		return err
	}

	if !player.LoggedIn {
		return errors.New("You are already logged out")
	}
	player.LoggedIn = false
	if err := engine.World.SetPlayer(player); err != nil {
		return err
	}

	entity, err := engine.World.Entity(player.EntityId)
	if entity == nil {
		return errors.New(fmt.Sprintf(`Player's entity does not exist: "%d"`, player.EntityId))
	} else if err != nil {
		return err
	}
	if err := engine.World.RemoveEntityFromChunk(entity.Id, entity.Location.Chunk); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) ([]byte, error) {
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14); err == nil {
		return hashedPassword, nil
	} else {
		return nil, errors.Wrap(err, "Failed to hash password")
	}
}

func makeToken(secret []byte, name string, expiry int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["exp"] = expiry

	if tokenString, err := token.SignedString(secret); err == nil {
		return tokenString, nil
	} else {
		return "", err
	}
}

func readToken(secret []byte, tokenString string) (string, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name := claims["name"].(string)
		exp := int64(claims["exp"].(float64))
		return name, exp, nil
	} else {
		return "", 0, err
	}

}