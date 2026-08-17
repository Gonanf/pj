package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectInitAndSave(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".pm", "project.toml")); os.IsNotExist(err) {
		t.Fatal("project.toml not created")
	}
	item := s.NewItem("Test item", "")
	if err := item.Save(s); err != nil {
		t.Fatal(err)
	}
	items, err := s.LoadItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}