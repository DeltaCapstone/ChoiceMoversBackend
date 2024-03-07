package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// SHOULD THIS ID BE USER ID OR TOKEN ID???????
type JwtCustomClaims struct {
	Id                   int    `json:"id"`
	UserName             string `json:"username"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims `json:"claims"`
}

// SAME????
type JwtRefreshClaims struct {
	Id                   int `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func MakeTokenPair(id int, username string, role string) (map[string]string, error) {
	t, err := MakeToken(id, username, role)
	if err != nil {
		return nil, err
	}
	rt, err := MakeRefreshToken(id)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"accessToken":  t,
		"refreshToken": rt,
	}, nil
}

func MakeToken(id int, username string, role string) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
		id,
		username,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func MakeRefreshToken(id int) (string, error) {
	// Set custom claims
	claims := &JwtRefreshClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
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
