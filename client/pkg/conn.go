package aculo

import (
	"io"
)

//go:generate mockery --filename=mock_conn.go --name=Conn --dir=. --structname MockConn  --inpackage=true
type Conn interface {
	io.WriteCloser
}
type ConfigCore struct {
	Dst Destination
}
