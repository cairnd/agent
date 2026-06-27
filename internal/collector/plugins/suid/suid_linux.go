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

func init() {
	collector.Register(
		collector.Plugin{
			Name:        "SUID",
			Description: "checks for setuid binaries and scripts",
		},
		Config{
			Roots: []string{"/usr/bin"},
		},
		collect,
	)
}

type Config struct {
	Roots []string `json:"roots"`
}

func collect(cfg Config) (collector.Entries, error) {
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

				entry := suidFile{
					Path:   path,
					Perm:   info.Mode().String(),
					Sha256: plugins.HashFile(path),
					UID:    uid,
				}

				entries = append(entries, collector.Entry{
					Key: path,
					Fingerprint: collector.Fingerprint(
						entry.Sha256,
						entry.Perm,
						strconv.FormatUint(uint64(uid), 10),
					),
					Raw: entry,
				})
			}
			return nil
		})
	}
	return entries, nil
}

type suidFile struct {
	Path   string `json:"path"`
	Perm   string `json:"perm"`
	Sha256 string `json:"sha256"`
	UID    uint32 `json:"UID"`
}
