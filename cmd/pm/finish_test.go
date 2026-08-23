package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chaos/pj/internal/model"
	"github.com/chaos/pj/internal/store"
)

// chdirTemp creates a fresh pj project in a temp dir and moves the process
// there (all commands resolve .pm via the current working directory).
func chdirTemp(t *testing.T) (*store.Store, string) {
	t.Helper()
	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	s := store.NewStore(filepath.Join(dir, ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	return s, dir
}

func writeItem(t *testing.T, s *store.Store, id int, title, state, created string) {
	t.Helper()
	it := model.Item{
		ID: id, Title: title, Description: "",
		State: state, Created: created, Updated: created,
	}
	if err := it.Save(s); err != nil {
		t.Fatal(err)
	}
}

func TestFinishCountsMixedStates(t *testing.T) {
	items := []model.Item{
		{ID: 1, Title: "First", State: "todo", Created: "2026-08-10T10:00:00Z"},
		{ID: 2, Title: "Second", State: "done", Created: "2026-08-11T10:00:00Z"},
		{ID: 3, Title: "Third", State: "discarded", Created: "2026-08-12T10:00:00Z"},
		{ID: 4, Title: "Fourth", State: "in specification", Created: "2026-08-13T10:00:00Z"},
		{ID: 5, Title: "Fifth", State: "todo", Created: "2026-08-14T10:00:00Z"},
	}
	today := time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC)
	got := buildSummary("demo", items, today)

	for _, want := range []string{
		"Project: demo",
		"todo: 2",
		"in specification: 1",
		"done: 1",
		"closed: 0",
		"discarded: 1",
		"Completion: 25% (1 of 4 active)",
		"Duration: 12 days",
		"[1] First — todo",
		"[4] Fourth — in specification",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("summary missing %q\ngot:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"[2] Second", "[3] Third"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("non-open item listed as open: %q\ngot:\n%s", unwanted, got)
		}
	}
}

func TestFinishSaveWritesSummaryMD(t *testing.T) {
	s, dir := chdirTemp(t)
	writeItem(t, s, 1, "Only", "todo", "2026-08-01T00:00:00Z")

	if err := runFinish(true); err != nil {
		t.Fatalf("runFinish(save): %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, ".pm", "SUMMARY.md"))
	if err != nil {
		t.Fatalf("SUMMARY.md not written: %v", err)
	}
	text := string(data)
	for _, want := range []string{"Project: project", "[1] Only — todo", "Completion:", "Duration:"} {
		if !strings.Contains(text, want) {
			t.Errorf("SUMMARY.md missing %q\ngot:\n%s", want, text)
		}
	}

	// Without --save nothing is written.
	s2, dir2 := chdirTemp(t)
	writeItem(t, s2, 1, "Only", "todo", "2026-08-01T00:00:00Z")
	if err := runFinish(false); err != nil {
		t.Fatalf("runFinish(no save): %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir2, ".pm", "SUMMARY.md")); !os.IsNotExist(err) {
		t.Error("SUMMARY.md written without --save")
	}
}

func TestFinishDoesNotModifyItems(t *testing.T) {
	s, _ := chdirTemp(t)
	writeItem(t, s, 1, "Alpha task", "todo", "2026-08-01T00:00:00Z")
	writeItem(t, s, 2, "Beta task", "in progress", "2026-08-02T00:00:00Z")

	snap := func() map[string]string {
		m := map[string]string{}
		entries, err := os.ReadDir(s.ItemsDir())
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			b, err := os.ReadFile(filepath.Join(s.ItemsDir(), e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			m[e.Name()] = string(b)
		}
		return m
	}

	before := snap()
	if err := runFinish(false); err != nil {
		t.Fatalf("runFinish(no save): %v", err)
	}
	if err := runFinish(true); err != nil {
		t.Fatalf("runFinish(save): %v", err)
	}
	after := snap()

	if len(before) != len(after) {
		t.Fatalf("item file set changed: %d -> %d", len(before), len(after))
	}
	for name, want := range before {
		got, ok := after[name]
		if !ok {
			t.Errorf("item file disappeared: %s", name)
			continue
		}
		if want != got {
			t.Errorf("finish modified item file: %s", name)
		}
	}
}

func TestFinishEmptyProject(t *testing.T) {
	chdirTemp(t)
	got := buildSummary("empty", nil, time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC))
	for _, want := range []string{
		"Project: empty",
		"Completion: 0% (0 of 0 active)",
		"Duration: n/a",
		"Open items (not finished):",
		"(none)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\ngot:\n%s", want, got)
		}
	}
	if err := runFinish(false); err != nil {
		t.Fatalf("runFinish on empty project: %v", err)
	}
}
