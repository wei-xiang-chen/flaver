package api

import (
	"errors"
	"flaver/lib/constants"
	"flaver/lib/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	DEFAULT_PAGE      = 1
	DEFAULT_PAGE_SIZE = 10
)

type PaginationResult[T interface{}] struct {
	Data       T   `json:"data"`
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
}

type HasNextResult[T interface{}] struct {
	Data   T    `json:"data"`
	NextId *int `json:"nextId"`
}

func Result(err *APIError, data interface{}, msg string, c *gin.Context) {
	if err.ErrorCode != SUCCESS {
		c.JSON(err.HttpStatus, map[string]interface{}{"error": ResponseWithError{
			Status: http.StatusText(err.HttpStatus),
			Code:   err.ErrorCode.String(),
			Msg:    msg,
		}})
	} else {
		c.JSON(http.StatusOK, data)
	}
}

const (
	// error code
	SUCCESS                           ErrorCode = "0"
	ERROR                             ErrorCode = "1"
	INVALID_ARGUMENT                  ErrorCode = "2"
	NOT_FOUND                         ErrorCode = "3"
	ALREADY_EXISTS                    ErrorCode = "4"
	PERMISSION_DENIED                 ErrorCode = "5"
	UNIMPLEMENTED                     ErrorCode = "6"
	TRANSACTION_RESULT_PARSING_FAILED ErrorCode = "7"

	TOPIC_COUNT_ILLEGAL ErrorCode = "000000"

	ID_TOKEN_PARSE_FAILED     ErrorCode = "100000"
	JWT_TOKEN_GENERATE_FAILED ErrorCode = "100001"
	JWT_TOKEN_PARSE_FAILED    ErrorCode = "100002"
)

var (
	Success                        = &APIError{ErrorCode: SUCCESS, Message: "SUCCESS", HttpStatus: http.StatusOK}
	Error                          = &APIError{ErrorCode: ERROR, Message: "ERROR", HttpStatus: http.StatusInternalServerError}
	InvalidArgument                = &APIError{ErrorCode: INVALID_ARGUMENT, Message: "INVALID_ARGUMENT", HttpStatus: http.StatusBadRequest}
	NotFound                       = &APIError{ErrorCode: NOT_FOUND, Message: "NOT_FOUND", HttpStatus: http.StatusNotFound}
	RecordNotFound                 = &APIError{ErrorCode: NOT_FOUND, Message: "RECORD_NOT_FOUND", HttpStatus: http.StatusNotFound}
	PermissionDenied               = &APIError{ErrorCode: PERMISSION_DENIED, Message: "PERMISSION_DENIED", HttpStatus: http.StatusForbidden}
	Unimplemented                  = &APIError{ErrorCode: UNIMPLEMENTED, Message: "UNIMPLEMENTED", HttpStatus: http.StatusNotImplemented}
	TransactionResultParsingFailed = &APIError{ErrorCode: TRANSACTION_RESULT_PARSING_FAILED, Message: "TRANSACTION_RESULT_PARSING_FAILED", HttpStatus: http.StatusInternalServerError}

	TopicCountIllegal = &APIError{ErrorCode: TOPIC_COUNT_ILLEGAL, Message: "TOPIC_COUNT_ILLEGAL", HttpStatus: http.StatusBadRequest}

	IdTokenParseFailed     = &APIError{ErrorCode: ID_TOKEN_PARSE_FAILED, Message: "ID_TOKEN_PARSE_FAILED", HttpStatus: http.StatusForbidden}
	JWTTokenGenerateFailed = &APIError{ErrorCode: JWT_TOKEN_GENERATE_FAILED, Message: "JWT_TOKEN_GENERATE_FAILED", HttpStatus: http.StatusInternalServerError}
	JWTTokenParseFailed    = &APIError{ErrorCode: JWT_TOKEN_PARSE_FAILED, Message: "JWT_TOKEN_PARSE_FAILED", HttpStatus: http.StatusForbidden}
)

type ResponseWithError struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Msg    string `json:"message"`
}

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

type APIError struct {
	ErrorCode  ErrorCode
	Message    string
	HttpStatus int
}

func (err *APIError) Error() string {
	return fmt.Sprintf("%s (%d)", err.Message, err.ErrorCode)
}

func SendResult(err error, data interface{}, c *gin.Context) {
	msg := "操作成功"
	if err != nil {
		if len(err.Error()) > 0 {
			msg = err.Error()
		} else {
			msg = "操作失败"
		}
	}
	if utils.GetRunTimeEnv() == constants.EnvProduction {
		msg = ""
	}
	Result(getErrorCodeStruct(err), data, msg, c)
}

func getErrorCodeStruct(err error) *APIError {
	if err == nil {
		return Success
	} else if len(err.Error()) == 0 {
		return Error
	}

	if IsArgCustomError(err) {
		return InvalidArgument
	}

	if apiError, ok := err.(*APIError); ok {
		return apiError
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return RecordNotFound
	}

	return Error
}

type ArgCustomError struct {
	CustomError error
	Message     string
}

func (e ArgCustomError) Error() string {
	return e.Message
}

func IsArgCustomError(err error) bool {
	_, ok := err.(*ArgCustomError)
	return ok
}
