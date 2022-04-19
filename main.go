package main

import (
	PostsControllers "go-auth/controllers/Posts"
	UserControllers "go-auth/controllers/Users"
	mongoconnect "go-auth/database"
	"go-auth/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// public routes
	e.POST("/accounts/register", UserControllers.RegisterUser)
	e.POST("/accounts/login", UserControllers.LoginUser)

	config := middleware.JWTConfig {
		Claims:     &types.CustomJWTClaims{},
		SigningKey: []byte("secret"),
	}

	// protected routes
	r := e.Group("/posts")
	r.Use(middleware.JWTWithConfig(config))

	// posts
	r.GET("", PostsControllers.GetUserPosts)
	r.POST("", PostsControllers.CreatePost)

	e.Logger.Fatal(e.Start(":5000"))

	defer mongoconnect.DisconnectDb()
}