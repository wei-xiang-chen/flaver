package response

type Topic struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	ImgUrl string `json:"imgUrl"`
}

type TopicGroup struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	ImgUrl string `json:"imgUrl"`

	Topics []*Topic `json:"topics"`
}
