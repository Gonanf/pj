package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chaos/pj/internal/model"
)

func TestLoadItemsMixedFormats(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	legacy := "id = 1\ntitle = 'Legacy'\ndescription = ''\nstate = 'todo'\ncreated = 'c'\nupdated = 'u'\n"
	if err := os.WriteFile(filepath.Join(s.ItemsDir(), "001-legacy.toml"), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	hybrid := "+++\nid = 2\ntitle = 'Hybrid'\ndescription = ''\nstate = 'in progress'\ncreated = 'c'\nupdated = 'u'\n+++\n\nFree markdown body.\n"
	if err := os.WriteFile(filepath.Join(s.ItemsDir(), "002-hybrid.md"), []byte(hybrid), 0o644); err != nil {
		t.Fatal(err)
	}
	items, err := s.LoadItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d: %+v", len(items), items)
	}
	if items[0].Title != "Legacy" || items[1].Title != "Hybrid" {
		t.Errorf("bad load: %+v", items)
	}
	if items[1].Body != "Free markdown body.\n" {
		t.Errorf("hybrid body lost: %q", items[1].Body)
	}
}

func TestFindItemReadsHybridFile(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	hybrid := "+++\nid = 4\ntitle = 'Four'\ndescription = ''\nstate = 'done'\ncreated = 'c'\nupdated = 'u'\n+++\n\nbody\n"
	if err := os.WriteFile(filepath.Join(s.ItemsDir(), "004-four.md"), []byte(hybrid), 0o644); err != nil {
		t.Fatal(err)
	}
	it, err := s.FindItem(4)
	if err != nil {
		t.Fatal(err)
	}
	if it.Title != "Four" || it.Body != "body\n" {
		t.Errorf("bad find: %+v", it)
	}
}

func TestSaveMigratesLegacyTOMLToHybrid(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	legacy := "id = 1\ntitle = 'Old'\ndescription = ''\nstate = 'todo'\ncreated = 'c'\nupdated = 'u'\n"
	oldPath := filepath.Join(s.ItemsDir(), "001-old.toml")
	if err := os.WriteFile(oldPath, []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}

	it, err := s.FindItem(1)
	if err != nil {
		t.Fatal(err)
	}
	it.State = "done"
	if err := it.Save(s); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Errorf("legacy .toml should be removed after migration, stat err=%v", err)
	}
	newPath := filepath.Join(s.ItemsDir(), "001-old.md")
	data, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "+++\n") || !strings.Contains(string(data), "state = 'done'") {
		t.Errorf("expected hybrid format with new state:\n%s", data)
	}

	// And the migrated project still loads cleanly.
	items, err := s.LoadItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].State != "done" {
		t.Errorf("post-migration load broken: %+v", items)
	}
}

func TestEditItemPreservesBody(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	orig := &model.Item{ID: 1, Title: "With body", State: "todo", Created: "c", Updated: "u", Body: "keep me"}
	if err := orig.Save(s); err != nil {
		t.Fatal(err)
	}
	loaded, err := s.FindItem(1)
	if err != nil {
		t.Fatal(err)
	}
	loaded.State = "testing"
	if err := s.EditItem(1, loaded); err != nil {
		t.Fatal(err)
	}
	got, err := s.FindItem(1)
	if err != nil {
		t.Fatal(err)
	}
	if got.Body != "keep me" {
		t.Errorf("body lost during edit: %q", got.Body)
	}
}
