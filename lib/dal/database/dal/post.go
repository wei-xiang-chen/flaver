package dal

import (
	"flaver/api"
	"flaver/globals"
	"flaver/lib/constants"
	"flaver/lib/utils"
	"flaver/models"
	"strconv"
)

type IPostDal interface {
	GetPostPaged(userUid string, search *string, nextId *int) (*api.HasNextResult[[]*models.Post], error)
	GetPost(userUid string, preloads []string, where string, args ...interface{}) (*models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(updates map[string]interface{}, where string, args ...interface{}) error

	CreatePostLike(postLike *models.PostLike) error
	DeletePostLike(userUid string, postId int) error
	GetPostsLikeCount(postIds []int) (map[int]int, error)
	GetPostLikeCount(userUid string, postId int) (int, error)
}

func (this *Dal) GetPostPaged(userUid string, search *string, nextId *int) (*api.HasNextResult[[]*models.Post], error) {

	posts := make([]*models.Post, 0)

	query := this.db.Preload("Owner").Preload("Restaurant").Preload("PostLikes", "user_uid = ?", userUid)
	if search != nil {
		query.Where("title like ? OR  content like ? ", "%"+*search+"%", "%"+*search+"%")
	}

	if nextId != nil {
		query.Where("id >= ?", *nextId)
	}

	if err := query.Order("id DESC").Limit(api.DEFAULT_PAGE_SIZE + 1).Find(&posts).Error; err != nil {
		return nil, err
	}

	nextId = nil
	if len(posts) > api.DEFAULT_PAGE_SIZE {
		nextId = &posts[api.DEFAULT_PAGE_SIZE].Id
	}

	dataLength := len(posts)
	if nextId != nil {
		dataLength = len(posts) - 1
	}
	return &api.HasNextResult[[]*models.Post]{
		Data:   posts[:dataLength],
		NextId: nextId,
	}, nil
}

func (this *Dal) GetPost(userUid string, preloads []string, where string, args ...interface{}) (*models.Post, error) {

	post := models.Post{}

	query := this.db.Model(&models.Post{}).Preload("PostLikes", "user_uid = ?", userUid)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Where(where, args...).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (this *Dal) CreatePost(post *models.Post) error {
	return this.db.Model(&models.Post{}).Create(post).Error
}

func (this *Dal) UpdatePost(updates map[string]interface{}, where string, args ...interface{}) error {
	return this.db.Model(&models.Post{}).Where(where, args...).Updates(updates).Error
}

func (this *Dal) CreatePostLike(postLike *models.PostLike) error {
	if err := this.db.Create(postLike).Error; err != nil {
		return err
	} else if err := this.redis.HIncrBy(constants.POST_LIKE_COUNT_KEY, strconv.Itoa(postLike.PostId), -1).Err(); err != nil {
		globals.GetLogger().Warnf("[CreatePostLike] Cache error: %v", err)
	}

	return nil
}

func (this *Dal) DeletePostLike(userUid string, postId int) error {

	if err := this.db.Where("user_uid = ? AND post_id = ?", userUid, postId).Delete(&models.PostLike{}).Error; err != nil {
		return err
	} else if err := this.redis.HIncrBy(constants.POST_LIKE_COUNT_KEY, strconv.Itoa(postId), -1).Err(); err != nil {
		globals.GetLogger().Warnf("[DeletePostLike] Cache error: %v", err)
	}

	return nil
}

func (this *Dal) GetPostsLikeCount(postIds []int) (map[int]int, error) {

	if resultMap, notFoundIds, err := this.getPostsLikeCountFromCache(postIds); err != nil {
		if result, err := this.getPostsLikeCountFromDB(postIds); err != nil {
			return nil, err
		} else {
			this.setPostsLikeCountToCache(result)

			return result, nil
		}
	} else if len(notFoundIds) > 0 {
		if result, err := this.getPostsLikeCountFromDB(notFoundIds); err != nil {
			return nil, err
		} else {
			this.setPostsLikeCountToCache(result)

			if len(result) > len(resultMap) {
				for k, v := range resultMap {
					result[k] = v
				}
				return result, nil
			} else {
				for k, v := range result {
					resultMap[k] = v
				}
				return resultMap, nil
			}
		}
	} else {
		return resultMap, nil
	}

}

// daily del this key
func (this *Dal) setPostsLikeCountToCache(data map[int]int) {
	fields := make(map[string]interface{})
	for k, v := range data {
		fields[strconv.Itoa(k)] = v
	}
	if err := this.redis.HMSet(constants.POST_LIKE_COUNT_KEY, fields).Err(); err != nil {
		globals.GetLogger().Warnf("[setPostsLikeCountToCache] Cache error: %v", err)
	}
}

func (this *Dal) getPostsLikeCountFromDB(postIds []int) (map[int]int, error) {
	var results []map[string]interface{}

	if err := this.db.Model(&models.PostLike{}).Select("post_id, COUNT(*) AS count").Where("post_id IN ?", postIds).Group("post_id").Find(&results).Error; err != nil {
		return nil, err
	}

	resultMap := map[int]int{}
	for _, result := range results {
		resultMap[result["post_id"].(int)] = int(result["count"].(int64))
	}

	return resultMap, nil
}

func (this *Dal) getPostsLikeCountFromCache(postIds []int) (map[int]int, []int, error) {

	result := make(map[int]int)
	notFoundIds := make([]int, 0)

	if redisResult, err := this.redis.HMGet(constants.POST_LIKE_COUNT_KEY, utils.IntArrToStrArr(postIds)...).Result(); err != nil {
		return nil, nil, err
	} else {
		for i, val := range redisResult {
			if val == nil {
				notFoundIds = append(notFoundIds, postIds[i])
			} else {
				result[postIds[i]] = val.(int)
			}
		}
	}

	return result, notFoundIds, nil
}

func (this *Dal) GetPostLikeCount(userUid string, postId int) (int, error) {

	var count int64

	query := this.db.Model(&models.PostLike{})

	if err := query.Where("user_uid = ? AND post_id = ?", userUid, postId).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
