package event

import (
	"aculo/frontend-restapi/internal/config"
	"aculo/frontend-restapi/internal/server"
	eservice "aculo/frontend-restapi/internal/service/event"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Shared state of all */event/* routes

type EventGroup interface {
	server.Group
	makeEvent(gctx *gin.Context)
	getEvent(gctx *gin.Context)
}

func NewEventGroup(ctx context.Context, config config.Config, service eservice.EventService) server.Group {
	return &eventGroup{
		service: service,
	}
}

type eventGroup struct {
	service eservice.EventService
}

func (egroup *eventGroup) Attach(root *gin.RouterGroup) {
	_ = egroup.Chain(root)

}
func (g *eventGroup) Chain(root *gin.RouterGroup) *gin.RouterGroup {
	eGroup := root.Group("event")
	eGroup.POST("", g.makeEvent)
	eGroup.GET(":eid", g.getEvent)
	return eGroup

}
func (g *eventGroup) makeEvent(gctx *gin.Context) {
	panic("frontend shout not implement  post reqsts")

}
func (g *eventGroup) getEvent(gctx *gin.Context) {
	reqEID := gctx.Param("eid")
	if reqEID == "" {
		gctx.JSON(http.StatusBadRequest, gin.H{
			"message": "eid required",
		})
	}
	rsp, err := g.service.GetEvent(context.TODO(), eservice.GetEventRequest{
		EID: reqEID})

	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	gctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"event":   rsp.Event,
	})

}
