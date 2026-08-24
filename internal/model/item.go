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
	Type        string `toml:"type,omitempty"`
	Created     string `toml:"created"`
	Updated     string `toml:"updated"`
	// Body is the free markdown text after the TOML frontmatter. It never
	// round-trips inside the frontmatter itself (hence "-").
	Body string `toml:"-"`
}

// Dir is the minimal surface Save needs from a store (avoids an import cycle).
type Dir interface {
	ItemsDir() string
}

func (i *Item) Save(d Dir) error {
	newPath := filepath.Join(d.ItemsDir(), fmt.Sprintf("%03d-%s.md", i.ID, slug(i.Title)))
	data, err := MarshalItem(i)
	if err != nil {
		return err
	}
	if err := os.WriteFile(newPath, data, 0o644); err != nil {
		return err
	}
	// Transparent migration / rename: drop any other file carrying this ID
	// (legacy .toml or a stale slug). Only runs after the new file is on disk.
	old, err := filepath.Glob(filepath.Join(d.ItemsDir(), fmt.Sprintf("%03d-*", i.ID)))
	if err != nil {
		return err
	}
	for _, p := range old {
		if p != newPath {
			if err := os.Remove(p); err != nil {
				return err
			}
		}
	}
	return nil
}

const fmDelim = "+++"

// MarshalItem renders an Item as hybrid format: TOML frontmatter between +++
// delimiters followed by the free markdown body.
func MarshalItem(i *Item) ([]byte, error) {
	fm, err := toml.Marshal(i)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString(fmDelim + "\n")
	b.Write(fm)
	if !strings.HasSuffix(b.String(), "\n") {
		b.WriteString("\n")
	}
	b.WriteString(fmDelim + "\n")
	if i.Body != "" {
		b.WriteString("\n")
		b.WriteString(i.Body)
	}
	return []byte(b.String()), nil
}

// UnmarshalItem parses either the hybrid format (+++ TOML frontmatter + body)
// or legacy plain-TOML files. Body is empty for legacy items.
func UnmarshalItem(data []byte) (*Item, error) {
	var it Item
	text := string(data)
	if !strings.HasPrefix(text, fmDelim+"\n") && text != fmDelim {
		if err := toml.Unmarshal(data, &it); err != nil {
			return nil, err
		}
		return &it, nil
	}
	rest := strings.TrimPrefix(text, fmDelim+"\n")
	end := strings.Index(rest, "\n"+fmDelim)
	if end < 0 {
		return nil, fmt.Errorf("unterminated +++ frontmatter")
	}
	fm := rest[:end]
	// After "+++": first \n closes the delimiter line, second \n is the
	// canonical blank separator before the body.
	body := strings.TrimPrefix(rest[end+1+len(fmDelim):], "\n")
	body = strings.TrimPrefix(body, "\n")
	if err := toml.Unmarshal([]byte(fm), &it); err != nil {
		return nil, err
	}
	it.Body = body
	return &it, nil
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// ValidStates son todos los estados válidos en orden de ciclo.
var ValidStates = []string{
	"todo",
	"in progress",
	"testing",
	"blocked",
	"done",
	"closed",
	"in specification",
	"discarded",
}

// ValidTypes son los tipos de tarea válidos.
var ValidTypes = []string{"feat", "chore", "fix", "docs"}

// IsValidType verifica si un string es un tipo válido. Vacío es válido
// (backward compatible con items que no tienen tipo).
func IsValidType(s string) bool {
	if s == "" {
		return true
	}
	for _, v := range ValidTypes {
		if v == s {
			return true
		}
	}
	return false
}

// IsValidState verifica si un string es un estado válido.
func IsValidState(s string) bool {
	for _, v := range ValidStates {
		if v == s {
			return true
		}
	}
	return false
}

// NextState devuelve el siguiente estado en el ciclo.
func NextState(current string) string {
	for i, s := range ValidStates {
		if s == current {
			return ValidStates[(i+1)%len(ValidStates)]
		}
	}
	return "todo"
}

// StateIndex devuelve el índice numérico (1-based) de un estado.
func StateIndex(s string) int {
	for i, v := range ValidStates {
		if v == s {
			return i + 1
		}
	}
	return 1
}

func slug(s string) string {
	return strings.Trim(nonAlnum.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
