package models

import (
	"flaver/api/response"
	"time"

	"gorm.io/gorm"
)

type Restaurant struct {
	GooglePlaceId string  `gorm:"column:place_id; primary_key"`
	Name          string  `gorm:"column:name"`
	Phone         *string `gorm:"column:phone"`
	Address       string  `gorm:"column:address"`
	// Location      string         `gorm:"column:location; type:geography(point,4326)"`
	Lat       float64        `gorm:"column:lat"`
	Lng       float64        `gorm:"column:lng"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

func (this *Restaurant) SerializeTo(buffer interface{}) bool {
	if theBuffer, ok := buffer.(*response.Restaurant); ok {
		theBuffer.GooglePlaceId = this.GooglePlaceId
		theBuffer.Name = this.Name
		theBuffer.Phone = this.Phone
		theBuffer.Address = this.Address
		theBuffer.Lat = this.Lat
		theBuffer.Lng = this.Lng

		return true
	}

	return false
}
