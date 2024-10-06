package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Custom token errors
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of a token
type Payload struct {
	ID     uuid.UUID `json:"id"`
	UserID int64     `json:"user_id"`
	// Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(userID int64, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:     tokenID,
		UserID: userID,
		// Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is expired or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// Implementing the jwt CLaims
func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: p.ExpiredAt,
	}, nil
}
func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: p.IssuedAt,
	}, nil
}
func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: p.IssuedAt,
	}, nil
}
func (p *Payload) GetIssuer() (string, error) {
	return "", nil
}
func (p *Payload) GetSubject() (string, error) {
	return "", nil
}
func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{""}, nil
}
