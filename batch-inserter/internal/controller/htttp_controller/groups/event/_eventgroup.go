package event

import (
	//"aculo/frontend-restapi/domain"
	//"aculo/frontend-restapi/internal/config"
	//"aculo/frontend-restapi/internal/controller/groups"
	//"aculo/frontend-restapi/internal/request"
	//eservice "aculo/frontend-restapi/internal/service/event"
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/controller"
	"aculo/batch-inserter/internal/controller/htttp_controller/groups"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/unified/unifaces"
	"context"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

// Shared state of all */event/* routes
var _ RegisterGroup = (*registerGroup)(nil)

type RegisterGroup interface {
	groups.Group
}

func NewRegisterGroup(ctx context.Context, config config.Config, l logger.Logger, topichandler unifaces.Tx[controller.ControllerAPI]) *registerGroup {
	return &registerGroup{
		masternode: topichandler,
		l:          l,
	}
}

type registerGroup struct {
	masternode unifaces.Tx[controller.ControllerAPI]
	l          logger.Logger
}

/* func (g *registerGroup) Attach(root *gin.RouterGroup) {
	_ = g.Chain(root)

} */
/* func (g *registerGroup) Chain(root *gin.RouterGroup) *gin.RouterGroup {
	eGroup := root.Group("event")
	eGroup.POST("", g.makeEvent)
	eGroup.GET(":eid", g.getEvent)
	return eGroup

} */
// ============================================================= TODO MAKE SWAGGER INFO ====================================================
// @Summary      Get event by EID(UUID)
// @Description  TODO : Add description test
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        eid   path      string  true  "Enevent ID"
// @Success      200  {object}  Rsp200
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /event/{eid} [get]
func (group *registerGroup) register(gctx *gin.Context) {
	// reqEID := gctx.Param("eid")
	ctlAPI, txclose, err := group.masternode.Tx()
	if err != nil {
		group.l.Error(err)
		errbytes, _ := sonic.Marshal(&RegisterError{
			Message: err.Error(),
		})
		gctx.Writer.Write(errbytes)
	}
	defer txclose()
	ctlAPI. МАСТЕРНОДА НЕ РЕАЛИЗОВАЛА HANDLE TOPIC 

	/* 	if err != nil {
		gctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	} */
	/* 	gctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"event":   rsp.Event,
	}) */

}

type RegisterError struct {
	Message string `json:"message"`
}
type Rsp200 struct {
	Message string       `json:"message"`
	Event   domain.Event `json:"event"`
}
