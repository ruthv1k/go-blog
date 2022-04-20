package controllers

import (
	mongoconnect "go-auth/database"
	"go-auth/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection, ctx, cancel = mongoconnect.GetCollection("users")

func RegisterUser(c echo.Context) (err error) {
	u := new(types.PublicUser)

	if err = c.Bind(u); err != nil {
		return err
	}

	if userDetails, _ := GetUserByEmail(*usersCollection, ctx, u.Email); userDetails.Email == u.Email {
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

	if _, err := usersCollection.InsertOne(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{ "message": err.Error() })
	}

	return c.JSON(http.StatusOK, echo.Map{ "message": "Account created" })
}