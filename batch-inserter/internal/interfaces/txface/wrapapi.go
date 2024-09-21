package txface

// For any Tx() returnnig struct, not an interface.
//
// Tx()(*obj,txclose,err) => Tx()(ObjAPI,txclose,err)
//
//go:generate mockery --filename=mock_apiwrapper.go --name=ApiWrapper --dir=. --structname MockApiWrapper  --inpackage=true
type ApiWrapper[API any] interface {
	WrapAPI() Tx[API]
}

/* Example
...
func (b *broker) Tx() (*broker, unifaces.TxClose, error)
...

	type wrap struct {
		*broker
	}

	func (w *wrap) Tx()(BrokerAPI, unifaces.TxClose, error) {
		return w.broker.Tx()
	}

	func (b *broker) WrapAPI(unifaces.Tx[*broker]) unifaces.Tx[BrokerAPI] {
		return &wrap{
			broker: b,
		}

	}

*/
