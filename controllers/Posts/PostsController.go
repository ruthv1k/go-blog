package controllers

import (
	"context"
	mongoconnect "go-auth/database"
	"go-auth/types"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePost(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*types.CustomJWTClaims)
	userId := claims.UserId

	post := new(types.PublicPost)

	if err := c.Bind(post); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := mongoconnect.GetDatabase()
	postsCollection := db.Collection("posts")

	post = &types.PublicPost{
		AuthorId: userId,
		Title: post.Title,
		Description: post.Description,
	}
	
	if _, err := postsCollection.InsertOne(ctx, post); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{ "message": err.Error() })
	}

	return c.JSON(http.StatusOK, echo.Map{ "message": "Post created" })
}

func GetUserPosts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*types.CustomJWTClaims)
	userId := claims.UserId

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := mongoconnect.GetDatabase()
	postsCollection := db.Collection("posts")

	documents := postsCollection.FindOne(ctx, bson.M{ "authorid": userId });

	if documents.Err() == mongo.ErrNoDocuments || documents.Err() != nil  {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": "Posts not found" })
	}

	var posts bson.M = make(bson.M)
	documents.Decode(&posts)

	return c.JSON(http.StatusOK, echo.Map{ "posts": posts })
}

func UpdatePost(c echo.Context) error {
	postId := c.Param("post_id")
	post := new(types.PublicPost)

	if err := c.Bind(post); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := mongoconnect.GetDatabase()
	postsCollection := db.Collection("posts")

	post = &types.PublicPost{
		Title: post.Title,
		Description: post.Description,
	}

	updatedPost := postsCollection.FindOneAndUpdate(ctx, bson.M{ "_id": postId }, post)	
	
	if updatedPost.Err() == mongo.ErrNoDocuments || updatedPost.Err() != nil  {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": "Post not found" })
	}

	return c.JSON(http.StatusOK, echo.Map{ "post": updatedPost })
}