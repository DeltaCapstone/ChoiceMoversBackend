package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var Config = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(JwtCustomClaims)
	},
	SigningMethod: jwt.SigningMethodHS256.Name,
	SigningKey:    []byte("secret"),
}

type JwtCustomClaims struct {
	UserName             string `json:"username"`
	Role                 string `json:"role"`
	TokenID              uuid.UUID
	jwt.RegisteredClaims `json:"claims"`
}

type JwtRefreshClaims struct {
	UserName             string `json:"username"`
	TokenID              uuid.UUID
	jwt.RegisteredClaims `json:"claims"`
}

func MakeAccessToken(username string, role string) (string, *JwtCustomClaims, error) {
	// Set custom claims
	newTokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}
	claims := &JwtCustomClaims{
		username,
		role,
		newTokenID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", nil, err
	}

	return signedToken, claims, nil
}

func MakeRefreshToken(username string) (string, *JwtRefreshClaims, error) {
	// Set custom claims
	newTokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}
	claims := &JwtRefreshClaims{
		username,
		newTokenID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", nil, err
	}

	return signedToken, claims, nil
}

///////////////////////////////////////////////////////////////
//Pretty sure this just does what echojwt does so not neccessary

/*
// JWTMiddleware validates the JWT token and sets the user role in the context.

	func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := extractJWTToken(c.Request())
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			claims := &JwtCustomClaims{}
			// Validate and parse the token into claims

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

	func parseAndValidateToken(tokenString string, claims *JwtCustomClaims) error {
		// Implement your JWT validation logic here
		// Use a JWT library to validate and parse the token
		// Example:
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil || !token.Valid || token.Method != jwt.SigningMethodHS256 {
			return errors.New("invalid token")
		}

		return nil
	}
*/
