package event

import (
	"aculo/connector-restapi/internal/config"
	"aculo/connector-restapi/internal/controller/groups"
	"aculo/connector-restapi/internal/service"
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Shared state of all */event/* routes

/* type EventGroup interface {
	server.Group
	makeEvent(gctx *gin.Context)
	getEvent(gctx *gin.Context)
} */

func NewEventGroup(ctx context.Context, config config.Config, service service.Service) groups.Group {
	return &eventGroup{
		service: service,
	}
}

type eventGroup struct {
	service service.Service
}

func (g *eventGroup) Attach(root *gin.RouterGroup) {
	_ = g.Chain(root)

}
func (g *eventGroup) Chain(root *gin.RouterGroup) *gin.RouterGroup {
	router := root.Group("event")
	router.POST("/", g.sendSingleEvent)
	return router

}

// @Summary      Send single event via REST
// @Description  Least preferable variant, but still works
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event   body       string   true  "Raw event to send"
// @Param        topic   query      string   true  "Topic to send event"
// @Success      200  {object}  Rsp200
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /event/ [post]
func (g *eventGroup) sendSingleEvent(gctx *gin.Context) {
	body, err := io.ReadAll(gctx.Request.Body)
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{
			"status": err.Error(),
		})
		return
	}
	topic := gctx.Query("topic")
	if topic == "" {
		gctx.JSON(http.StatusBadRequest, gin.H{
			"status": "topic is required",
		})
		return
	}
	// TODO : resp may be useless here
	_, err = g.service.SendEvent(context.TODO(), service.SendEventRequest{
		Topic: topic,
		Event: body,
	})
	if err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{
			"status": err.Error(),
		})
		return
	}

	gctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}

type Rsp200 struct {
	Status string `json:"status"`
}
type Error struct {
	Status string `json:"status"`
}
