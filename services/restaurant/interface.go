package restaurant

import "flaver/lib/dal/database/dal"

type IRestaurantService interface {
	GetRestaurantRating(placeIds []string) (map[string]*dal.RestaurantRating, error)
}