package main

import (
	"encoding/json"
	"fmt"

	"github.com/monjiapawne/cairn/internal/collector"
	// Plugin list
	_ "github.com/monjiapawne/cairn/internal/collector/plugins/file"
	_ "github.com/monjiapawne/cairn/internal/collector/plugins/suid"
)

func main() {
	colCfg := collector.NewConfig(map[string]json.RawMessage{}, []string{"file", "suid"})
	res := collector.CollectAll(colCfg)

	j, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(j))
}
