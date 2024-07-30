package app

import (
	"aculo/connector-restapi/internal/server"
)

func chain(elems ...server.Chainable) server.Chain {
	return elems
}
