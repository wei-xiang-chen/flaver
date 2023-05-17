package response

type PostList struct {
	Id             int      `json:"id"`
	Title          string   `json:"title"`
	ImgUrls        []string `json:"imgUrls"`
	RestaurantName string   `json:"restaurantName"`
	LikeCount      int      `json:"likeCount"`
	HasLiked       bool     `json:"hasLiked"`
	HasCollected   bool     `json:"hasCollected"`
	Owner          User     `json:"owner"`
}

type PostDetail struct {
	Id         int        `json:"id"`
	Title      string     `json:"title"`
	ImgUrls    []string   `json:"imgUrls"`
	Content    string     `json:"content"`
	LikeCount  int        `json:"likeCount"`
	HasLiked   bool       `json:"hasLiked"`
	TopicIds   []int64    `json:"topicIds"`
	Rating     int        `json:"rating"`
	Restaurant Restaurant `json:"restaurant"`
	Owner      User       `json:"owner"`
}
