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

type BlackUser struct {
	UID              string
	Name             string
	BeginDisableTime string
	EndDisableTime   string
}

type QueryDisableUserParam struct {
	OperationID *string `form:"operationID" json:"operationID" binding:"required"`
}

func QueryDisableUser(c *gin.Context) {
	token := c.Request.Header.Get("token")
	param := &QueryDisableUserParam{}
	if c.Bind(param) != nil {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Json.ErrCode,
			"errorMsg":  ErrInfo.Json.ErrMsg,
		})
		return
	}
	fmt.Printf("QueryDisableUser param = %v", param)

	if token != utils.Md5(config.Config.Secret) {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Token.ErrCode,
			"errorMsg":  ErrInfo.Token.ErrMsg,
		})
		return
	}

	dbConn, _ := db.DB.DefaultGormDB()
	sql := fmt.Sprintf("delete from black_list where end_disable_time <= '%s'", time.Now().Format("2006-01-02 15:04:05"))
	dbConn.Exec(sql)

	var blackUser BlackUser
	var blackList []BlackUser
	rows, _ := dbConn.Raw("select uid, begin_disable_time, end_disable_time from black_list").Rows()
	for rows.Next() {
		rows.Scan(&blackUser.UID, &blackUser.BeginDisableTime, &blackUser.EndDisableTime)
		blackList = append(blackList, blackUser)
	}

	for k, v := range blackList {
		sql := fmt.Sprintf("select name from user where uid = '%s'", v.UID)

		type NameResult struct {
			Name string
		}
		var name NameResult
		dbConn.Raw(sql).Scan(&name)
		blackList[k].Name = name.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"errorCode": 0,
		"errorMsg":  "",
		"data": gin.H{
			"users": blackList,
		},
	})
}
