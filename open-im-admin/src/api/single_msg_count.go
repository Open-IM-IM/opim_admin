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

type SingleMsgCountParam struct {
	OptionID  *string `form:"optionID" json:"optionID" binding:"required"`
	BeginTime *string `form:"beginTime" json:"beginTime" binding:"required"`
	EndTime   *string `form:"endTime" json:"endTime" binding:"required"`
}

func SingleMsgCount(c *gin.Context) {
	token := c.Request.Header.Get("token")
	param := &QueryNewUserParam{}
	if c.Bind(param) != nil {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Json.ErrCode,
			"errorMsg":  ErrInfo.Json.ErrMsg,
		})
		return
	}
	fmt.Printf("SingleMsgCount param = %v", param)

	if token != utils.Md5(config.Config.Secret) {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Token.ErrCode,
			"errorMsg":  ErrInfo.Token.ErrMsg,
		})
		return
	}

	dbConn, _ := db.DB.DefaultGormDB()

	var counts []int
	for *param.BeginTime < *param.EndTime {
		nextTime := *param.BeginTime + 3600*24
		if nextTime > *param.EndTime {
			nextTime = *param.EndTime
		}

		var result Result
		sql := fmt.Sprintf("select count(*) as count from chat_log where msg_from = 0 and send_time >= '%s' and send_time < '%s'",
			time.Unix(*param.BeginTime, 0).Format("2006-01-02 15:04:05"), time.Unix(nextTime, 0).Format("2006-01-02 15:04:05"))
		dbConn.Raw(sql).Scan(&result)
		counts = append(counts, result.Count)

		*param.BeginTime = *param.BeginTime + 3600*24
	}

	var singleAll Result
	sql := fmt.Sprintf("select count(*) as count from chat_log where session_type = 1")
	dbConn.Raw(sql).Scan(&singleAll)

	c.JSON(http.StatusOK, gin.H{
		"errorCode": 0,
		"errorMsg":  "",
		"data": gin.H{
			"count": counts,
			"subData": gin.H{
				"singleAll": singleAll.Count,
			},
		},
	})
}
