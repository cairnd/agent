package main

import (
	"encoding/json"
	"fmt"

	"github.com/monjiapawne/cairn/internal/collector"
	// Plugin list
	_ "github.com/monjiapawne/cairn/internal/collector/plugins/suid"
)

func main() {
	res := collector.CollectAll(collector.Config{})

	// testing print
	for _, r := range res {
		if r.Err != nil {
			fmt.Printf("%s: ERROR %v", r.Plugin, r.Err)
			continue
		}
		pretty, err := json.MarshalIndent(r.Data, "", " ")
		if err != nil {
			fmt.Printf("%s: error bad json %v", r.Plugin, err)
			continue
		}
		fmt.Printf("%s\n%s\n", r.Plugin, pretty)
	}
}
