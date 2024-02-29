package token

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// TODO: what to do here?
func getKey() string {
	var key string
	return key
}

var config = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(jwtCustomClaims)
	},
	SigningMethod: jwt.SigningMethodHS256.Name,
	SigningKey:    []byte("secret"),
}

type jwtCustomClaims struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates the JWT token and sets the user role in the context.
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := extractJWTToken(c.Request())
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		claims := &jwtCustomClaims{}
		// Validate and parse the token into claims
		// You should use a JWT library (e.g., github.com/dgrijalva/jwt-go) for this
		if err := parseAndValidateToken(token, claims); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		// Set the user role in the context
		c.Set("role", claims.Role)

		return next(c)
	}
}

// extractJWTToken extracts the JWT token from the request header.
func extractJWTToken(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// parseAndValidateToken validates and parses the JWT token.
// You should use a JWT library for this (e.g., github.com/dgrijalva/jwt-go).
func parseAndValidateToken(tokenString string, claims *jwtCustomClaims) error {
	// Implement your JWT validation logic here
	// Use a JWT library to validate and parse the token
	// Example:
	// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
	//     return []byte("your-secret-key"), nil
	// })
	// if err != nil || !token.Valid {
	//     return errors.New("invalid token")
	// }
	return nil
}
