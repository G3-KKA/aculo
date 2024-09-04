package swaggergroup

import (
	"aculo/frontend-restapi/internal/controller/groups"
	// In the future we might want to change docs in code,
	// but now we just need to import it, for swagger to work
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
