package user

import (
	"flaver/api/request"
	"flaver/api/response"
	"flaver/models"
)

type IUserService interface {
	Login(arg *request.LoginArg) (*response.Token, error)
	RefreshToken(userUid, refreshToken string)(*response.Token, error)
	RegistUser(data *request.RegisteArg) (*response.Token, error)
	GetUser(uid string) (*models.User, error)
	Update(uid string, data *request.UpdateUserArg) error
}
