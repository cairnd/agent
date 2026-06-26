package plugins

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/monjiapawne/cairn/internal/collector"
)

type SUID struct {
	Roots []string
}

// Registers name and function to collect
func init() { collector.Register(SUID{Roots: []string{"/usr/bin"}}) }

func (SUID) Description() collector.Description {
	return collector.Description{
		Name:    "SUID",
		Purpose: "Scan scripts/binaries with setuid permission set",
	}
}

func (c SUID) Collect() (any, error) {
	var files []SUIDFile

	for _, root := range c.Roots {
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
				stat, ok := info.Sys().(*syscall.Stat_t)
				uid := stat.Uid
				if !ok {
					// for now just set to 1 but this needs to have error collection logic
					uid = 1
				}

				files = append(files, SUIDFile{
					Path:   path,
					Perm:   info.Mode().String(),
					Sha256: hashFile(path),
					UID:    uid,
				})
			}
			return nil
		})
	}
	return files, nil
}

type SUIDFile struct {
	Path   string `json:"path"`
	Perm   string `json:"perm"`
	Sha256 string `json:"sha256"`
	UID    uint32 `json:"UID"`
}
