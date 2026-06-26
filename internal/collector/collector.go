package collector

import (
	"encoding/json"
)

type Collector interface {
	Description() Description
	Collect() (any, error)
}

var registry []Collector

func Register(c Collector) {
	registry = append(registry, c)
}

func CollectAll() []Result {
	out := make([]Result, 0, len(registry))
	for _, c := range registry {
		out = append(out, run(c))
	}
	return out
}

func run(c Collector) (r Result) {
	r.Plugin = c.Description()
	data, err := c.Collect() // TODO: concurrent collectors
	if err != nil {
		r.Err = err
		return
	}
	r.Data, r.Err = json.Marshal(data)
	return
}
