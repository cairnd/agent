package collector

import "encoding/json"

// CollectFunc is a plugin: given its raw config section (empty -> defaults), it
// does the collection and returns the entries upstream expects.
type CollectFunc func(raw json.RawMessage) (Entries, error)

type registration struct {
	pluginInfo PluginInfo
	collect    CollectFunc
}

var registry []registration

// Register adds a plugin to the list of plugins to be used.
//
//   - PluginInfo: metadata about a plugin
//   - Defaults: config passed to the collect function
//   - Collect: the main function that creates entries
func Register[T any](pluginInfo PluginInfo, defaults T, collect func(cfg T) (Entries, error)) {
	registry = append(registry, registration{
		pluginInfo: pluginInfo,
		collect:    addPluginConfig(collect, defaults),
	})
}

// addPluginConfig wraps around a plugin to decode and pass in thw json config, so the collect method
// can stay typed.
func addPluginConfig[T any](collect func(cfg T) (Entries, error), defaults T) CollectFunc {
	return func(raw json.RawMessage) (Entries, error) {
		cfg := defaults
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &cfg); err != nil {
				return nil, err
			}
		}
		return collect(cfg)
	}
}
