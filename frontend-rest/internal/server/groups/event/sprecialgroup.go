package event

import (
	"aculo/frontend-restapi/internal/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type specialGroup struct {
}

func NewSpecialGroup() server.Group {
	return &specialGroup{}
}
func (g *specialGroup) getSpecial(gctx *gin.Context) {
	gctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})

}
func (g *specialGroup) Attach(root *gin.RouterGroup) {
	_ = g.Chain(root)

}
func (g *specialGroup) Chain(root *gin.RouterGroup) *gin.RouterGroup {
	specialGroup := &specialGroup{}
	sGroup := root.Group("special")
	sGroup.GET("", specialGroup.getSpecial)
	return sGroup
}
