package controller

import "errors"

var ErrConfigSameAddresses = errors.New("grpc and http server's addresses are the same")
