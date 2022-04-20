package controllers

import (
	mongoconnect "go-auth/database"
	"go-auth/types"
	"net/http"

	"github.com/golang-jwt/jwt"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	postOptions := options.FindOne().SetProjection(bson.D{{"_id", 0}})

	documents := postsCollection.FindOne(ctx, bson.M{ "authorid": userId }, postOptions);

	if documents.Err() == mongo.ErrNoDocuments || documents.Err() != nil  {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": "Posts not found" })
	}

	var posts bson.M = make(bson.M)
	documents.Decode(&posts)

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
	opts := options.FindOneAndUpdate().SetUpsert(true)	
	updatedPost := postsCollection.FindOneAndUpdate(ctx, filter, update, opts)
	
	if updatedPost.Err() != nil  {
		return c.JSON(http.StatusNotFound, echo.Map{ "message": updatedPost.Err() })
	}

	return c.JSON(http.StatusOK, echo.Map{ "post": updatedPost })
}