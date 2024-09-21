package metadata

import (
	"aculo/batch-inserter/domain"
	"aculo/batch-inserter/internal/config"
	"aculo/batch-inserter/internal/controller/controller/groups"
	"aculo/batch-inserter/internal/interfaces/txface"
	"aculo/batch-inserter/internal/logger"
	"aculo/batch-inserter/internal/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Shared state of all */event/* routes
var _ MetadataGroup = (*metadataHandler)(nil)

type MetadataGroup interface {
	groups.Group
}

func NewMetadataGroup(ctx context.Context, config config.Config, logger logger.Logger, service txface.Tx[service.ServiceAPI]) *metadataHandler {
	return &metadataHandler{
		service: service,
		logger:  logger,
	}
}

type metadataHandler struct {
	service txface.Tx[service.ServiceAPI]
	logger  logger.Logger
}

func (handler *metadataHandler) Attach(root *gin.RouterGroup) {
	_ = handler.Chain(root)

}
func (handler *metadataHandler) Chain(root *gin.RouterGroup) *gin.RouterGroup {
	eGroup := root.Group("metadata")
	eGroup.GET("", handler.Register)
	return eGroup

}

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
func (handler *metadataHandler) Register(gctx *gin.Context) {
	service, txclose, err := handler.service.Tx()
	if err != nil {
		handler.logger.Error(err.Error())
		regerr := RegisterError{
			Message: err.Error(),
		}
		gctx.JSON(http.StatusInternalServerError, &regerr)
	}
	defer txclose()
	topic, err := service.HandleNewClient(gctx.Request.Context())
	md := metadata{
		Topic: topic,
	}
	gctx.JSON(http.StatusOK, &md)

}

type metadata struct {
	Topic string `json:"topic"`
}

type RegisterError struct {
	Message string `json:"message"`
}
type Rsp200 struct {
	Message string     `json:"message"`
	Event   domain.Log `json:"event"`
}
