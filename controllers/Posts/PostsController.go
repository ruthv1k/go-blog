package controllers

import (
	mongoconnect "go-auth/database"
	"go-auth/types"
	"net/http"

	"github.com/golang-jwt/jwt"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreatePost(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*types.CustomJWTClaims)
	userId := claims.UserId

	post := new(types.PublicPost)

	if err := c.Bind(post); err != nil {
		return err
	}
	
	post = &types.PublicPost{
		PostId: uuid.NewString(),
		AuthorId: userId,
		Title: post.Title,
		Description: post.Description,
	}

	postsCollection, ctx, cancel := mongoconnect.GetCollection("posts")
	defer cancel()
	
	if _, err := postsCollection.InsertOne(ctx, post); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{ "message": err.Error() })
	}

	return c.JSON(http.StatusCreated, echo.Map{ "message": "Post created" })
}

func GetUserPosts(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*types.CustomJWTClaims)
	userId := claims.UserId

	postsCollection, ctx, cancel := mongoconnect.GetCollection("posts")
	defer cancel()

	postOptions := options.Find().SetProjection(bson.D{{"_id", 0}})

	cursor, err := postsCollection.Find(ctx, bson.M{ "authorid": userId }, postOptions);

	if err != nil  {
		return c.JSON(http.StatusOK, echo.Map{ "message": "No posts found" })
	}

	var posts []types.PublicPost

	if err = cursor.All(ctx, &posts); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{ "message": "Error fetching the posts" })
	}

	return c.JSON(http.StatusOK, echo.Map{ "posts": posts })
}

func UpdatePost(c echo.Context) error {
	postId := c.Param("post_id")
	postFromReq := new(types.PublicPost)

	if err := c.Bind(postFromReq); err != nil {
		return err
	}

	postFromReq = &types.PublicPost{
		Title: postFromReq.Title,
		Description: postFromReq.Description,
	}

	postsCollection, ctx, cancel := mongoconnect.GetCollection("posts")
	defer cancel()

	filter := bson.M{"postid": postId}
	update := bson.D{{"$set", bson.D{{"title", postFromReq.Title}, {"description", postFromReq.Description}}}}
	document := postsCollection.FindOneAndUpdate(ctx, filter, update)

	if document.Err() != nil  {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": "Post not found" })
	}

	var updatedPost types.PublicPost
	document.Decode(&updatedPost)

	return c.JSON(http.StatusOK, echo.Map{ "message": "Post updated", "post": updatedPost })
}

func DeletePost(c echo.Context) error {
	postId := c.Param("post_id")
	postsCollection, ctx, cancel := mongoconnect.GetCollection("posts")
	defer cancel()
	filter := bson.M{"postid": postId}
	document := postsCollection.FindOneAndDelete(ctx, filter)

	if document.Err() != nil {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": "Post not found" })
	}

	var deletedPost types.PublicPost
	document.Decode(&deletedPost)

	return c.JSON(http.StatusOK, echo.Map{ "message": "Post deleted", "postid": deletedPost.PostId })
}