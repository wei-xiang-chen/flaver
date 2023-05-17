package constants

import "fmt"

const (
	uSER_PROFILE_KEY         = "USER:PROFILE:%s"
	USER_PROFILE_TTL_IN_HOUR = 24

	POST_LIKE_COUNT_KEY = "POST:LIKE:COUNT"
)

func GetUserProfileKey(uid string) string {
	return fmt.Sprintf(uSER_PROFILE_KEY, uid)
}
