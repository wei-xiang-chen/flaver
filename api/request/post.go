package request

type QueryPostInfoArg struct {
	Search *string `form:"search"`
	NextId *int    `form:"nextId"`
}

type CreatePostArg struct {
	UserUid       string   `json:"userUid"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	ImgUrls       []string `json:"imgUrls"`
	Rating        int      `json:"rating"`
	TopicIds      []int64  `json:"topicIds"`
	GooglePlaceId string   `json:"googlePlaceId"`
}

type UpdatePostArg struct {
	Title         *string   `json:"title"`
	Content       *string   `json:"content"`
	ImgUrls       *[]string `json:"imgUrls"`
	Rating        *int      `json:"rating"`
	TopicIds      *[]int64  `json:"topicIds"`
	GooglePlaceId *string   `json:"googlePlaceId"`
}

type PostLikeType string

const (
	PostLikeTypeLike   PostLikeType = "like"
	PostLikeTypeRevoke PostLikeType = "revoke"
)

type PostLikeArg struct {
	Type PostLikeType `json:"type"`
}
