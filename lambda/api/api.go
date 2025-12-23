package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) *ApiHandler {
	return &ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var event types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, fmt.Errorf("Invalid request body: %v", err)
	}
	if event.Username == "" || event.Password == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Username and password are required",
		}, nil
	}

	userExists, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error checking user existence: %v", err),
		}, fmt.Errorf("Error checking user existence: %v", err)
	}
	if userExists {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusConflict,
			Body:       "User already exists",
		}, nil
	}

	user, err := types.NewUser(event.Username, event.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error creating user: %v", err),
		}, fmt.Errorf("Error creating user: %v", err)
	}

	err = api.dbStore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error inserting user: %v", err),
		}, fmt.Errorf("Error inserting user: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User registered successfully",
	}, nil
}

func (api ApiHandler) HealthCheckHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "API is healthy",
	}, nil
}

func (api ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	type LoginUser struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var event LoginUser
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid request body",
		}, fmt.Errorf("Invalid request body: %v", err)
	}
	if event.Username == "" || event.Password == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Username and password are required",
		}, nil
	}

	user, err := api.dbStore.GetUser(event.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error checking user existence: %v", err),
		}, fmt.Errorf("Error checking user existence: %v", err)
	}

	if user.Username != event.Username {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Invalid username or password",
		}, nil
	}

	isPasswordValid, err := types.ValidatePassword(user.PasswordHash, event.Password)

	if err != nil || !isPasswordValid {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Invalid username or password",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Login successful",
	}, nil
}
