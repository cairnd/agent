package collector

import (
	"encoding/json"
	"strings"
)

// Config
type Config struct {
	Plugins map[string]json.RawMessage `json:"plugins"`
	// EnabledPlugins list of enabled plugins
	EnabledPlugins map[string]bool `json:"enabled_plugins"`
}

func NewConfig(plugins map[string]json.RawMessage, enabledPlugins []string) Config {
	return Config{
		Plugins:        plugins,
		EnabledPlugins: toSet(enabledPlugins),
	}

}

// CollectAll runs all registers configs, collating into formatted results
func CollectAll(cfg Config) Results {
	var results Results

	// Loop over all registered functions,
	// if enabled run the collect
	for _, p := range registry {
		// Only use enabled plugins if list has entries
		if !cfg.EnabledPlugins[strings.ToLower(p.pluginInfo.Name)] && len(cfg.EnabledPlugins) > 0 {
			continue
		}

		res := Result{Plugin: p.pluginInfo}

		entries, err := p.collect(cfg.Plugins[p.pluginInfo.Name])
		if err != nil {
			res.Error = err
			results = append(results, res)
			continue
		}
		res.Entries = entries

		results = append(results, res)
	}

	return results
}

func toSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[strings.ToLower(v)] = true
	}
	return m
}
