package main

import (
	"fmt"
	"net/http"
	"time"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/DeltaCapstone/ChoiceMoversBackend/token"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

func VerifyRefreshToken(r string) (*token.JwtRefreshClaims, error) {
	rt, err := jwt.ParseWithClaims(r, &token.JwtRefreshClaims{}, token.GetKey)
	if err != nil {
		zap.L().Sugar().Errorf("Could not verify refresh token: ", err.Error())
		switch err {
		case jwt.ErrTokenExpired:
			// Handle validation errors
			return nil, echojwt.ErrJWTInvalid
		default:
			// Handle other errors
			return nil, echojwt.ErrJWTMissing
		}
	}
	return rt.Claims.(*token.JwtRefreshClaims), nil
}

func renewAccessToken(c echo.Context) error {
	var req renewAccessTokenRequest
	if err := c.Bind(&req); err != nil {
		zap.L().Sugar().Errorf("Could not bind refresh token: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	refreshClaims, err := VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		zap.L().Sugar().Errorf("Could not verify refresh token: ", err.Error())
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}

	session, err := DB.PgInstance.GetSession(c.Request().Context(), refreshClaims.TokenID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if session.ID != refreshClaims.TokenID {
		return c.String(http.StatusNotFound, fmt.Sprintf("Session not found: %v", refreshClaims.TokenID))
	}

	if session.IsBlocked {
		return c.JSON(http.StatusUnauthorized, "Blocked Session")
	}

	if session.Username != refreshClaims.Username {
		return c.JSON(http.StatusUnauthorized, "Incorrect Session user")
	}

	if session.RefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, "session token missmatch")
	}

	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, "Session Expired")
	}

	accessToken, accessClaims, err := token.MakeAccessToken(session.Username, session.Role)
	if err != nil {
		zap.L().Sugar().Errorf("problem making new access token: ", err.Error())
		return c.String(http.StatusInternalServerError, "Error creating new access token")
	}
	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,
	}
	return c.JSON(http.StatusOK, rsp)
}
