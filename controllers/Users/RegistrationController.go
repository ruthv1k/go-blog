package controllers

import (
	mongoconnect "go-auth/database"
	"go-auth/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c echo.Context) (err error) {
	u := new(types.PublicUser)

	if err = c.Bind(u); err != nil {
		return err
	}

	usersCollection, ctx, cancel := mongoconnect.GetCollection("users")
	defer cancel()

	userDetails := usersCollection.FindOne(ctx, bson.M{ "email": u.Email })

	var user types.User
	userDetails.Decode(&user)

	if user.Email == u.Email {
		return c.JSON(http.StatusBadRequest, echo.Map{ "message":  "User already exists" })
	}
	
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 13)
	if err != nil {
		return err
	}

	user = types.User {
		UserId: uuid.NewString(),
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