package model

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestMarshalHybrid(t *testing.T) {
	it := Item{ID: 1, Title: "Hybrid", Description: "d", State: "todo", Body: "# Hello\n\nfree text"}
	data, err := MarshalItem(&it)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.HasPrefix(s, "+++\n") {
		t.Errorf("expected +++ frontmatter delimiter, got:\n%s", s)
	}
	if !strings.Contains(s, "title = 'Hybrid'") && !strings.Contains(s, "title = \"Hybrid\"") {
		t.Errorf("frontmatter missing title:\n%s", s)
	}
	if !strings.Contains(s, "# Hello") {
		t.Errorf("body missing after frontmatter:\n%s", s)
	}
	if strings.Count(s, "+++") != 2 {
		t.Errorf("expected exactly two +++ delimiters:\n%s", s)
	}
}

func TestUnmarshalHybrid(t *testing.T) {
	raw := "+++\nid = 7\ntitle = 'T'\ndescription = ''\nstate = 'todo'\ntype = 'feat'\ncreated = '2026-01-01T00:00:00Z'\nupdated = '2026-01-01T00:00:00Z'\n+++\n\nSome **markdown** body.\n"
	it, err := UnmarshalItem([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if it.ID != 7 || it.Title != "T" || it.Type != "feat" {
		t.Errorf("bad parse: %+v", it)
	}
	if it.Body != "Some **markdown** body.\n" {
		t.Errorf("bad body: %q", it.Body)
	}
}

func TestUnmarshalLegacyTOML(t *testing.T) {
	raw := "id = 3\ntitle = 'Legacy'\nstate = 'todo'\ncreated = 'x'\nupdated = 'y'\n"
	it, err := UnmarshalItem([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if it.ID != 3 || it.Title != "Legacy" || it.Body != "" {
		t.Errorf("legacy parse broken: %+v", it)
	}
}

func TestHybridRoundTripPreservesBody(t *testing.T) {
	it := Item{ID: 2, Title: "RT", State: "done", Body: "line one\nline two\n"}
	data, err := MarshalItem(&it)
	if err != nil {
		t.Fatal(err)
	}
	got, err := UnmarshalItem(data)
	if err != nil {
		t.Fatal(err)
	}
	if got.Body != it.Body || got.Title != it.Title || got.State != it.State {
		t.Errorf("round trip lost data: %+v vs %+v", got, &it)
	}
}

func TestSaveWritesHybridMD(t *testing.T) {
	dir := t.TempDir()
	st := stubDir(dir)
	it := Item{ID: 5, Title: "Saved", State: "todo", Body: "body here"}
	if err := it.Save(st); err != nil {
		t.Fatal(err)
	}
	entries, _ := filepath.Glob(dir + "/005-*")
	if len(entries) != 1 {
		t.Fatalf("expected exactly one file, got %v", entries)
	}
	if !strings.HasSuffix(entries[0], ".md") {
		t.Errorf("expected .md file, got %s", entries[0])
	}
}
