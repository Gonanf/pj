package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chaos/pj/internal/model"
	"github.com/pelletier/go-toml/v2"
)

type Store struct {
	pmDir string
}

func NewStore(pmDir string) *Store {
	return &Store{pmDir: pmDir}
}

// ItemsDir returns the directory where item files live.
// Exposed so model.Item.Save can write without an import cycle.
func (s *Store) ItemsDir() string {
	return filepath.Join(s.pmDir, "items")
}

func (s *Store) Init() error {
	if err := os.MkdirAll(s.ItemsDir(), 0o755); err != nil {
		return err
	}
	project := struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
	}{Name: "project", Version: "0.1.0"}
	data, err := toml.Marshal(project)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.pmDir, "project.toml"), data, 0o644)
}

func (s *Store) NewItem(title, desc string) *model.Item {
	items, err := s.LoadItems()
	if err != nil {
		items = nil
	}
	id := 1
	for _, it := range items {
		if it.ID >= id {
			id = it.ID + 1
		}
	}
	now := time.Now().Format(time.RFC3339)
	return &model.Item{
		ID:          id,
		Title:       title,
		Description: desc,
		State:       "todo",
		Created:     now,
		Updated:     now,
	}
}

func (s *Store) LoadItems() ([]model.Item, error) {
	files, err := filepath.Glob(filepath.Join(s.ItemsDir(), "*.toml"))
	if err != nil {
		return nil, err
	}
	mdFiles, err := filepath.Glob(filepath.Join(s.ItemsDir(), "*.md"))
	if err != nil {
		return nil, err
	}
	files = append(files, mdFiles...)
	items := make([]model.Item, 0, len(files))
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		it, err := model.UnmarshalItem(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		items = append(items, *it)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}

// itemFile returns the path of the file holding the item with the given ID.
// Filenames embed a slug we can't recompute here (model.slug is unexported),
// so files are located by their zero-padded ID prefix.
func (s *Store) itemFile(id int) (string, error) {
	matches, err := filepath.Glob(filepath.Join(s.ItemsDir(), fmt.Sprintf("%03d-*.toml", id)))
	if err != nil {
		return "", err
	}
	mdMatches, err := filepath.Glob(filepath.Join(s.ItemsDir(), fmt.Sprintf("%03d-*.md", id)))
	if err != nil {
		return "", err
	}
	matches = append(matches, mdMatches...)
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return "", fmt.Errorf("item %d not found", id)
	default:
		return "", fmt.Errorf("item %d has %d files, run `pj status`", id, len(matches))
	}
}

// FindItem returns the item with the given ID.
func (s *Store) FindItem(id int) (*model.Item, error) {
	p, err := s.itemFile(id)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	it, err := model.UnmarshalItem(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", p, err)
	}
	return it, nil
}

// DeleteItem removes the file holding the item with the given ID.
func (s *Store) DeleteItem(id int) error {
	p, err := s.itemFile(id)
	if err != nil {
		return err
	}
	return os.Remove(p)
}

// EditItem validates state and title before touching disk, saves edited
// (renaming the file when the title changes), and always refreshes `updated`.
// The old file is removed only after the renamed one was written.
func (s *Store) EditItem(id int, edited *model.Item) error {
	orig, err := s.FindItem(id)
	if err != nil {
		return err
	}
	_ = orig // kept for API symmetry; Save() handles rename cleanup
	if _, err := s.itemFile(id); err != nil {
		return err
	}
	edited.ID = id
	if !model.IsValidState(edited.State) {
		return fmt.Errorf("invalid state %q (valid: %v)", edited.State, model.ValidStates)
	}
	if strings.TrimSpace(edited.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}
	edited.Updated = time.Now().Format(time.RFC3339)
	if err := edited.Save(s); err != nil {
		return err
	}
	// Save() already removed the old-named/legacy file; nothing else to do.
	return nil
}
