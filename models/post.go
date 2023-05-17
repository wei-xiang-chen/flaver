package models

import (
	"flaver/api/response"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	Id            int            `gorm:"column:id; primary_key"`
	Title         string         `gorm:"column:title"`
	Content       string         `gorm:"column:content"`
	ImgUrls       pq.StringArray `gorm:"column:img_urls; type:text[]"`
	Rating        int            `gorm:"column:rating"`
	TopicIds      pq.Int64Array  `gorm:"column:topic_ids; type:integer[]"`
	OwnerUid      string         `gorm:"column:owner_uid"`
	GooglePlaceId string         `gorm:"column:google_place_id"`
	InstantlyRole string         `gorm:"column:instantly_role"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at"`

	Owner      *User       `gorm:"foreignKey:OwnerUid; references:Uid"`
	Restaurant *Restaurant `gorm:"foreignKey:GooglePlaceId; references:PlaceId"`
	PostLikes  []*PostLike `gorm:"foreignKey:PostId; references:Id"`
}

func (Post) TableName() string {
	return "posts"
}

func (this *Post) SerializeTo(buffer interface{}, postsLikeCount map[int]int) bool {
	if theBuffer, ok := buffer.(*response.PostList); ok {
		theBuffer.Id = this.Id
		theBuffer.Title = this.Title
		theBuffer.ImgUrls = this.ImgUrls

		if this.Restaurant != nil {
			theBuffer.RestaurantName = this.Restaurant.Name
		}

		if this.Owner != nil {
			this.Owner.SerializeTo(&theBuffer.Owner)
		}

		if len(this.PostLikes) > 0 {
			theBuffer.HasLiked = true
		}

		if val, ok := postsLikeCount[this.Id]; ok {
			theBuffer.LikeCount = val
		}

		return true
	} else if theBuffer, ok := buffer.(*response.PostDetail); ok {
		theBuffer.Id = this.Id
		theBuffer.Title = this.Title
		theBuffer.ImgUrls = this.ImgUrls
		theBuffer.Content = this.Content
		theBuffer.Rating = this.Rating
		theBuffer.TopicIds = this.TopicIds
		if this.Restaurant != nil {
			this.Restaurant.SerializeTo(&theBuffer.Restaurant)
		}
		if this.Owner != nil {
			this.Owner.SerializeTo(&theBuffer.Owner)
		}

		if len(this.PostLikes) > 0 {
			theBuffer.HasLiked = true
		}

		if val, ok := postsLikeCount[this.Id]; ok {
			theBuffer.LikeCount = val
		}

		return true
	}

	return false
}
