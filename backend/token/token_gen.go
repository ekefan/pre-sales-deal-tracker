package token

import (
	"time"
)

// TokenGenerator manages token creation and verification
type TokenGenerator interface {
	// GenerateToken creates a new token for a specific user_id and duration
	GenerateToken(user_id int64, role string, duration time.Duration) (string, error)

	// Verify checks if the token is valid or not
	VerifyToken(token string)(*Payload, error)
}