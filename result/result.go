package common

import (
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Result API返回结果
type Result struct {
	ErrorCode string `json:"code"`
	Payload   any    `json:"payload,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (r Result) Error() string {
	return r.Message
}

func Fails(err error) *Result {
	return &Result{
		ErrorCode: "UNEXPECTED_ERROR",
		Message:   err.Error(),
	}
}

func FailsWithPayLoad(payload any, err error) *Result {
	result := Fails(err)
	result.Payload = payload
	return result
}

// FailsWithErrorCode 通过错误码返回国际化资源
func FailsWithErrorCode(ctx *gin.Context, errorCode string, params map[string]string) *Result {
	var content string
	if params != nil {
		content = ginI18n.MustGetMessage(ctx, &i18n.LocalizeConfig{
			MessageID:    errorCode,
			TemplateData: params,
		})
	} else {
		content = ginI18n.MustGetMessage(ctx, errorCode)
	}

	return &Result{
		ErrorCode: errorCode,
		Message:   content,
	}
}
