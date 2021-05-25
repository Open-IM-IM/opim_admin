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

type Result struct {
	Count int
}

type QueryNewUserParam struct {
	OperationID *string `form:"operationID" json:"operationID" binding:"required"`
	BeginTime   *int64  `form:"beginTime" json:"beginTime" binding:"required"`
	EndTime     *int64  `form:"endTime" json:"endTime" binding:"required"`
}

func QueryNewUser(c *gin.Context) {
	token := c.Request.Header.Get("token")
	param := &QueryNewUserParam{}
	if c.Bind(param) != nil {
		c.JSON(http.StatusOK, gin.H{
			"ErrCode": ErrInfo.Json.ErrCode,
			"ErrMsg":  ErrInfo.Json.ErrMsg,
		})
		return
	}
	fmt.Printf("QueryNewUser param = %v", param)

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
		daysBeginTime := time.Unix(*param.BeginTime, 0).Format("2006-01-02 15:04:05")
		daysEndTime := time.Unix(*param.BeginTime+3600*24, 0).Format("2006-01-02 15:04:05")

		var result Result
		sql := fmt.Sprintf("select count(*) as count from user where created_time >= '%s' and created_time < '%s'", daysBeginTime, daysEndTime)
		dbConn.Raw(sql).Scan(&result)
		counts = append(counts, result.Count)

		*param.BeginTime = *param.BeginTime + 3600*24
	}

	var all Result
	sql := fmt.Sprintf("select count(*) as count from user ")
	dbConn.Raw(sql).Scan(&all)

	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	timeNumber := t.Unix() - 24*3600

	tm1 := time.Unix(timeNumber, 0)
	tm2 := time.Unix(t.Unix(), 0)

	var add Result
	sql = fmt.Sprintf("select count(*) as count from user where created_time >= '%s' and created_time < '%s'",
		tm1.Format("2006-01-02 15:04:05"), tm2.Format("2006-01-02 15:04:05"))
	dbConn.Raw(sql).Scan(&add)

	var active Result
	sql = fmt.Sprintf("select count(*) as count from chat_log where send_time >= '%s' and send_time <'%s'",
		tm1.Format("2006-01-02 15:04:05"), tm2.Format("2006-01-02 15:04:05"))
	dbConn.Raw(sql).Scan(&active)

	c.JSON(http.StatusOK, gin.H{
		"errorCode": 0,
		"errorMsg":  "",
		"data": gin.H{
			"count": counts,
			"subData": gin.H{
				"add":    add.Count,
				"active": active.Count,
				"all":    all.Count,
			},
		},
	})
}
