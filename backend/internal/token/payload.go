package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Payload holds fields required to hold data for JWT authentication
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	IssAt    time.Time `json:"created_at"`
	ExpAt    time.Time `json:"expiry_at"`
}


// NewPayload takes username, role and duration 
// to create a new instance of a payload
func NewPayload(username, role string, duration time.Duration) (*Payload, error){
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID: tokenId,
		Username: username,
		Role: role,
		IssAt: time.Now(),
		ExpAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token is not expired
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpAt) {
		return errors.New("token has expired")
	}
	return nil
}