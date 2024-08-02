package controller

import (
	"aculo/connector-restapi/internal/config"
	"aculo/connector-restapi/internal/controller/groups"
	"context"

	"github.com/gin-gonic/gin"
)

func newMux(
	ctx context.Context,
	config config.Config,
	engine *gin.Engine,
	groups []groups.Attachable,
	chains []groups.Chain,
) groups.AppMux {
	attachEndpoints(&engine.RouterGroup, groups)
	attachChainedEndpoints(&engine.RouterGroup, chains)
	return engine
}
func attachEndpoints(rootGroup *gin.RouterGroup, groups []groups.Attachable) {
	for _, g := range groups {
		g.Attach(rootGroup)
	}
}
func attachChainedEndpoints(rootGroup *gin.RouterGroup, chains []groups.Chain) {
	for _, chain := range chains {
		todo := rootGroup
		for _, elem := range chain {
			todo = elem.Chain(todo)
		}
	}
}
