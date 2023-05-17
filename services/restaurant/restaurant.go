package restaurant

import "flaver/lib/dal/database/dal"

type RestaurantService struct {
	restaurantDal dal.IRestaurantDal
}

type RestaurantServiceOption func(*RestaurantService)

func NewRestaurantServiceOption(options ...func(*RestaurantService)) IRestaurantService {
	service := RestaurantService{}

	for _, option := range options {
		option(&service)
	}

	return &service
}

func WithRestaurantDal(dal dal.IRestaurantDal) RestaurantServiceOption {
	return func(service *RestaurantService) {
		service.restaurantDal = dal
	}
}

func (this *RestaurantService) GetRestaurantRating(placeIds []string) (map[string]*dal.RestaurantRating, error) {
	return this.restaurantDal.GetRestaurantRating(placeIds)
}
