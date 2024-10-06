package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtGenerator struct {
	secretKey []byte
}


const minSecretKeyLen = 32
func NewJwtGenerator(secretKey string) (TokenGenerator, error) {
	if len(secretKey) < minSecretKeyLen {
		return nil, fmt.Errorf("length of secret key too short, min length: %v", minSecretKeyLen)
	}
	return &JwtGenerator{
		secretKey: []byte(secretKey),
	}, nil
}

func (jg *JwtGenerator) GenerateToken(user_id int64, role string, duration time.Duration) (string, error) {
	payload,  err := NewPayload(user_id, role, duration)
	if err !=  nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(jg.secretKey)
}

func (jg *JwtGenerator) VerifyToken(token string)(*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return jg.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
