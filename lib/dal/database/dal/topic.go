package dal

import "flaver/models"

type ITopicDal interface {
	GetTopicGroups() ([]*models.TopicGroup, error)
}

func (this *Dal) GetTopicGroups() ([]*models.TopicGroup, error) {

	topicGroups := make([]*models.TopicGroup, 0)

	if err := this.db.Model(&models.TopicGroup{}).Preload("Topics").Find(&topicGroups).Error; err != nil {
		return nil, err
	} else {
		return topicGroups, nil
	}

}
