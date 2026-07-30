package diy

import (
	_ "embed"
	"encoding/json"
)

//go:embed defaults.json
var defaultsJSON []byte

type defaultsFile struct {
	DefaultItems       map[string]any `json:"defaultItems"`
	DefaultPage        map[string]any `json:"defaultPage"`
	CenterDefaultItems map[string]any `json:"centerDefaultItems"`
}

func loadDefaults() defaultsFile {
	var d defaultsFile
	_ = json.Unmarshal(defaultsJSON, &d)
	if d.DefaultItems == nil {
		d.DefaultItems = map[string]any{}
	}
	if d.DefaultPage == nil {
		d.DefaultPage = map[string]any{}
	}
	return d
}
