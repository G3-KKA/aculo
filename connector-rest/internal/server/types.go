package server

import "github.com/gin-gonic/gin"

type Chain = []Chainable
type AppMux = *gin.Engine

type Group interface {
	Attachable
	Chainable
}
type Attachable interface {
	Attach(*gin.RouterGroup)
}
type Chainable interface {
	Chain(root *gin.RouterGroup) *gin.RouterGroup
}
