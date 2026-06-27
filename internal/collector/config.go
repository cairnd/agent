package collector

import "encoding/json"

// Config maps to each plugin by name, it's decoded by each plugin
// decoded by the plugin that owns its schema, keyed by Plugin.Name.
type Config struct {
	Plugins map[string]json.RawMessage `json:"plugins"`
	Enabled []string
}
