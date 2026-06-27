package file

import "github.com/monjiapawne/cairn/internal/collector"

type Config struct {
	Files []string `json:"files"`
}

func init() {
	collector.Register(
		collector.Plugin{
			Name:        "file",
			Description: "Watch a file's contents, permissions, etc.",
		},
		Config{
		}
	)

}
