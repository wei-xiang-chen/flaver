package dal

import (
	"flaver/lib/dal/database"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Dal struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewDal() IDal {
	dal := &Dal{
		db:    database.GetClientDB(),
		redis: database.GetRedisClient(),
	}
	return dal
}

func (this *Dal) SetConn(db *gorm.DB) {
	this.db = db
}

func (this *Dal) GetConn() *gorm.DB {
	return this.db
}

type IDal interface {
	database.GormGetSettable

	IRedis

	IUserDal
	IPostDal
	IRestaurantDal
	ITopicDal
}
