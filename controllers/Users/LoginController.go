package controllers

import (
	"context"
	"go-auth/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(c echo.Context) (err error) {
	u := new(types.PublicUser)

	if err = c.Bind(u); err != nil {
		return err
	}

	defer cancel()

	userDetails, _ := GetUserByEmail(*usersCollection, ctx, u.Email)

	if userDetails.Email != u.Email {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": "User not found" })
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{ "message": "Invalid email / password" })
	}
	
	claims := &types.CustomJWTClaims {
		UserId:      userDetails.UserId.Hex(),
		DisplayName:    userDetails.DisplayName,
		Email:   userDetails.Email,
		Role: userDetails.Role,
		ExpiresAt: time.Now().Add(time.Hour * 1).UnixNano() / int64(time.Millisecond),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func GetUserByEmail(collection mongo.Collection, ctx context.Context, email string) (types.User, error) {
	var userdetails types.User
	
	err := collection.FindOne(ctx, bson.M{ "email": email }).Decode(&userdetails)

	return userdetails, err
}
