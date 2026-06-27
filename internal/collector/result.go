package collector

type Results []Result

type Result struct {
	Plugin  PluginInfo `json:"plugin"`
	Entries Entries    `json:"entries"`
	Error   error      `json:"error,omitempty,omitzero"`
}
