package models

import "time"

type UserTopic struct {
	UserUid   string    `gorm:"column:user_uid; primary_key"`
	TopicId   int       `gorm:"column:topic_id; primary_key"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserTopic) TableName() string {
	return "user_topics"
}
