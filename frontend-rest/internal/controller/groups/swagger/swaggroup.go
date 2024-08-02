package swaggergroup

import (
	"aculo/frontend-restapi/internal/controller/groups"

	_ "aculo/frontend-restapi/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewSwaggerGroup() groups.Group {
	return &swaggerGroup{}
}

type swaggerGroup struct{}

func (g *swaggerGroup) Attach(root *gin.RouterGroup) {
	_ = g.Chain(root)

}
func (g *swaggerGroup) Chain(root *gin.RouterGroup) *gin.RouterGroup {
	swgGroup := root.Group("swagger")
	swgGroup.GET("*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return swgGroup

}
