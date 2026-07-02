package main

import (
	"encoding/json"
	"fmt"

	"github.com/carind/agent/internal/collector"
	// Plugin list
	_ "github.com/carind/agent/internal/collector/plugins/file"
	_ "github.com/carind/agent/internal/collector/plugins/suid"
	_ "github.com/carind/agent/internal/collector/plugins/users"
)

func main() {
	manifests := collector.Manifests()
	j, _ := json.MarshalIndent(manifests, "", " ")
	fmt.Println(string(j))

	res := collector.CollectAll(getCollectorConfig())
	j, _ = json.MarshalIndent(res, "", " ")
	fmt.Println(string(j))
}
