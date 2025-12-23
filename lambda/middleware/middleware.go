package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

// extract the headers
// extract the claims
// verify the token

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Extract the Authorization header
		authHeader := request.Headers["Authorization"]
		if authHeader == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Missing Authorization header",
			}, nil
		}

		// Extract the token from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Invalid Authorization header format",
			}, nil
		}

		// Parse and validate the token

		secretKey := []byte("your-secret-key")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Invalid token",
			}, fmt.Errorf("Invalid token: %v", err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Invalid token claims",
			}, fmt.Errorf("Invalid token claims")
		}

		// Check token expiration
		if exp, ok := claims["expires"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusUnauthorized,
					Body:       "Token has expired",
				}, nil
			}
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body:       "Invalid token expiration",
			}, fmt.Errorf("Invalid token expiration")
		}

		// Token is valid, proceed to the next handler
		return next(request)
	}
}
