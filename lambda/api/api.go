package api

import (
	"context"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient
}

func NewApiHandler(dbStore database.DynamoDBClient) *ApiHandler {
	return &ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(ctx context.Context, event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("Username and Password are required")
	}

	userExists, err := api.dbStore.DoesUserExist(ctx, event.Username)
	if err != nil {
		return fmt.Errorf("Error checking user existence: %w", err)
	}
	if userExists {
		return fmt.Errorf("User already exists")
	}

	err = api.dbStore.InsertUser(ctx, event)
	if err != nil {
		return fmt.Errorf("Error inserting user: %w", err)
	}

	return nil
}
