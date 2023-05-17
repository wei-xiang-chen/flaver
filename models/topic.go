package models

import "flaver/api/response"

type Topic struct {
	Id           int    `gorm:"column:id; primary_key"`
	Name         string `gorm:"column:name"`
	ImgUrl       string `gorm:"column:img_url"`
	TopicGroupId int    `gorm:"column:topic_group_id"`

	TopicGroup TopicGroup `gorm:"foreignKey:TopicGroupId; references:Id"`
}

func (Topic) TableName() string {
	return "topics"
}

type TopicGroup struct {
	Id     int    `gorm:"column:id;primary_key"`
	Name   string `gorm:"column:name"`
	ImgUrl string `gorm:"column:img_url"`

	Topics []*Topic `gorm:"foreignKey:Id; references:TopicGroupId"`
}

func (TopicGroup) TableName() string {
	return "topic_groups"
}

func (this *Topic) SerializeTo(buffer interface{}) bool {

	if theBuffer, ok := buffer.(*response.Topic); ok {
		theBuffer.Id = this.Id
		theBuffer.Name = this.Name
		theBuffer.ImgUrl = this.ImgUrl

		return true
	}

	return false
}

func (this *TopicGroup) SerializeTo(buffer interface{}) bool {

	if theBuffer, ok := buffer.(*response.TopicGroup); ok {
		theBuffer.Id = this.Id
		theBuffer.Name = this.Name
		theBuffer.ImgUrl = this.ImgUrl

		topics := make([]*response.Topic, 0)
		for _, topic := range this.Topics {
			respTopic := response.Topic{}
			topic.SerializeTo(respTopic)
			topics = append(topics, &respTopic)
		}
		theBuffer.Topics = topics

		return true
	}

	return false
}
