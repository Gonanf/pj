package model

import (
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestTypeRoundTrip(t *testing.T) {
	// With type: marshal -> unmarshal preserves it.
	withType := Item{ID: 1, Title: "Fix crash", State: "todo", Type: "fix"}
	data, err := toml.Marshal(withType)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "type") || !strings.Contains(string(data), "fix") {
		t.Errorf("marshaled TOML missing type field:\n%s", data)
	}
	var got Item
	if err := toml.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Type != "fix" {
		t.Errorf("expected Type 'fix' after round-trip, got %q", got.Type)
	}

	// Without type (legacy item): round-trip must not break.
	noType := Item{ID: 2, Title: "Old item", State: "todo"}
	data, err = toml.Marshal(noType)
	if err != nil {
		t.Fatal(err)
	}
	got = Item{}
	if err := toml.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal of legacy item without type failed: %v", err)
	}
	if got.Type != "" {
		t.Errorf("expected empty Type after round-trip, got %q", got.Type)
	}

	// Legacy TOML without a type field unmarshals cleanly.
	got = Item{}
	if err := toml.Unmarshal([]byte("id = 3\ntitle = 'Legacy'\nstate = 'todo'\n"), &got); err != nil {
		t.Fatalf("legacy TOML without type failed to unmarshal: %v", err)
	}
	if got.Type != "" {
		t.Errorf("expected empty Type for legacy TOML, got %q", got.Type)
	}
}

func TestInvalidTypeRejected(t *testing.T) {
	for _, v := range ValidTypes {
		if !IsValidType(v) {
			t.Errorf("IsValidType(%q) = false, want true", v)
		}
	}
	// Empty is valid (backward compatible with items that have no type).
	if !IsValidType("") {
		t.Error(`IsValidType("") = false, want true`)
	}
	for _, bad := range []string{"bug", "feature", "FEAT", "refactor"} {
		if IsValidType(bad) {
			t.Errorf("IsValidType(%q) = true, want false", bad)
		}
	}
}

func TestDetectType(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"feat: add new feature", "feat"},
		{"feat(ui): add button", "feat"},
		{"feat(core/engine): fast parser", "feat"},
		{"feat!: breaking change", "feat"},
		{"feat(api)!: breaking api change", "feat"},
		{"fix: fix crash", "fix"},
		{"fix(store): handle empty dir", "fix"},
		{"chore: update dependencies", "chore"},
		{"chore(deps): update go-toml", "chore"},
		{"docs: update readme", "docs"},
		{"docs(setup): install guide", "docs"},
		{"FEAT: uppercase feat", "feat"},
		{"Fix(UI): mixed case fix", "fix"},
		{"Regular title without prefix", ""},
		{"feat without colon", ""},
		{"unknown: title with invalid type", ""},
		{"refactor(core): unsupported type", ""},
		{"  feat: leading spaces  ", "feat"},
	}

	for _, tc := range tests {
		got := DetectType(tc.title)
		if got != tc.expected {
			t.Errorf("DetectType(%q) = %q, want %q", tc.title, got, tc.expected)
		}
	}
}

