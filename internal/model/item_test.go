package model

import (
	"os"
	"path/filepath"
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

type stubDir string

func (d stubDir) ItemsDir() string { return string(d) }

func TestItemSaveFilename(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "items")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	item := Item{ID: 3, Title: "  Hola Mundo!!  ", State: "todo"}
	if err := item.Save(stubDir(dir)); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "003-hola-mundo.toml")
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected %s: %v", want, err)
	}
}