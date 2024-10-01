package api

import (
	"context"
	"log"

	db "github.com/ekefan/pre-sales-deal-tracker/backend/db/sqlc"
)

const (
	initUsername     = "josh"
	initUserRole     = "admin"
	initUserFullname = "Joshua Olufinlua"
	initUserEmail    = "josh@vastech.ng.com"
)

// InitUser checks if an admin user exists, on error, sends
// If none, creates the initial user,
// on error creating the user, stops server with exist status 1
func initUser(store db.Store) {
	ctx := context.Background()
	numAdminUsers, err := store.GetNumberOfAdminUsers(ctx, initUserRole)
	if err != nil {
		log.Fatalf("can not determine number of admin users %v", err)
	}
	if numAdminUsers < 1 {
		hash, err := HashPassword(db.DefaultUserPassword)
		if err != nil {
			log.Fatalf("can not hash default password %v", err)
		}
		createUserArg := db.CreateUserParams{
			Username: initUsername,
			Role:     initUserRole,
			FullName: initUserFullname,
			Email:    initUserEmail,
			Password: hash,
		}
		_, err = store.CreateUser(ctx, createUserArg)
		if err != nil {
			log.Fatalf("can not create an initial admin %v", err)
		}
	}
}
