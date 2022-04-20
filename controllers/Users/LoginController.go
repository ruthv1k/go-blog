package controllers

import (
	mongoconnect "go-auth/database"
	"go-auth/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(c echo.Context) (err error) {
	u := new(types.PublicUser)

	if err = c.Bind(u); err != nil {
		return err
	}

	usersCollection, ctx, cancel := mongoconnect.GetCollection("users")
	defer cancel()

	userDetails := usersCollection.FindOne(ctx, bson.M{ "email": u.Email })
	
	if userDetails.Err() != nil  {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": "User not found" })
	}

	var user types.User
	userDetails.Decode(&user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{ "message": "Invalid email / password" })
	}
	
	claims := &types.CustomJWTClaims {
		UserId:      user.UserId,
		DisplayName:    user.DisplayName,
		Email:   user.Email,
		Role: user.Role,
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
