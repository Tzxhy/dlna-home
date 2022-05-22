package utils

import (
	"net/http"

	"gitee.com/tzxhy/dlna-home/constants"
	"github.com/gin-gonic/gin"
)

func ReturnJSON(code int, message string, data *gin.H) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
}

func ReturnParamNotValid(c *gin.Context) {
	c.JSON(http.StatusOK, ReturnJSON(constants.CODE_PARAMS_NOT_VALID_TIPS.Code, constants.CODE_PARAMS_NOT_VALID_TIPS.Tip, nil))
}
