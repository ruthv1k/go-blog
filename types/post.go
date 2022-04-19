package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type PublicPost struct {
	PostId primitive.ObjectID `bson:"_id,omitempty" json:"post_id"`
	AuthorId string `json:"author_id"`
	Title string `json:"title"`
	Description string `json:"description"`
}