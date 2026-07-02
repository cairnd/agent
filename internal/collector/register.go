package collector

import "encoding/json"

// CollectFunc is a plugin: given its raw config section (empty -> defaults), it
// does the collection and returns the entries upstream expects.
type CollectFunc func(raw json.RawMessage) (Entries, error)

type registration struct {
	pluginInfo PluginInfo
	collect    CollectFunc
	spec       Spec
}

var registry []registration

// Register adds a plugin to the list of plugins to be used.
//
//   - PluginInfo: metadata about a plugin
//   - Defaults: config passed to the collect function
//   - Collect: the main function that creates entries
//   - Spec(optional): defines the mapping for the stream server to state changes in the snapshots
func Register[T any](pluginInfo PluginInfo, defaults T, collect func(cfg T) (Entries, error), spec ...Spec) {
	r := registration{
		pluginInfo: pluginInfo,
		collect:    addPluginConfig(collect, defaults),
	}
	if len(spec) > 0 {
		r.spec = spec[0]
	}
	registry = append(registry, r)
}

type Manifest struct {
	Plugin PluginInfo `json:"plugin"`
	Spec   Spec       `json:"spec,omitempty"`
}

func Manifests() []Manifest {
	ms := make([]Manifest, 0, len(registry))
	for _, r := range registry {
		ms = append(ms, Manifest{Plugin: r.pluginInfo, Spec: r.spec})
	}
	return ms
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
