package passwd

import (
	"strings"
)

const PasswdDefaultPath = "/etc/passwd"

// Passwd is a unix user account
type Passwd struct {
	// USERNAME:PASSWORD:UID:GID:GECOS:HOME_DIR:SHELL
	// guest:x:1001:1001:guest,,,:/home/guest:/bin/bash
	Username string
	Password string // Placeholder can ignore
	UID      string
	GID      string
	Gecos    string
	HomeDir  string
	Shell    string
}

// ParseLine takes a typical unix passwd line and marshals into a typed struct
func ParseLine(line string) (Passwd, bool) {
	if line == "" || strings.HasPrefix(line, "#") {
		return Passwd{}, false
	}
	parts := strings.Split(line, ":")
	if len(parts) != 7 {
		return Passwd{}, false
	}
	return Passwd{
		Username: parts[0],
		Password: parts[1],
		UID:      parts[2],
		GID:      parts[3],
		Gecos:    parts[4],
		HomeDir:  parts[5],
		Shell:    parts[6],
	}, true
}

// helper
