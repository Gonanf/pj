package model

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestItemMarshal(t *testing.T) {
	item := Item{
		ID:    1,
		Title: "Planificar",
		State: "todo",
	}
	data, err := toml.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Item
	if err := toml.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ID != item.ID || decoded.Title != item.Title {
		t.Errorf("roundtrip mismatch: got %+v", decoded)
	}
}