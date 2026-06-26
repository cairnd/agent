package collector

import "encoding/json"

type Result struct {
	Plugin Description
	Data   json.RawMessage
	Err    error
}
