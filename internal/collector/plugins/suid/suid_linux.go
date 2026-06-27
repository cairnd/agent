package suid

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/monjiapawne/cairn/internal/collector"
	"github.com/monjiapawne/cairn/internal/collector/plugins"
)

type config struct {
	Roots []string `json:"roots"`
}

func init() {
	collector.Register(
		collector.PluginInfo{
			Name:        "SUID",
			Description: "checks for setuid binaries and scripts",
		},
		config{
			Roots: []string{"/usr/bin"},
		},
		collect,
	)
}

type suidFile struct {
	Perm     string `json:"perm"`
	Contents string `json:"contents"`
	UID      uint32 `json:"UID"`
}

func collect(cfg config) (collector.Entries, error) {
	entries := collector.Entries{}
	for _, root := range cfg.Roots {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			// perm /6000 == setuid or setgid bit set
			if info.Mode()&(os.ModeSetuid|os.ModeSetgid) != 0 {
				uid := uint32(1) // fallback
				if stat, ok := info.Sys().(*syscall.Stat_t); ok {
					uid = stat.Uid
				}

				// Use error in entry
				contents, contentsErr := plugins.HashFile(path)

				sf := suidFile{
					Perm:     info.Mode().String(),
					Contents: contents,
					UID:      uid,
				}
				e := collector.Entry{
					Key: path,
					Fingerprint: collector.Fingerprint(
						sf.Contents,
						sf.Perm,
						strconv.FormatUint(uint64(uid), 10),
					),
					Snapshot: sf,
				}
				if contentsErr != nil {
					e.Error = contentsErr.Error()
				}

				entries = append(entries, e)
			}
			return nil
		})
	}
	return entries, nil
}
