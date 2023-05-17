package controllers

import (
	"flaver/api"
	"flaver/api/request"
	"flaver/api/response"
	"flaver/globals"
	"flaver/lib/dal"
	"flaver/lib/utils"
	"flaver/services/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	userService user.IUserService

	transactionContext dal.TransactionContext
}

func NewUserController() UserController {
	dal := dal.NewDal()

	return UserController{
		userService: user.NewUserServiceOption(
			user.WithUserDal(dal.Dal),
			user.WithFirebaseAuthUtil(&utils.FirebaseAuthUtil{}),
		),

		transactionContext: dal.ExecTransaction,
	}
}

func (this UserController) Login(c *gin.Context) {

	arg := request.LoginArg{}
	if err := c.ShouldBindJSON(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	result, err := this.userService.Login(&arg)
	if err != nil {
		api.SendResult(err, nil, c)
		return
	} else {
		api.SendResult(nil, result, c)
	}
}

func (this UserController) RefreshToken(c *gin.Context) {

	userUid := c.MustGet("userUid").(string)

	arg := request.RefreshTokenArg{}
	if err := c.ShouldBindJSON(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	result, err := this.userService.RefreshToken(userUid, arg.RefreshToken)
	if err != nil {
		api.SendResult(err, nil, c)
		return
	} else {
		api.SendResult(nil, result, c)
	}
}

func (this UserController) RegistUser(c *gin.Context) {

	arg := request.RegisteArg{}
	if err := c.ShouldBindJSON(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	transaction := func(tx *gorm.DB) (interface{}, error) {
		result, err := this.userService.RegistUser(&arg)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	if result, err := this.transactionContext(transaction); err != nil {
		globals.GetLogger().Warnf("[RegistUser] transaction error: %v", err)
		api.SendResult(err, nil, c)
		return
	} else if result, ok := result.(*response.Token); !ok {
		api.SendResult(api.TransactionResultParsingFailed, nil, c)
	} else {
		api.SendResult(nil, result, c)
	}
}

func (this UserController) GetDetail(c *gin.Context) {
	userUid := c.Param("userUid")

	user, err := this.userService.GetUser(userUid)
	if err != nil {
		globals.GetLogger().Warnf("[GetUserDetail] error: %v", err)
		api.SendResult(err, nil, c)
		return
	}

	result := &response.UserDetail{}
	user.SerializeTo(result)

	api.SendResult(nil, result, c)
}

func (this UserController) Update(c *gin.Context) {
	userUid := c.MustGet("userUid").(string)

	arg := request.UpdateUserArg{}
	if err := c.ShouldBindJSON(&arg); err != nil {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	transaction := func(tx *gorm.DB) (interface{}, error) {
		err := this.userService.Update(userUid, &arg)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
	if _, err := this.transactionContext(transaction); err != nil {
		globals.GetLogger().Warnf("[UpdateUser] transaction error: %v", err)
		api.SendResult(err, nil, c)
		return
	} else {
		api.SendResult(nil, nil, c)
	}
}
