package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chaos/pj/internal/model"
	"github.com/chaos/pj/internal/store"
	"github.com/spf13/cobra"
)

func TestShellCompletionGenerators(t *testing.T) {
	shells := []string{"bash", "zsh", "fish"}
	for _, sh := range shells {
		var buf bytes.Buffer
		cmd := rootCmd
		cmd.SetOut(&buf)
		cmd.SetArgs([]string{"completion", sh})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("completion %s failed: %v", sh, err)
		}
		out := buf.String()
		if len(out) == 0 {
			t.Errorf("expected completion script for %s to be non-empty", sh)
		}
		if !strings.Contains(out, "pj") {
			t.Errorf("expected completion script for %s to contain 'pj', got %q", sh, out[:min(len(out), 100)])
		}
	}
}

func TestItemIDCompletion(t *testing.T) {
	// 1. Without .pm directory (fail-safe)
	emptyDir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	_ = os.Chdir(emptyDir)
	defer os.Chdir(old)

	comps, directive := itemIDCompletion(doneCmd, []string{}, "")
	if len(comps) != 0 {
		t.Errorf("expected 0 completions in empty directory, got %v", comps)
	}
	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
	}

	// 2. With items in .pm
	s := store.NewStore(filepath.Join(emptyDir, ".pm"))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	it1 := model.Item{ID: 1, Title: "First task", State: "todo"}
	it2 := model.Item{ID: 2, Title: "Second task", State: "in progress"}
	_ = it1.Save(s)
	_ = it2.Save(s)

	comps, directive = itemIDCompletion(doneCmd, []string{}, "")
	if len(comps) != 2 {
		t.Fatalf("expected 2 completions, got %d: %v", len(comps), comps)
	}
	if !strings.HasPrefix(comps[0], "1\t") || !strings.Contains(comps[0], "First task") {
		t.Errorf("expected '1\tFirst task', got %q", comps[0])
	}
	if !strings.HasPrefix(comps[1], "2\t") || !strings.Contains(comps[1], "Second task") {
		t.Errorf("expected '2\tSecond task', got %q", comps[1])
	}

	// 3. When an ID argument has already been provided, don't complete further args
	compsAfter, _ := itemIDCompletion(doneCmd, []string{"1"}, "")
	if len(compsAfter) != 0 {
		t.Errorf("expected 0 completions when args already provided, got %v", compsAfter)
	}
}
