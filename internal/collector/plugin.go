package collector

import (
	"crypto/sha256"
	"encoding/hex"
)

// Plugin is the requirements to register a plugin
type Plugin struct {
	Name        string
	Description string
}

type Entries []Entry

// Entry is one tracked thing a plugin emits per run
type Entry struct {
	// Key is how each entry is tracked
	Key string `json:"key"`
	// Fingerprint is what gets compared to see if there's a diff
	Fingerprint string `json:"fingerprint"`
	// Raw should correlate to changes in the fingerprint and should
	// identify clearly the artifact that has been changed
	Raw any `json:"raw"`
}

func Fingerprint(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))
}
