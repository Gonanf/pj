package main

import (
	"testing"
)

func TestParseInteractiveDescription(t *testing.T) {
	input := `
# Please enter the description for item: Test
# Lines starting with '#' will be ignored.

This is a multi-line
description of the feature.

- Point 1
- Point 2

# Another comment
`
	expected := "This is a multi-line\ndescription of the feature.\n\n- Point 1\n- Point 2"
	got := parseInteractiveDescription(input)
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}

	// All comments
	allComments := "# Only\n# Comments\n"
	if got := parseInteractiveDescription(allComments); got != "" {
		t.Errorf("expected empty string for comment-only input, got %q", got)
	}

	// Empty input
	if got := parseInteractiveDescription(""); got != "" {
		t.Errorf("expected empty string for empty input, got %q", got)
	}
}

func TestParseInteractiveTitleAndDescription(t *testing.T) {
	input := `
# Please enter item title on first line
feat: my interactive feature

Detailed description here
across multiple lines.
# Ignored footer comment
`
	title, desc := parseInteractiveTitleAndDescription(input)
	if title != "feat: my interactive feature" {
		t.Errorf("got title %q, want %q", title, "feat: my interactive feature")
	}
	expectedDesc := "Detailed description here\nacross multiple lines."
	if desc != expectedDesc {
		t.Errorf("got desc %q, want %q", desc, expectedDesc)
	}

	// Title only
	titleOnly := "# Comment\nSingle line title\n# More comments\n"
	tOnly, dOnly := parseInteractiveTitleAndDescription(titleOnly)
	if tOnly != "Single line title" {
		t.Errorf("got title %q, want %q", tOnly, "Single line title")
	}
	if dOnly != "" {
		t.Errorf("got desc %q, want empty", dOnly)
	}
}
