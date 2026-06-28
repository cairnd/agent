package users

import (
	"bufio"
	"os"

	"github.com/carind/agent/internal/collector"
	"github.com/carind/agent/internal/collector/plugins/internal/linux/group"
	"github.com/carind/agent/internal/collector/plugins/internal/linux/passwd"
)

type config struct {
}

func init() {
	collector.Register(
		collector.PluginInfo{
			Name:        "users",
			Description: "tracks users",
		},
		config{},
		collect,
	)
}

func collect(c config) (collector.Entries, error) {
	uf, err := os.Open(passwd.PasswdDefaultPath)
	if err != nil {
		return nil, err
	}
	defer uf.Close()

	// Users
	users := []passwd.Passwd{}
	scanner := bufio.NewScanner(uf)
	for scanner.Scan() {
		u, ok := passwd.ParseLine(scanner.Text())
		if !ok {
			continue
		}
		users = append(users, u)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	groups, err := group.ParseGroupFile(group.GroupDefaultPath)
	if err != nil {
		// We don't need to quit here..
		return nil, err
	}
	_ = groups

	return nil, nil
}
