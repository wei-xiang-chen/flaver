package post

import (
	"flaver/api"
	"flaver/api/request"
	"flaver/models"
)

type IPostService interface {
	GetPosts(userUid string, search *string, nextId *int) (*api.HasNextResult[[]*models.Post], error)
	GetPost(userUid string, id int) (*models.Post, error)
	Create(data request.CreatePostArg) (*models.Post, error)
	Update(id int, userUid string, data request.UpdatePostArg) error
	
	Like(userUid string, postId int, postLikeType request.PostLikeType) error
	GetPostsLikeCount(postIds []int) (map[int]int, error)
}
