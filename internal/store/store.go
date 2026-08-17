package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	items := make([]model.Item, 0, len(files))
	for _, f := range files {
		var it model.Item
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if err := toml.Unmarshal(data, &it); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		items = append(items, it)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}