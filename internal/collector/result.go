package collector

import "encoding/json"

type Results []Result

type Result struct {
	Plugin Plugin
	Data   json.RawMessage
	Err    error
}
