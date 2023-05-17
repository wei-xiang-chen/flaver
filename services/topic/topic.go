package topic

import (
	"flaver/lib/dal/database/dal"
	"flaver/models"
)

type TopicService struct {
	topicDal dal.ITopicDal
}

type TopicServiceOption func(*TopicService)

func NewTopicServiceOption(options ...func(*TopicService)) ITopicService {
	service := TopicService{}

	for _, option := range options {
		option(&service)
	}

	return &service
}

func WithTopicDal(dal dal.ITopicDal) TopicServiceOption {
	return func(service *TopicService) {
		service.topicDal = dal
	}
}

func (this *TopicService) GetTopicGroups() ([]*models.TopicGroup, error) {
	return this.topicDal.GetTopicGroups()
}
