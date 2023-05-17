package user

import (
	"errors"
	"flaver/api"
	"flaver/api/request"
	"flaver/api/response"
	"flaver/globals"
	"flaver/lib/dal/database/dal"
	"flaver/lib/utils"
	"flaver/models"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	userDal dal.IUserDal

	firebaseAuthUtil utils.IFirebaseAuthUtil
}

type UserServiceOption func(*UserService)

func NewUserServiceOption(options ...func(*UserService)) IUserService {
	service := UserService{}

	for _, option := range options {
		option(&service)
	}

	return &service
}

func WithUserDal(dal dal.IUserDal) UserServiceOption {
	return func(service *UserService) {
		service.userDal = dal
	}
}

func WithFirebaseAuthUtil(util utils.IFirebaseAuthUtil) UserServiceOption {
	return func(service *UserService) {
		service.firebaseAuthUtil = util
	}
}

func (this *UserService) Login(arg *request.LoginArg) (*response.Token, error) {
	var token string
	if arg.GoogleIdToken != nil {
		token = *arg.GoogleIdToken
	} else if arg.AppleIdToken != nil {
		token = *arg.AppleIdToken
	} else {
		return nil, api.InvalidArgument
	}

	if idTokenInfo, err := this.firebaseAuthUtil.ParseIdToken(token); err != nil {
		return nil, api.IdTokenParseFailed
	} else if user, err := this.userDal.GetUser(nil, "google_id_token = ? OR apple_id_token = ?", idTokenInfo.TokenUID); err != nil {
		return nil, api.NotFound
	} else if accessToken, err := utils.JwtTokenGenerator(user.Uid, string(user.Role), 1); err != nil {
		return nil, api.JWTTokenGenerateFailed
	} else if refreshToken, err := utils.JwtTokenGenerator(user.Uid, string(user.Role), 720); err != nil {
		return nil, api.JWTTokenGenerateFailed
	} else {
		return &response.Token{
			AccessToken:  accessToken.Token,
			RefreshToken: refreshToken.Token,
		}, nil
	}
}

func (this *UserService) RefreshToken(userUid, refreshToken string) (*response.Token, error) {

	if jwtTokenInfo, err := utils.ParseJWTToken(refreshToken); err != nil {
		return nil, api.JWTTokenParseFailed
	} else if claims, _ := jwtTokenInfo.Claims.(*utils.FlaverCliaims); claims.Audience != userUid || claims.ExpiresAt < time.Now().Unix() {
		return nil, api.PermissionDenied
	} else if user, err := this.userDal.GetUserProfileByUid(userUid); err != nil {
		return nil, api.NotFound
	} else if accessToken, err := utils.JwtTokenGenerator(user.Uid, string(user.Role), 1); err != nil {
		return nil, api.JWTTokenGenerateFailed
	} else if refreshToken, err := utils.JwtTokenGenerator(user.Uid, string(user.Role), 720); err != nil {
		return nil, api.JWTTokenGenerateFailed
	} else {
		return &response.Token{
			AccessToken:  accessToken.Token,
			RefreshToken: refreshToken.Token,
		}, nil
	}
}

func (this *UserService) RegistUser(data *request.RegisteArg) (*response.Token, error) {

	var token string
	if data.GoogleIdToken != nil {
		token = *data.GoogleIdToken
	} else if data.AppleIdToken != nil {
		token = *data.AppleIdToken
	} else {
		return nil, api.InvalidArgument
	}

	if idTokenInfo, err := this.firebaseAuthUtil.ParseIdToken(token); err != nil {
		return nil, api.IdTokenParseFailed
	} else {

		user := models.User{
			Nickname:     data.Nickname,
			AvatarImgUrl: data.AvatarImgUrl,
			Email:        idTokenInfo.Email,
			Role:         models.GeneralRole,
		}
		if data.GoogleIdToken != nil {
			user.GoogleIdToken = &idTokenInfo.TokenUID
		} else if data.AppleIdToken != nil {
			user.AppleIdToken = &idTokenInfo.TokenUID
		}
		if err := this.userDal.CreateUser(&user); err != nil {
			return nil, err
		} else if err = this.replaceUserTopics(user.Uid, data.TopicIds); err != nil {
			return nil, err
		}

		if accessToken, err := utils.JwtTokenGenerator(user.Uid, string(user.Role), 1); err != nil {
			return nil, api.JWTTokenGenerateFailed
		} else if refreshToken, err := utils.JwtTokenGenerator(user.Uid, string(user.Role), 720); err != nil {
			return nil, api.JWTTokenGenerateFailed
		} else {
			return &response.Token{
				AccessToken:  accessToken.Token,
				RefreshToken: refreshToken.Token,
			}, nil
		}
	}
}

func (this *UserService) GetUser(uid string) (*models.User, error) {
	if user, err := this.userDal.GetUserProfileByUid(uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.RecordNotFound
		} else {
			return nil, err
		}
	} else {
		return user, nil
	}
}

func (this *UserService) Update(userUid string, data *request.UpdateUserArg) error {

	if _, err := this.userDal.GetUserProfileByUid(userUid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return api.RecordNotFound
		} else {
			return err
		}
	}

	updates := map[string]interface{}{}

	if data.Nickname != nil {
		updates["nickname"] = *data.Nickname
	}

	if data.Birthday != nil {
		updates["birthday"] = *data.Birthday
	}

	if data.AvatarImgUrl != nil {
		updates["avatar_img_url"] = *data.AvatarImgUrl
	}

	if len(data.TopicIds) > 0 {
		if err := this.replaceUserTopics(userUid, data.TopicIds); err != nil {
			return err
		}
	}

	return this.userDal.UpdateUser(updates, "uid = ?", userUid)
}

func (this *UserService) replaceUserTopics(userUid string, topicIds []int) error {

	if len(topicIds) > globals.GetViper().GetInt("topic.max_count") {
		return api.TopicCountIllegal
	} else if err := this.userDal.DeleteUserTopics(userUid); err != nil {
		return err
	} else {
		userTopics := []*models.UserTopic{}
		for _, topicId := range topicIds {
			userTopics = append(userTopics, &models.UserTopic{UserUid: userUid, TopicId: topicId})
		}
		if err = this.userDal.CreateUserTopics(userTopics); err != nil {
			return err
		}
	}

	return nil
}
