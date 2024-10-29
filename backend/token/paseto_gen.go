package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoGenerator is a PASETO token maker
type PasetoGenerator struct {
	pasetoVersion *paseto.V2
	symmetricKey  []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoGenerator(symmetricKey string) (TokenGenerator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	tokenGenerator := &PasetoGenerator{
		pasetoVersion: paseto.NewV2(),
		symmetricKey:  []byte(symmetricKey),
	}

	return tokenGenerator, nil
}

// CreateToken creates a new token for a specific username and duration
func (pg *PasetoGenerator) GenerateToken(user_id int64, role string, duration time.Duration) (string, error) {
	payload, err := NewPayload(user_id, role, duration)
	if err != nil {
		return "", err
	}
	token, err := pg.pasetoVersion.Encrypt(pg.symmetricKey, payload, nil)
	return token, err
}

// VerifyToken checks if the token is valid or not
func (pg *PasetoGenerator) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := pg.pasetoVersion.Decrypt(token, pg.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
