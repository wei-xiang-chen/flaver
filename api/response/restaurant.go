package response

type Restaurant struct {
	GooglePlaceId string  `json:"googlePlaceId"`
	Name          string  `json:"name"`
	Phone         *string `json:"phone"`
	Address       string  `json:"address"`
	Lat           float64 `json:"lat"`
	Lng           float64 `json:"lng"`
	GeneralRating float64 `json:"generalRating"`
	FlaverRating  float64 `json:"flaverRating"`
}
