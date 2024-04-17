package ginx

import (
    "github.com/dulisoft/spirit/core/errorx/agcodes"
    "github.com/dulisoft/spirit/core/errorx/agerrors"
    "github.com/gin-gonic/gin"
    "net/http"
)

type HttpError struct {
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Solution    string      `json:"solution,omitempty"`
	Cause       string      `json:"cause,omitempty"`
	Detail      interface{} `json:"detail,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

// success Json Response
func ResOKJson(c *gin.Context, data interface{}) {

	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, data)
}

// list Response
func ResList(c *gin.Context, list interface{}, totalCount int) {

	c.JSON(http.StatusOK, gin.H{
		"entries":     list,
		"total_count": totalCount,
	})

}

// failed Json Response
func ResErrJson(c *gin.Context, err error) {

	var (
		code = agerrors.Code(err)
	)
	if err != nil {
		if code == agcodes.CodeNil {
			code = agcodes.CodeInternalError
		}
	} else if c.Writer.Status() > 0 && c.Writer.Status() != http.StatusOK {
		switch c.Writer.Status() {
		case http.StatusNotFound:
			code = agcodes.CodeNotFound
		case http.StatusForbidden:
			code = agcodes.CodeNotAuthorized
		default:
			code = agcodes.CodeInternalError
		}
	} else {
		code = agcodes.CodeOK
	}

	c.JSON(c.Writer.Status(), HttpError{
		Code:        code.GetErrorCode(),
		Description: code.GetDescription(),
		Solution:    code.GetSolution(),
		Cause:       code.GetCause(),
		Detail:      code.GetErrorDetails(),
	})
}
