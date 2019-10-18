package engine

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/felzix/huyilla/types"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
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
	defaultLocation := types.NewAbsolutePoint(0, 0, 0, 0, 0, 1)

	entity, err := engine.World.CreateEntity(ENTITY["human"], name, defaultLocation)
	if err != nil {
		return err
	}
	entity.PlayerName = name

	return engine.World.CreatePlayer(name, hashedPassword, entity.Id, defaultLocation)
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

	if token, err := makeToken(engine.Secret, name, time.Now().Add(time.Hour*24).Unix()); err == nil {
		originalToken := player.Token
		player.Token = token
		if len(originalToken) > 0 {
			return token, nil // player is already logged in; they just needed a new token
		}
	} else {
		return "", err
	}

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

	return player.Token, nil
}

func (engine *Engine) LogOut(name string) error {
	player, err := engine.World.Player(name)
	if player == nil {
		return errors.New(fmt.Sprintf(`No such player "%s"`, name))
	} else if err != nil {
		return err
	}

	if len(player.Token) == 0 {
		return errors.New("You are already logged out")
	}
	player.Token = ""
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

func (engine *Engine) UserExists(name string) (bool, error) {
	if player, err := engine.World.Player(name); err == nil {
		return player != nil, nil
	} else {
		return false, err
	}
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
	if id, err := uuid.NewV4(); err == nil {
		claims["tokenId"] = id
	} else {
		return "", err
	}

	if tokenString, err := token.SignedString(secret); err == nil {
		return tokenString, nil
	} else {
		return "", err
	}
}

func readToken(secret []byte, tokenString string) (string, string, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", "", 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name := claims["name"].(string)
		exp := int64(claims["exp"].(float64))
		tokenId := claims["tokenId"].(string)
		return name, tokenId, exp, nil
	} else {
		return "", "", 0, err
	}

}
