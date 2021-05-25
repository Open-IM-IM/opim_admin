package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"open-im-admin/src/common/config"
	ErrInfo "open-im-admin/src/common/error"
	"open-im-admin/src/utils"
)

type TokenParam struct {
	OperationID *string `form:"operationID" json:"operationID" binding:"required"`
	Secret      string  `form:"secret" json:"secret" binding:"required"`
}

func GetToken(c *gin.Context) {
	param := &TokenParam{}
	if c.Bind(param) != nil {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Json.ErrCode,
			"errorMsg":  ErrInfo.Json.ErrMsg,
		})
		return
	}
	fmt.Printf("GetToken param = %v", param)

	if param.Secret != config.Config.Secret {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Secret.ErrCode,
			"errorMsg":  ErrInfo.Secret.ErrMsg,
		})
		return
	}

	tokenString := utils.Md5(config.Config.Secret)

	c.JSON(http.StatusOK, gin.H{
		"errorCode": 0,
		"errorMsg":  "",
		"data": gin.H{
			"token": tokenString,
		},
	})
}
