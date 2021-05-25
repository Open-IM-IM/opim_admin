package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"open-im-admin/src/common/config"
	"open-im-admin/src/common/db"
	ErrInfo "open-im-admin/src/common/error"
	"open-im-admin/src/utils"
	"time"
)

type DisableUserParam struct {
	OperationID   *string `form:"operationID" json:"operationID" binding:"required"`
	UID           string  `form:"uid" json:"uid" binding:"required"`
	DisableSecond *int64  `form:"disable_second" json:"disable_second" binding:"required"`
	EX            *string `form:"ex" json:"ex" binding:"required"`
}

func DisableUser(c *gin.Context) {
	token := c.Request.Header.Get("token")
	param := &DisableUserParam{}
	if c.Bind(param) != nil {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Json.ErrCode,
			"errorMsg":  ErrInfo.Json.ErrMsg,
		})
		return
	}
	fmt.Printf("DisableUser param = %v", param)

	if token != utils.Md5(config.Config.Secret) {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Token.ErrCode,
			"errorMsg":  ErrInfo.Token.ErrMsg,
		})
		return
	}

	dbConn, _ := db.DB.DefaultGormDB()
	if *param.DisableSecond <= 0 {
		sql := fmt.Sprintf("delete from black_list where uid = '%s'", param.UID)
		dbConn.Exec(sql)
	} else {
		beginDisableSecond := time.Now().Unix()
		endDisableSecond := beginDisableSecond + *param.DisableSecond
		tmBeginDisableTime := time.Unix(beginDisableSecond, 0)
		tmEndDisableTime := time.Unix(endDisableSecond, 0)
		dbConn.Exec("replace into black_list(uid, begin_disable_time, end_disable_time, ex) values(?, ?, ?, ?)",
			param.UID, tmBeginDisableTime.Format("2006-01-02 15:04:05"), tmEndDisableTime.Format("2006-01-02 15:04:05"), *param.EX)
	}

	c.JSON(http.StatusOK, gin.H{
		"errorCode": 0,
		"errorMsg":  "",
	})
}
