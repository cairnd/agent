package main

import (
	"encoding/json"
	"fmt"

	"github.com/carind/agent/internal/collector"
	// Plugin list
	_ "github.com/carind/agent/internal/collector/plugins/file"
	_ "github.com/carind/agent/internal/collector/plugins/suid"
)

func main() {
	res := collector.CollectAll(getCollectorConfig())
	j, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(j))
}
