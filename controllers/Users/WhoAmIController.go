package controllers

import (
	"go-auth/types"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func WhoAmI(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*types.CustomJWTClaims)
	userId := claims.UserId
	email := claims.Email
	displayname := claims.DisplayName

	return c.JSON(http.StatusCreated, echo.Map{ "userid": userId, "email": email, "displayname": displayname })
}