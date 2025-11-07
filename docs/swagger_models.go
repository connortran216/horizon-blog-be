package docs

import "go-crud/schemas"

// swagger:route POST /posts posts CreatePost
// Creates a new post with title and content.
// responses:
//   201: PostResponse
//   400: ErrorResponse
//   500: ErrorResponse

// swagger:route GET /posts posts ListPosts
// Get a paginated list of posts.
// responses:
//   200: ListPostsResponse
//   500: ErrorResponse

// swagger:route GET /posts/{id} posts GetPost
// Get a single post by ID.
// responses:
//   200: PostResponse
//   400: ErrorResponse
//   404: ErrorResponse

// swagger:route PUT /posts/{id} posts UpdatePost
// Update a complete post by ID.
// responses:
//   200: PostResponse
//   400: ErrorResponse
//   404: ErrorResponse

// swagger:route PATCH /posts/{id} posts PartialUpdatePost
// Partially update a post by ID.
// responses:
//   200: PostResponse
//   400: ErrorResponse
//   404: ErrorResponse

// swagger:route DELETE /posts/{id} posts DeletePost
// Delete a post by ID.
// responses:
//   200: MessageResponse
//   400: ErrorResponse
//   404: ErrorResponse

// swagger:parameters CreatePost UpdatePost
type PostRequestBody struct {
	// in:body
	Body schemas.CreatePostRequest
}

// swagger:parameters PartialUpdatePost
type PatchPostRequestBody struct {
	// in:body
	Body schemas.PatchPostRequest
}

// swagger:parameters GetPost UpdatePost PartialUpdatePost DeletePost
type PostIDParam struct {
	// Post ID
	// in:path
	// required: true
	ID uint `json:"id"`
}

// swagger:parameters ListPosts
type ListPostsParams struct {
	// Page number
	// in:query
	// default: 1
	Page int `json:"page"`
	// Items per page
	// in:query
	// default: 10
	Limit int `json:"limit"`
}
