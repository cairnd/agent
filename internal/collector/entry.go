package collector

// Entries is a collection from a single plugin
type Entries []Entry

// Entry is one tracked thing a plugin emits per run
type Entry struct {
	// Key is how each entry is tracked
	Key string `json:"key"`
	// Fingerprint is what gets compared to see if there's a difference.
	// Typically the hashed version of the snapshot, but up the plugin's discretion what gets
	// included in the diff check.
	Fingerprint string `json:"fingerprint"`
	// Snapshot is a typed object that shows the state and details of the key
	// It should identify why the fingerprint has changed
	Snapshot any `json:"artifact"`
	// Error signifies there was an error processing a key.
	Error string `json:"error,omitzero"`
}
