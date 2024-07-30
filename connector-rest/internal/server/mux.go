package server

import (
	"aculo/connector-restapi/internal/config"
	"context"

	"github.com/gin-gonic/gin"
)

func NewMux(ctx context.Context, config config.Config, engine *gin.Engine, groups []Attachable, chains []Chain) AppMux {
	attachEndpoints(&engine.RouterGroup, groups)
	attachChainedEndpoints(&engine.RouterGroup, chains)
	return engine
}
func attachEndpoints(rootGroup *gin.RouterGroup, groups []Attachable) {
	for _, g := range groups {
		g.Attach(rootGroup)
	}
}
func attachChainedEndpoints(rootGroup *gin.RouterGroup, chains []Chain) {
	for _, chain := range chains {
		todo := rootGroup
		for _, elem := range chain {
			todo = elem.Chain(todo)
		}
	}
}
