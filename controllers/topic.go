package controllers

import (
	"flaver/api"
	"flaver/api/response"
	"flaver/globals"
	"flaver/lib/dal"
	"flaver/services/topic"

	"github.com/gin-gonic/gin"
)

type TopicController struct {
	topicService topic.ITopicService

	transactionContext dal.TransactionContext
}

func NewTopicController() TopicController {
	dal := dal.NewDal()

	return TopicController{
		topicService: topic.NewTopicServiceOption(
			topic.WithTopicDal(dal.Dal),
		),
		transactionContext: dal.TransactionContext,
	}
}

func (this TopicController) GetList(c *gin.Context) {

	result, err := this.topicService.GetTopicGroups()
	if err != nil {
		globals.GetLogger().Warnf("[GetTopicList] error: %v", err)
		api.SendResult(err, nil, c)
		return
	}

	resultData := make([]*response.TopicGroup, 0)

	for _, data := range result {
		topicGroup := response.TopicGroup{}
		data.SerializeTo(topicGroup)
		resultData = append(resultData, &topicGroup)
	}

	api.SendResult(nil, resultData, c)
}
