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
