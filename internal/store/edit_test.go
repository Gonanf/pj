package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chaos/pj/internal/model"
)

// newEditFixture creates an initialized store with one item whose timestamps
// are pinned to 2020 so timestamp assertions are deterministic.
func newEditFixture(t *testing.T) (*Store, model.Item) {
	t.Helper()
	s := NewStore(filepath.Join(t.TempDir(), ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	it := s.NewItem("Old title", "some description")
	it.Created = "2020-01-01T00:00:00Z"
	it.Updated = "2020-01-01T00:00:00Z"
	if err := it.Save(s); err != nil {
		t.Fatal(err)
	}
	return s, *it
}

func TestEditRenamesFileOnTitleChange(t *testing.T) {
	s, orig := newEditFixture(t)

	edited := orig
	edited.Title = "New title here"
	if err := s.EditItem(orig.ID, &edited); err != nil {
		t.Fatal(err)
	}

	oldPath := filepath.Join(s.ItemsDir(), "001-old-title.md")
	newPath := filepath.Join(s.ItemsDir(), "001-new-title-here.md")
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("old file should be gone after rename: %s", oldPath)
	}
	data, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := s.FindItem(orig.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != orig.ID {
		t.Errorf("ID changed: want %d, got %d", orig.ID, loaded.ID)
	}
	if loaded.Title != "New title here" {
		t.Errorf("title not persisted: got %q", loaded.Title)
	}
	if !strings.Contains(string(data), "new-title") && !strings.Contains(string(data), "New title here") {
		// content sanity: marshaled TOML holds the new title
		t.Errorf("new file does not contain new title:\n%s", data)
	}
}

func TestEditRejectsInvalidState(t *testing.T) {
	s, orig := newEditFixture(t)
	path, err := s.itemFile(1)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	edited := orig
	edited.State = "banana"
	err = s.EditItem(orig.ID, &edited)
	if err == nil {
		t.Fatal("expected error for invalid state, got nil")
	}
	if !strings.Contains(err.Error(), "banana") {
		t.Errorf("error should mention the invalid state, got: %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("original file must remain untouched: %v", err)
	}
	if string(before) != string(after) {
		t.Error("invalid edit must not write anything to disk")
	}
}

func TestEditRejectsEmptyTitle(t *testing.T) {
	s, orig := newEditFixture(t)

	edited := orig
	edited.Title = ""
	err := s.EditItem(orig.ID, &edited)
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestEditFailsOnUnknownID(t *testing.T) {
	s, _ := newEditFixture(t)

	err := s.EditItem(99, &model.Item{Title: "whatever", State: "todo"})
	if err == nil {
		t.Fatal("expected error for unknown ID, got nil")
	}
	if !strings.Contains(err.Error(), "99") {
		t.Errorf("error should mention the ID, got: %v", err)
	}
}

func TestEditUpdatesTimestamp(t *testing.T) {
	s, orig := newEditFixture(t)

	edited := orig
	if err := s.EditItem(orig.ID, &edited); err != nil {
		t.Fatal(err)
	}
	loaded, err := s.FindItem(orig.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Updated == "2020-01-01T00:00:00Z" {
		t.Fatal("updated was not refreshed")
	}
	ts, err := time.Parse(time.RFC3339, loaded.Updated)
	if err != nil {
		t.Fatalf("updated is not RFC3339: %v", err)
	}
	if !ts.After(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("updated should be recent, got %s", loaded.Updated)
	}
	if loaded.Created != "2020-01-01T00:00:00Z" {
		t.Errorf("created must be preserved, got %s", loaded.Created)
	}
}
