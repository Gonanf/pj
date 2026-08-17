package model

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Item struct {
	ID          int    `toml:"id"`
	Title       string `toml:"title"`
	Description string `toml:"description"`
	State       string `toml:"state"`
	Created     string `toml:"created"`
	Updated     string `toml:"updated"`
}

// Dir is the minimal surface Save needs from a store (avoids an import cycle).
type Dir interface {
	ItemsDir() string
}

func (i *Item) Save(d Dir) error {
	path := filepath.Join(d.ItemsDir(), fmt.Sprintf("%03d-%s.toml", i.ID, slug(i.Title)))
	data, err := toml.Marshal(i)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func slug(s string) string {
	return strings.Trim(nonAlnum.ReplaceAllString(strings.ToLower(s), "-"), "-")
}