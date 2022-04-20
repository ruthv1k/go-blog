package types

type PublicPost struct {
	PostId string `json:"post_id"`
	AuthorId string `json:"author_id"`
	Title string `json:"title"`
	Description string `json:"description"`
}