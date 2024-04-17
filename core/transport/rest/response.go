package rest

import (
    "github.com/dulisoft/spirit/core/errorx/agcodes"
    "github.com/dulisoft/spirit/core/errorx/agerrors"
    "github.com/gin-gonic/gin"
    "net/http"
)

type HttpError struct {
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Solution    string      `json:"solution"`
	Cause       string      `json:"cause"`
	Detail      interface{} `json:"detail,omitempty"`
}

// success Json Response
func ResOKJson(c *gin.Context, data interface{}) {

	c.JSON(http.StatusOK, data)
}

// failed Json Response
func ResErrJson(c *gin.Context, err error) {
	var (
		code       = agerrors.Code(err)
		statusCode = 400
	)
	if err != nil {
		if code == agcodes.CodeNil {
			code = agcodes.CodeInternalError
		}
	} else if c.Writer.Status() > 0 && c.Writer.Status() != http.StatusOK {
		//switch c.Writer.Status() {
		//case http.StatusNotFound:
		//    code = agcodes.CodeNotFound
		//case http.StatusForbidden:
		//    code = agcodes.CodeNotAuthorized
		//
		//default:
		//    code = agcodes.CodeInternalError
		//}
		statusCode = c.Writer.Status()
	} else {
		code = agcodes.CodeOK
		statusCode = 200
	}

	c.JSON(statusCode, HttpError{
		Code:        code.GetErrorCode(),
		Description: code.GetDescription(),
		Solution:    code.GetSolution(),
		Cause:       code.GetCause(),
		Detail:      code.GetErrorDetails(),
	})
}
