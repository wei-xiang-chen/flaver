package topic

import "flaver/models"

type ITopicService interface {
	GetTopicGroups() ([]*models.TopicGroup, error)
}