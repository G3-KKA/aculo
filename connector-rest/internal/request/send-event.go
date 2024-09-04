package request

type SendEventRequest struct {
	Topic string
	Event []byte
}
type SendEventResponse struct{}
