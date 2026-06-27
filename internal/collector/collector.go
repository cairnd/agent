package collector

import "encoding/json"

// CollectFunc is a plugin: given its raw config section (empty ⇒ defaults), it
// does the collection and returns the entries upstream expects. That's the whole
// contract — plugins stay dumb.
type CollectFunc func(raw json.RawMessage) (Entries, error)

type registration struct {
	plugin  Plugin
	collect CollectFunc
}

var registry []registration

// Register adds a plugin whose collect receives its config decoded into T (from
// defaults). An empty config section keeps the defaults. The raw/JSON boundary
// stays hidden here so plugins only deal with their own typed config.
func Register[T any](p Plugin, defaults T, collect func(cfg T) (Entries, error)) {
	registry = append(registry, registration{
		plugin: p,
		collect: func(raw json.RawMessage) (Entries, error) {
			cfg := defaults
			if len(raw) > 0 {
				if err := json.Unmarshal(raw, &cfg); err != nil {
					return nil, err
				}
			}
			return collect(cfg)
		},
	})
}

// CollectAll runs all registers configs, collating into formatted results
func CollectAll(cfg Config) Results {
	var results Results
	enabled := toSet(cfg.Enabled)
	for _, r := range registry {
		if len(enabled) > 0 && !enabled[r.plugin.Name] {
			continue
		}
		res := Result{Plugin: r.plugin}

		entries, err := r.collect(cfg.Plugins[r.plugin.Name])
		if err != nil {
			res.Err = err
			results = append(results, res)
			continue
		}

		res.Data, res.Err = json.Marshal(entries)
		results = append(results, res)
	}

	return results
}

func toSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}
