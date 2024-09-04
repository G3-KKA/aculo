package request

import "aculo/frontend-restapi/domain"

type GetEventRequest struct {
	EID string
}
type GetEventResponse struct {
	Event domain.Event
}
