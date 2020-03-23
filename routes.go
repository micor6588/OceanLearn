package main

import (
	"micor/ginessential/controller"

	"github.com/gin-gonic/gin"
)

// CollectRoute 服务器的路由设置
func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	return r
}
