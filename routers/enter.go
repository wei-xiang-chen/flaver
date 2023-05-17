package routers

type RouterGroup struct {
	PostGroup  PostRouter
	UserGroup  UserRouter
	FileGroup  FileRouter
	TopicGroup TopicRouter
}

var (
	RouterGroupApp = new(RouterGroup)
)
