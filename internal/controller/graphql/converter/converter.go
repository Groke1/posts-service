package converter

import (
	"posts-service/internal/entity"
	"posts-service/pkg/generated/posts/graphql/model"
)

func ToGraphQLComment(comment *entity.Comment) *model.Comment {
	return &model.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		AuthorID:  comment.AuthorID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func ToGraphQLComments(comments []*entity.Comment) []*model.Comment {
	result := make([]*model.Comment, len(comments))
	for i, comment := range comments {
		result[i] = ToGraphQLComment(comment)
	}
	return result
}

func ToGraphQLPost(post *entity.Post) *model.Post {
	return &model.Post{
		ID:              post.ID,
		Title:           post.Title,
		Text:            post.Text,
		AuthorID:        post.AuthorID,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	}
}

func ToGraphQLPosts(posts []*entity.Post) []*model.Post {
	result := make([]*model.Post, len(posts))
	for i, post := range posts {
		result[i] = ToGraphQLPost(post)
	}
	return result
}
