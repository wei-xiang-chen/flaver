package models

import "time"

type PostLike struct {
	UserUid   string    `gorm:"column:user_uid; primary_key"`
	PostId    int       `gorm:"column:post_id; primary_key"`
	CreatedAt time.Time `gorm:"column:created_at"`

	User *User `gorm:"constrain:foreignKey:UserUid; references:Uid"`
	Post *Post `gorm:"constrain:foreignKey:PostId; references:Id"`
}

func (PostLike) TableName() string {
	return "post_likes"
}
