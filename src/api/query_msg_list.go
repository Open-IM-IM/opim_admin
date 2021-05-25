package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"open-im-admin/src/common/config"
	"open-im-admin/src/common/db"
	ErrInfo "open-im-admin/src/common/error"
	"open-im-admin/src/utils"
	"strconv"
)

type chatMsg struct {
	SessionType  int
	MsgType      int
	SendID       string
	SendName     string
	RecvID       string
	RecvName     string
	QueryContent string
	Content      string
	Time         int
}

type QueryMsgListParam struct {
	OperationID  *string `form:"operationID" json:"operationID" binding:"required"`
	BeginTime    *int64  `form:"beginTime" json:"beginTime" `
	EndTime      *int64  `form:"endTime" json:"endTime" `
	SessionType  *int    `form:"sessionType" json:"sessionType" `
	MsgType      *int    `form:"msgType" json:"msgType" `
	QueryContent *string `form:"queryContent" json:"queryContent" `
}

func QueryMsgList(c *gin.Context) {
	token := c.Request.Header.Get("token")

	param := &QueryMsgListParam{}
	if c.Bind(param) != nil {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Json.ErrCode,
			"errorMsg":  ErrInfo.Json.ErrMsg,
		})
		return
	}
	fmt.Printf("QueryMsgList param = %v", param)

	if token != utils.Md5(config.Config.Secret) {
		c.JSON(http.StatusOK, gin.H{
			"errorCode": ErrInfo.Token.ErrCode,
			"errorMsg":  ErrInfo.Token.ErrMsg,
		})
		return
	}

	dbConn, _ := db.DB.DefaultGormDB()

	sql := "select session_type, content_type, send_id, recv_id, content, send_time from chat_log "

	if param.BeginTime != nil || param.EndTime != nil || param.SessionType != nil || param.MsgType != nil || param.QueryContent != nil {
		sql = sql + " where "
	}

	if param.BeginTime != nil {
		if param.EndTime != nil || param.SessionType != nil || param.MsgType != nil || param.QueryContent != nil {
			sql = sql + " send_time > " + strconv.Itoa(int(*param.BeginTime)) + " and "
		} else {
			sql = sql + " send_time > " + strconv.Itoa(int(*param.BeginTime))
		}

	}

	if param.EndTime != nil {
		if param.SessionType != nil || param.MsgType != nil || param.QueryContent != nil {
			sql = sql + " send_time < " + strconv.Itoa(int(*param.EndTime)) + " and "
		} else {
			sql = sql + " send_time < " + strconv.Itoa(int(*param.EndTime))
		}

	}
	if param.SessionType != nil {
		if param.MsgType != nil || param.QueryContent != nil {
			sql = sql + " session_type = " + strconv.Itoa(*param.SessionType) + " and "
		} else {
			sql = sql + " session_type = " + strconv.Itoa(*param.SessionType)
		}
	}
	if param.MsgType != nil {
		if param.QueryContent != nil {
			sql = sql + " content_type = " + strconv.Itoa(*param.MsgType) + " and "
		} else {
			sql = sql + " content_type = " + strconv.Itoa(*param.MsgType)
		}

	}
	if param.QueryContent != nil {
		sql = sql + " content = '" + *param.QueryContent + "'"
	}

	var arrChat []chatMsg
	rows, _ := dbConn.Raw(sql).Rows()
	for rows.Next() {
		var chat chatMsg
		if param.QueryContent != nil {
			chat.QueryContent = *param.QueryContent
		}
		rows.Scan(&chat.SessionType, &chat.MsgType, &chat.SendID, &chat.RecvID, &chat.Content, &chat.Time)
		arrChat = append(arrChat, chat)
	}

	var uidToName map[string]string
	uidToName = make(map[string]string)
	for _, v := range arrChat {
		uidToName[v.SendID] = ""
		uidToName[v.RecvID] = ""
	}

	for k := range uidToName {
		sql := fmt.Sprintf("select name from user where uid = '%s'", k)

		type NameResult struct {
			Name string
		}
		var name NameResult
		dbConn.Raw(sql).Scan(&name)
		uidToName[k] = name.Name
	}

	for k, v := range arrChat {
		arrChat[k].SendName = uidToName[v.SendID]
		arrChat[k].RecvName = uidToName[v.RecvID]
	}

	c.JSON(http.StatusOK, gin.H{
		"errorCode": 0,
		"errorMsg":  "",
		"data": gin.H{
			"msg": arrChat,
		},
	})
}
