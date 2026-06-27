package file

import (
	"os"
	"strconv"

	"github.com/carind/agent/internal/collector"
	"github.com/carind/agent/internal/collector/plugins"
)

type Config struct {
	Files []string `json:"files"`
}

func init() {
	collector.Register(
		collector.PluginInfo{
			Name:        "file",
			Description: "Watch a file's contents, permissions, etc.",
		},
		Config{
			Files: []string{"/etc/passwd", "/etc/shadow"},
		},
		collect,
	)
}

type file struct {
	Perms    string `json:"perms"`
	Contents string `json:"contents"`
}

func collect(c Config) (collector.Entries, error) {
	entries := make(collector.Entries, 0, len(c.Files))

	for _, fp := range c.Files {
		info, err := os.Stat(fp)
		if err != nil {
			entries = append(entries, collector.Entry{Key: fp, Error: err.Error()})
			continue
		}

		perms := strconv.FormatUint(uint64(info.Mode().Perm()), 10)
		contents, contentsErr := plugins.HashFile(fp)
		f := file{
			Perms:    perms,
			Contents: contents,
		}

		e := collector.Entry{
			Key:         fp,
			Fingerprint: plugins.Fingerprint(f.Perms, f.Contents),
			Snapshot:    f,
		}
		if contentsErr != nil {
			e.Error = contentsErr.Error()
		}

		entries = append(entries, e)
	}

	return entries, nil
}
