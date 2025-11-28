package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	CodeSuccess      = 200
	CodeServerError  = 500
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
)

func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Success(c *gin.Context, data interface{}) {
	Result(c, CodeSuccess, "success", data)
}

func Error(c *gin.Context, code int, msg string) {
	Result(c, code, msg, nil)
}

func ServerError(c *gin.Context, err error) {
	Result(c, CodeServerError, err.Error(), nil)
}

func BadRequest(c *gin.Context, msg string) {
	Result(c, CodeBadRequest, msg, nil)
}

func Page(c *gin.Context, list interface{}, total int64) {
	Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
