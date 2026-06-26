package main

import (
	"encoding/json"
	"fmt"

	"github.com/monjiapawne/cairn/internal/collector"
	_ "github.com/monjiapawne/cairn/internal/collector/plugins" // registers plugins via init
)

func main() {
	res := collector.CollectAll()

	// testing print
	for _, r := range res {
		if r.Err != nil {
			fmt.Printf("%s: ERROR %v", r.Plugin.Name, r.Err)
			continue
		}
		pretty, err := json.MarshalIndent(r.Data, "", " ")
		if err != nil {
			fmt.Printf("%s: error bad json %v", r.Plugin.Name, err)
			continue
		}
		fmt.Printf("%s - %s\n%s\n", r.Plugin.Name, r.Plugin.Purpose, pretty)
	}
}
