package dal

import (
	"encoding/json"
	"flaver/globals"
	"flaver/lib/constants"
	"flaver/models"
	"time"
)

type IUserDal interface {
	GetUser(preloads []string, where string, args ...interface{}) (*models.User, error)
	GetUserProfileByUid(uid string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(updates map[string]interface{}, where string, args ...interface{}) error

	CreateUserTopics(data []*models.UserTopic) error
	DeleteUserTopics(userUid string) error
}

func (this *Dal) GetUser(preloads []string, where string, args ...interface{}) (*models.User, error) {
	user := models.User{}

	query := this.db.Model(&models.User{})

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Where(where, args...).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (this *Dal) GetUserProfileByUid(uid string) (*models.User, error) {

	if user, err := this.getUserProfileFromCache(uid); err != nil {
		if user, err = this.GetUser([]string{"UserTopics"}, "uid = ?", uid); err != nil {
			return nil, err
		} else {
			jsondata, _ := json.Marshal(user)
			if err = this.redis.Set(constants.GetUserProfileKey(uid), string(jsondata), time.Hour*constants.USER_PROFILE_TTL_IN_HOUR).Err(); err != nil {
				globals.GetLogger().Warn("[GetUserByUid] set to cache error: %+v", err)
			}

			return user, nil
		}
	} else {
		return user, nil
	}
}

func (this *Dal) getUserProfileFromCache(uid string) (*models.User, error) {
	var user models.User

	if result, err := this.redis.Get(constants.GetUserProfileKey(uid)).Result(); err != nil {
		return nil, err
	} else if err = json.Unmarshal([]byte(result), &user); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

func (this *Dal) CreateUser(user *models.User) error {
	return this.db.Model(&models.User{}).Create(user).Error
}

func (this *Dal) UpdateUser(updates map[string]interface{}, where string, args ...interface{}) error {

	if err := this.db.Model(&models.User{}).Where(where, args...).Updates(updates).Error; err != nil {
		return err
	}

	go func() {
		users := make([]models.User, 0)

		if err := this.db.Model(&models.User{}).Select("Uid").Where(where, args...).Find(&users).Error; err != nil {
			globals.GetLogger().Warn("[UpdateUser] db error: %+v", err)
		}

		keys := make([]string, len(users))
		for _, user := range users {
			keys = append(keys, constants.GetUserProfileKey(user.Uid))
		}
		if err := this.redis.Del(keys...).Err(); err != nil {
			globals.GetLogger().Warn("[UpdateUser] delete from cache error: %+v", err)
		}
	}()

	return nil
}

func (this *Dal) CreateUserTopics(data []*models.UserTopic) error {
	return this.db.Model(&models.UserTopic{}).CreateInBatches(data, len(data)).Error
}

func (this *Dal) DeleteUserTopics(userUid string) error {

	if err := this.db.Model(&models.UserTopic{}).Delete("user_uid = ?", userUid).Error; err != nil {
		return err
	}
	if err := this.redis.Del(constants.GetUserProfileKey(userUid)).Err(); err != nil {
		globals.GetLogger().Warn("[DeleteUserTopics] delete from cache error: %+v", err)
	}

	return nil
}
