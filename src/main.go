package main

import (
	"github.com/gin-gonic/gin"
	"open-im-admin/src/api"
	"open-im-admin/src/common/config"
	"open-im-admin/src/utils"
)

func main() {
	r := gin.Default()
	r.Use(utils.CorsHandler())
	r.POST("/get_token", api.GetToken)
	r.POST("/disable_user", api.DisableUser)
	r.POST("/query_disable_user", api.QueryDisableUser)
	r.POST("/query_user", api.QueryUser)
	r.POST("/single_msg_count", api.SingleMsgCount)
	r.POST("/query_new_user", api.QueryNewUser)
	r.POST("/query_msg_list", api.QueryMsgList)
	r.Run(config.Config.Api)
}
