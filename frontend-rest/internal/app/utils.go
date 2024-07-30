package app

import (
	"aculo/frontend-restapi/internal/server"
)

func chain(elems ...server.Chainable) server.Chain {
	return elems
}
