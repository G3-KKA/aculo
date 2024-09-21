package streampool

import "errors"

var (
	ErrPoolShuttedDhown   = errors.New("pool shutted down")
	ErrWorkerNotFound     = errors.New("worker not found")
	ErrWorkerAlreadyExist = errors.New("worker already exist ")
)
