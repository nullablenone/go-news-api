package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Error:   nil,
	})
}

func Error(c *gin.Context, code int, message string, errCode string) {
	c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
		Error:   errCode,
	})
}

func AbortError(c *gin.Context, code int, message string, errCode string) {
	c.AbortWithStatusJSON(code, APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
		Error:   errCode,
	})
}
