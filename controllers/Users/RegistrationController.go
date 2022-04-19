package controllers

import (
	"context"
	mongoconnect "go-auth/database"
	"go-auth/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c echo.Context) (err error) {
	u := new(types.PublicUser)

	if err = c.Bind(u); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := mongoconnect.GetDatabase()
	collection := db.Collection("users")

	if userDetails, _ := GetUserByEmail(*collection, ctx, u.Email); userDetails.Email == u.Email {
		return c.JSON(http.StatusBadRequest, echo.Map{ "message":  "User already exists" })
	}
	
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 13)
	if err != nil {
		return err
	}

	user :=  types.User {
		DisplayName: u.DisplayName,
		UserName: u.UserName,
		Email: u.Email,
		Password: string(password),
		Role: "writer",
	}

	if _, err := collection.InsertOne(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{ "message": err.Error() })
	}

	return c.JSON(http.StatusOK, echo.Map{ "message": "Account created" })
}