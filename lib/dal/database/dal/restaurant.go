package dal

import (
	"flaver/models"
	"fmt"
)

type IRestaurantDal interface {
	GetRestaurantRating(placeIds []string) (map[string]*RestaurantRating, error)
	CheckRestaurantIsExist(placeId string) (bool, error)
	CreateRestaurant(restaurant *models.Restaurant) error
}

type RestaurantRating struct {
	GeneralRating float64
	FlaverRating  float64
}

func (this *Dal) GetRestaurantRating(placeIds []string) (map[string]*RestaurantRating, error) {
	var results []map[string]interface{}
	if err := this.db.Model(&models.Post{}).Select("google_place_id, instantly_role, SUM(rating)::float/COUNT(*) AS average_rating").Where("google_place_id IN ?", placeIds).Group("google_place_id, instantly_role").Find(&results).Error; err != nil {
		return nil, err
	}

	resultMap := map[string]*RestaurantRating{}
	for _, result := range results {

		if _, ok := resultMap[result["google_place_id"].(string)]; !ok {
			resultMap[result["google_place_id"].(string)] = &RestaurantRating{}
		}

		restaurantRating := resultMap[result["google_place_id"].(string)]
		role := result["instantly_role"].(string)
		if role == string(models.GeneralRole) {
			restaurantRating.GeneralRating = result["average_rating"].(float64)
		} else if role == string(models.FlaverRole) {
			restaurantRating.FlaverRating = result["average_rating"].(float64)
		}
		fmt.Println(restaurantRating)
	}

	return resultMap, nil
}

func (this *Dal) CheckRestaurantIsExist(placeId string) (bool, error) {
	var count int64
	if err := this.db.Model(&models.Restaurant{}).Where("place_id = ?", placeId).Count(&count).Error; err != nil {
		return false, err
	} else if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (this *Dal) CreateRestaurant(restaurant *models.Restaurant) error {
	return this.db.Model(&models.Restaurant{}).Create(restaurant).Error
}
