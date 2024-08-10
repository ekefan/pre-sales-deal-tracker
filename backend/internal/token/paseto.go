package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type TokenMaker interface {
	CreateToken(username, role string, duration time.Duration) (string, error)
	VerifyToken(token string)(*Payload, error)
}

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// Create's a new paseto maker instance...
func NewPaseto(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("SymmetricKey too short should be: %v", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// / CreateToken takes username, role and duration to create an access token
func (maker *PasetoMaker) CreateToken(username, role string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
