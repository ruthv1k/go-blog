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
	postsRoutes := e.Group("/posts")
	usersRoutes := e.Group("/users")

	postsRoutes.Use(middleware.JWTWithConfig(config))
	usersRoutes.Use(middleware.JWTWithConfig(config))

	// posts routes
	postsRoutes.GET("", PostsControllers.GetUserPosts)
	postsRoutes.POST("", PostsControllers.CreatePost)
	postsRoutes.PUT("/:post_id", PostsControllers.UpdatePost)
	postsRoutes.DELETE("/:post_id", PostsControllers.DeletePost)

	// user routes
	usersRoutes.GET("/whoami", UserControllers.WhoAmI)

	e.Logger.Fatal(e.Start(":5000"))

	defer mongoconnect.DisconnectDb()
}