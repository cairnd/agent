package users

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/carind/agent/internal/collector"
	"github.com/carind/agent/internal/collector/plugins"
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

const usersPath = "/etc/passwd"
const groupsPath = "/etc/groups"

func collect(c config) (collector.Entries, error) {
	entries := collector.Entries{}

	f, err := os.Open(usersPath)
	if err != nil {
		return entries, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		u, err := parsePasswdLine(line)
		if err != nil {
			// completely skip not even a possible key from a corrupt line
			continue
		}

		e := collector.Entry{
			Key: u.UID,
			Fingerprint: plugins.Fingerprint(
				u.Username,
				u.HomeDir,
				// TODO: check group membership, finger print groups
			),
			Snapshot: u,
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

// USERNAME:PASSWORD:UID:GID:GECOS:HOME_DIR:SHELL
// guest:x:1001:1001:guest,,,:/home/guest:/bin/bash
type passwd struct {
	Username string
	Password string
	UID      string
	GID      string
	Geocos   string
	HomeDir  string
	Shell    string
}

func parsePasswdLine(line string) (passwd, error) {
	parts := strings.Split(line, ":")
	if len(parts) > 7 {
		return passwd{}, fmt.Errorf("malformed passwd line")
	}
	return passwd{
		Username: parts[0],
		Password: parts[1],
		UID:      parts[2],
		GID:      parts[3],
		Geocos:   parts[4],
		HomeDir:  parts[5],
		Shell:    parts[6],
	}, nil
}
