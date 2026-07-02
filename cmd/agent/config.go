package main

import (
	"encoding/json"

	"github.com/carind/agent/internal/collector"
)

func getCollectorConfig() collector.Config {
	// should later read from file
	enabledPlugins := []string{
		"suid",
		//		"file",
		//		"users",
	}
	return collector.NewConfig(map[string]json.RawMessage{}, enabledPlugins)
}
