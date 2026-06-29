package group

import (
	"bufio"
	"os"
	"strings"
)

const GroupDefaultPath = "/etc/group"

// NAME PW OR PLACEHOLDER GID LIST OF MEMBERS
// admins:x:1004:mj,max
type Group struct {
	Name     string
	Password string // place holder (x) likely ignore
	GID      string
	Members  []string
}

func ParseLine(line string) (Group, bool) {
	if line == "" || strings.HasPrefix(line, "#") {
		return Group{}, false
	}
	parts := strings.Split(line, ":")
	if len(parts) != 4 {
		return Group{}, false
	}

	return Group{
		Name:     parts[0],
		Password: parts[1],
		GID:      parts[2],
		Members:  strings.Split(parts[3], ","),
	}, true
}

func ParseGroupFile(fp string) ([]Group, error) {
	// Groups, enrich users further
	gf, err := os.Open(GroupDefaultPath)
	if err != nil {
		// TODO, we don't techically need this?
		return nil, err
	}

	groups := []Group{}
	scanner := bufio.NewScanner(gf)
	for scanner.Scan() {
		g, ok := ParseLine(scanner.Text())
		if !ok {
			continue
		}
		groups = append(groups, g)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
