package controller

import (
	"aculo/frontend-restapi/internal/config"
	"aculo/frontend-restapi/internal/controller/groups"
	"context"

	"github.com/gin-gonic/gin"
)

func NewMux(
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
		chainpart := rootGroup
		for _, elem := range chain {
			chainpart = elem.Chain(chainpart)
		}
	}
}
