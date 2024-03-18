package token

import (
	"time"

	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
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

func GetKey(token *jwt.Token) (interface{}, error) {
	if token.Method != jwt.SigningMethodHS256 {
		return nil, echo.ErrForbidden
	}
	return []byte("secret"), nil
}

type JwtCustomClaims struct {
	Username             string    `json:"username"`
	Role                 string    `json:"role"`
	TokenID              uuid.UUID `json:"tokenId"`
	jwt.RegisteredClaims `json:"claims"`
}

type JwtRefreshClaims struct {
	Username             string    `json:"username"`
	TokenID              uuid.UUID `json:"tokenId"`
	jwt.RegisteredClaims `json:"claims"`
}

type JwtEmployeeSignupClaims struct {
	Email                string    `json:"email"`
	TokenID              uuid.UUID `json:"tokenId"`
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(utils.ServerConfig.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, _ := GetKey(token)
	signedToken, err := token.SignedString(key)
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(utils.ServerConfig.RefreshTokenDuration)), //add this to a config file?
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, _ := GetKey(token)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", nil, err
	}

	return signedToken, claims, nil
}

func MakeEmployeeSignupToken(email string) (string, *JwtEmployeeSignupClaims, error) {
	newTokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}
	claims := &JwtEmployeeSignupClaims{
		email,
		newTokenID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(utils.ServerConfig.EmpSignupTokenDuration)), //add this to a config file?
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, _ := GetKey(token)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", nil, err
	}

	return signedToken, claims, nil

}
