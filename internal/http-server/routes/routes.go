package routes

import (
	"medods/internal/http-server/handlers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		Engine: gin.Default(),
	}
}

func (r *Router) SetupRouter(handlers handlers.AuthHandlers) {
	r.Engine.POST("/access", handlers.AccessHandler)
	r.Engine.POST("/refresh", handlers.RefreshHanadler)
}
