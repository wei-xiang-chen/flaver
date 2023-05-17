package request

import "time"

type LoginArg struct {
	GoogleIdToken *string `json:"googleIdToken"`
	AppleIdToken  *string `json:"appleIdToken"`
}

type RefreshTokenArg struct {
	RefreshToken string `json:"refreshToken"`
}

type RegisteArg struct {
	GoogleIdToken *string `json:"googleIdToken"`
	AppleIdToken  *string `json:"appleIdToken"`
	TopicIds      []int   `json:"topicIds"`
	Nickname      string  `json:"nickname"`
	AvatarImgUrl  string  `json:"avatarImgUrl"`
}

type UpdateUserArg struct {
	Nickname     *string    `json:"nickname"`
	Birthday     *time.Time `json:"birthday"`
	TopicIds     []int      `json:"topicIds"`
	AvatarImgUrl *string    `json:"avatarImgUrl"`
}
