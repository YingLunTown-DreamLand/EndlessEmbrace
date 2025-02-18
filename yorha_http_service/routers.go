package yorha_http_service

import (
	bot "EndlessEmbrace"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	GinEngine *gin.Engine
	BotClient *bot.Client
}

func NewRouter(botClient *bot.Client) *Router {
	router := &Router{
		GinEngine: gin.Default(),
	}

	qqAuthGroup := router.GinEngine.Group("/qq_auth")
	{
		qqAuthGroup.POST("/get_group_member_info", router.GetGroupMemberInfo)
		qqAuthGroup.POST("/notify_to_all_member", router.NotifyToAllMember)
	}

	// No router
	router.GinEngine.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	return router
}
