package controllers

import (
	"flaver/api"
	"flaver/api/request"
	"flaver/api/response"
	"flaver/globals"
	"flaver/lib/dal"
	"flaver/lib/utils"
	"flaver/services/post"
	"flaver/services/restaurant"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	postService       post.IPostService
	restaurantService restaurant.IRestaurantService

	transactionContext dal.TransactionContext
}

func NewPostController() PostController {
	dal := dal.NewDal()

	return PostController{
		postService: post.NewPostServiceOption(
			post.WithPostDal(dal.Dal),
			post.WithUserDal(dal.Dal),
			post.WithRestaurantDal(dal.Dal),

			post.WithMapUtil(&utils.GoogleMapUtil{}),
		),
		restaurantService: restaurant.NewRestaurantServiceOption(
			restaurant.WithRestaurantDal(dal.Dal),
		),
		transactionContext: dal.ExecTransaction,
	}
}

func (this PostController) GetList(c *gin.Context) {

	userUid := c.MustGet("userUid").(string)
	arg := request.QueryPostInfoArg{}

	if err := c.ShouldBindQuery(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	pageData, err := this.postService.GetPosts(userUid, arg.Search, arg.NextId)
	if err != nil {
		globals.GetLogger().Warnf("[GetPostInfo] error: %v", err)
		api.SendResult(err, nil, c)
		return
	}

	postIds := []int{}
	resultData := make([]*response.PostList, 0)

	for _, data := range pageData.Data {
		postIds = append(postIds, data.Id)
	}

	postsLikeCount, err := this.postService.GetPostsLikeCount(postIds)
	if err != nil {
		globals.GetLogger().Warnf("[GetPostInfo] error: %v", err)
		api.SendResult(err, nil, c)
		return
	}

	for _, data := range pageData.Data {
		post := &response.PostList{}
		data.SerializeTo(post, postsLikeCount)

		resultData = append(resultData, post)
	}

	api.SendResult(nil, api.HasNextResult[[]*response.PostList]{
		Data:   resultData,
		NextId: pageData.NextId,
	}, c)
}

func (this PostController) GetDetail(c *gin.Context) {

	userUid := c.MustGet("userUid").(string)
	idStr := c.Param("id")
	i, _ := strconv.Atoi(idStr)

	post, err := this.postService.GetPost(userUid, i)
	if err != nil {
		globals.GetLogger().Warnf("[GetPostDetail] error: %v", err)
		api.SendResult(err, nil, c)
		return
	}

	postsLikeCount, err := this.postService.GetPostsLikeCount([]int{post.Id})
	if err != nil {
		globals.GetLogger().Warnf("[GetPostDetail] error: %v", err)
		api.SendResult(err, nil, c)
		return
	}

	restaurantRatingMap, err := this.restaurantService.GetRestaurantRating([]string{post.GooglePlaceId})

	result := &response.PostDetail{}
	post.SerializeTo(result, postsLikeCount)

	if val, ok := restaurantRatingMap[post.GooglePlaceId]; ok {
		rating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", val.FlaverRating), 64)
		result.Restaurant.FlaverRating = rating
		rating, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", val.GeneralRating), 64)
		result.Restaurant.GeneralRating = rating
	}

	api.SendResult(nil, result, c)
}

func (this PostController) Create(c *gin.Context) {

	userUid := c.MustGet("userUid").(string)
	arg := request.CreatePostArg{UserUid: userUid}

	var err error
	if err = c.ShouldBindJSON(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	transaction := func(tx *gorm.DB) (interface{}, error) {
		_, err := this.postService.Create(arg)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
	if _, err := this.transactionContext(transaction); err != nil {
		globals.GetLogger().Warnf("[CreatePost] transaction error: %v", err)
		api.SendResult(err, nil, c)
		return
	} else {
		api.SendResult(nil, nil, c)
	}

	api.SendResult(nil, nil, c)
}

func (this PostController) Update(c *gin.Context) {
	userUid := c.MustGet("userUid").(string)
	idStr := c.Param("id")
	postId, _ := strconv.Atoi(idStr)

	arg := request.UpdatePostArg{}
	if err := c.ShouldBindJSON(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	err := this.postService.Update(postId, userUid, arg)
	if err != nil {
		api.SendResult(err, nil, c)
		return
	}

	api.SendResult(nil, nil, c)
}

func (this PostController) Like(c *gin.Context) {
	userUid := c.MustGet("userUid").(string)
	idStr := c.Param("id")
	postId, _ := strconv.Atoi(idStr)

	var err error
	arg := request.PostLikeArg{}
	if err = c.ShouldBindJSON(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	transaction := func(tx *gorm.DB) (interface{}, error) {
		err = this.postService.Like(userUid, postId, arg.Type)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
	if _, err = this.transactionContext(transaction); err != nil {
		globals.GetLogger().Warnf("[LikePost] transaction error: %v", err)
		api.SendResult(err, nil, c)
		return
	}

	api.SendResult(nil, nil, c)
}
