package models

import (
	"flaver/api/response"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	GeneralRole UserRole = "general"
	FlaverRole  UserRole = "flaver"
)

type User struct {
	Uid           string         `gorm:"column:uid; primary_key"`
	GoogleIdToken *string        `gorm:"column:google_id_token"`
	AppleIdToken  *string        `gorm:"column:apple_id_token"`
	Nickname      string         `gorm:"column:nickname"`
	AvatarImgUrl  string         `gorm:"column:avatar_img_url"`
	Birthday      time.Time      `gorm:"column:birthday"`
	Email         string         `gorm:"column:email"`
	Role          UserRole       `gorm:"column:role"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`

	UserTopics []*UserTopic `gorm:"foreignKey:Uid; references:UserUid"`
}

func (User) TableName() string {
	return "users"
}

func (this *User) SerializeTo(buffer interface{}) bool {
	if theBuffer, ok := buffer.(*response.User); ok {
		theBuffer.Uid = this.Uid
		theBuffer.Nickname = this.Nickname
		theBuffer.AvatarImgUrl = this.AvatarImgUrl

		return true
	} else if theBuffer, ok := buffer.(*response.UserDetail); ok {
		theBuffer.Uid = this.Uid
		theBuffer.Nickname = this.Nickname
		theBuffer.AvatarImgUrl = this.AvatarImgUrl
		theBuffer.Birthday = this.Birthday
		theBuffer.Email = this.Email

		topicIds := []*int{}
		for id := range this.UserTopics {
			topicIds = append(topicIds, &id)
		}
		theBuffer.TopicIds = topicIds

		return true
	}

	return false
}
