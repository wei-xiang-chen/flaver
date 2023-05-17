package response

import "time"

type User struct {
	Uid          string `json:"uid"`
	Nickname     string `json:"nickname"`
	AvatarImgUrl string `json:"avatarImgUrl"`
}

type UserDetail struct {
	Birthday time.Time `json:"birthday"`
	Email    string    `json:"email"`
	TopicIds []*int    `json:"topicIds"`
	User
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
