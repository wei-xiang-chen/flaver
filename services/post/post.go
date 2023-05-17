package post

import (
	"errors"
	"flaver/api"
	"flaver/api/request"
	"flaver/lib/dal/database/dal"
	"flaver/lib/utils"
	"flaver/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PostService struct {
	postDal       dal.IPostDal
	userDal       dal.IUserDal
	restaurantDal dal.IRestaurantDal

	mapUtil utils.IMapUtil
}

type PostServiceOption func(*PostService)

func NewPostServiceOption(options ...func(*PostService)) IPostService {
	service := PostService{}

	for _, option := range options {
		option(&service)
	}

	return &service
}

func WithPostDal(dal dal.IPostDal) PostServiceOption {
	return func(service *PostService) {
		service.postDal = dal
	}
}

func WithUserDal(dal dal.IUserDal) PostServiceOption {
	return func(service *PostService) {
		service.userDal = dal
	}
}

func WithRestaurantDal(dal dal.IRestaurantDal) PostServiceOption {
	return func(service *PostService) {
		service.restaurantDal = dal
	}
}

func WithMapUtil(util utils.IMapUtil) PostServiceOption {
	return func(service *PostService) {
		service.mapUtil = util
	}
}

func (this *PostService) GetPosts(userUid string, search *string, nextId *int) (*api.HasNextResult[[]*models.Post], error) {
	return this.postDal.GetPostPaged(userUid, search, nextId)
}

func (this *PostService) GetPost(userUid string, id int) (*models.Post, error) {
	return this.postDal.GetPost(userUid, []string{"Owner", "Restaurant"}, "id = ?", id)
}

func (this *PostService) Create(data request.CreatePostArg) (*models.Post, error) {

	if user, err := this.userDal.GetUserProfileByUid(data.UserUid); err != nil {
		return nil, err
	} else {
		if isRestaurantExist, err := this.restaurantDal.CheckRestaurantIsExist(data.GooglePlaceId); err != nil {
			return nil, err
		} else if !isRestaurantExist {
			if googlePlaceDetail, err := this.mapUtil.GetPlaceDetails(data.GooglePlaceId); err != nil {
				return nil, err
			} else if err = this.restaurantDal.CreateRestaurant(&models.Restaurant{
				GooglePlaceId: googlePlaceDetail.PlaceID,
				Name:          googlePlaceDetail.Name,
				Phone:         &googlePlaceDetail.FormattedPhoneNumber,
				Address:       googlePlaceDetail.FormattedAddress,
				// Location: googlePlaceDetail.Geometry.Location,
				Lat: googlePlaceDetail.Geometry.Location.Lat,
				Lng: googlePlaceDetail.Geometry.Location.Lng,
			}); err != nil {
				return nil, err
			}
		}

		var post = &models.Post{
			OwnerUid:      data.UserUid,
			Title:         data.Title,
			Content:       data.Content,
			ImgUrls:       data.ImgUrls,
			Rating:        data.Rating,
			TopicIds:      data.TopicIds,
			GooglePlaceId: data.GooglePlaceId,
			InstantlyRole: string(user.Role),
		}
		if err := this.postDal.CreatePost(post); err != nil {
			return nil, err
		}

		return post, nil
	}
}

func (this *PostService) Update(id int, userUid string, data request.UpdatePostArg) error {

	if post, err := this.postDal.GetPost(userUid, nil, "id = ?", id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.RecordNotFound
		} else {
			return err
		}
	} else if post.OwnerUid != userUid {
		return api.PermissionDenied
	}

	updates := map[string]interface{}{}

	if data.Title != nil {
		updates["title"] = *data.Title
	}
	if data.Content != nil {
		updates["content"] = *data.Content
	}
	if data.ImgUrls != nil {
		updates["img_urls"] = pq.Array(*data.ImgUrls)
	}
	if data.Rating != nil {
		updates["rating"] = *data.Rating
	}
	if data.TopicIds != nil {
		updates["topic_ids"] = pq.Array(*data.TopicIds)
	}
	// if data.GooglePlaceId != nil {
	// 	updates["google_place_id"] = *data.GooglePlaceId
	// }

	return this.postDal.UpdatePost(updates, "id = ?", id)
}

func (this *PostService) Like(userUid string, postId int, postLikeType request.PostLikeType) error {

	count, err := this.postDal.GetPostLikeCount(userUid, postId)
	if err != nil {
		return err
	} else if count > 0 && postLikeType == request.PostLikeTypeLike {
		return nil
	} else if count == 0 && postLikeType == request.PostLikeTypeRevoke {
		return nil
	}

	if postLikeType == request.PostLikeTypeLike {
		return this.postDal.CreatePostLike(&models.PostLike{
			UserUid: userUid,
			PostId:  postId,
		})
	} else if postLikeType == request.PostLikeTypeRevoke {
		return this.postDal.DeletePostLike(userUid, postId)
	} else {
		return nil
	}
}

func (this *PostService) GetPostsLikeCount(postIds []int) (map[int]int, error) {
	return this.postDal.GetPostsLikeCount(postIds)
}
