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

type QueryUserParam struct {
	OperationID *string `form:"operationID" json:"operationID" binding:"required"`
	UID         string  `form:"uid" json:"uid" `
}

type QUser struct {
	UID              string
	Name             string
	Icon             string
	CreatedTime      string
	BeginDisableTime string
	EndDisableTime   string
	Seal             int
}

func QueryUser(c *gin.Context) {
	token := c.Request.Header.Get("token")
	param := &QueryUserParam{}
	if c.Bind(param) != nil {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Json.ErrCode,
			"errorMsg":  ErrInfo.Json.ErrMsg,
		})
		return
	}
	fmt.Printf("QueryUser param = %v", param)

	if token != utils.Md5(config.Config.Secret) {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Token.ErrCode,
			"errorMsg":  ErrInfo.Token.ErrMsg,
		})
		return
	}

	dbConn, _ := db.DB.DefaultGormDB()

	var result QUser
	var allQUser []QUser
	var sql string
	if len(param.UID) == 0 {
		sql = fmt.Sprintf("select uid, name, icon, created_time from user")
	} else {
		sql = fmt.Sprintf("select uid, name, icon, created_time from user where uid = '%s'", param.UID)
	}

	rows, _ := dbConn.Raw(sql).Rows()

	for rows.Next() {
		rows.Scan(&result.UID, &result.Name, &result.Icon, &result.CreatedTime)
		allQUser = append(allQUser, result)
	}

	sql = fmt.Sprintf("delete from black_list where end_disable_time <= '%s'", time.Now().Format("2006-01-02 15:04:05"))
	dbConn.Exec(sql)

	for k, v := range allQUser {
		sql = fmt.Sprintf("select begin_disable_time, end_disable_time from black_list where uid = '%s'", v.UID)
		dbConn.Raw(sql).Row().Scan(&allQUser[k].BeginDisableTime, &allQUser[k].EndDisableTime)
		if len(allQUser[k].BeginDisableTime) > 0 && len(allQUser[k].EndDisableTime) > 0 {
			allQUser[k].Seal = 1
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errorCode": 0,
		"errorMsg":  "",
		"data": gin.H{
			"user": allQUser,
		},
	})
}
