// Package e2e runs the compiled pj binary end-to-end against real .pm/
// directories, simulating a full user workflow (CLI), plus a model-driven
// TUI integration test.
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/chaos/pj/internal/store"
	"github.com/chaos/pj/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	bin         string
	mockEditBin string
)

const mockEditCode = `package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no file provided")
		os.Exit(1)
	}
	filePath := os.Args[len(os.Args)-1]

	if appendText := os.Getenv("MOCK_EDIT_APPEND"); appendText != "" {
		data, _ := os.ReadFile(filePath)
		_ = os.WriteFile(filePath, []byte(string(data)+appendText), 0o644)
		return
	}

	if writeText := os.Getenv("MOCK_EDIT_WRITE"); writeText != "" {
		_ = os.WriteFile(filePath, []byte(writeText), 0o644)
		return
	}

	find := os.Getenv("MOCK_EDIT_FIND")
	replace := os.Getenv("MOCK_EDIT_REPLACE")
	if find != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		newText := strings.ReplaceAll(string(data), find, replace)
		if err := os.WriteFile(filePath, []byte(newText), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
}
`

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "pj-e2e-bin")
	if err != nil {
		panic(err)
	}
	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}
	bin = filepath.Join(dir, "pj"+exe)
	build := exec.Command("go", "build", "-o", bin, "./cmd/pm")
	build.Dir = repoRoot()
	if out, err := build.CombinedOutput(); err != nil {
		panic("build failed: " + string(out))
	}

	mockEditSrc := filepath.Join(dir, "mockedit.go")
	if err := os.WriteFile(mockEditSrc, []byte(mockEditCode), 0o644); err != nil {
		panic(err)
	}
	mockEditBin = filepath.Join(dir, "mockedit"+exe)
	buildMock := exec.Command("go", "build", "-o", mockEditBin, mockEditSrc)
	if out, err := buildMock.CombinedOutput(); err != nil {
		panic("build mockedit failed: " + string(out))
	}

	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func repoRoot() string {
	wd, _ := os.Getwd()
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found")
		}
		dir = parent
	}
}

func run(t *testing.T, dir string, env map[string]string, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "EDITOR="+mockEditBin)
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("pj %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	return stripANSI(string(out))
}

var ansiRe = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func newProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run(t, dir, nil, "init")
	return dir
}

func TestCLIEndToEndFlow(t *testing.T) {
	dir := newProject(t)

	run(t, dir, nil, "add", "First feature", "-d", "does things", "-t", "feat")
	run(t, dir, nil, "add", "Bugfix thing", "-t", "fix")
	run(t, dir, nil, "add", "feat(cli): auto detected feat")
	run(t, dir, nil, "add", "chore: auto detected chore")

	list := run(t, dir, nil, "list")
	for _, want := range []string{
		"[1] [feat] First feature",
		"[2] [fix] Bugfix thing",
		"[3] [feat] feat(cli): auto detected feat",
		"[4] [chore] chore: auto detected chore",
	} {
		if !strings.Contains(list, want) {
			t.Errorf("list missing %q:\n%s", want, list)
		}
	}

	run(t, dir, nil, "done", "1")
	list = run(t, dir, nil, "list")
	if !strings.Contains(list, "First feature — done") {
		t.Errorf("item 1 should be done:\n%s", list)
	}

	status := run(t, dir, nil, "status")
	if !strings.Contains(status, "healthy") {
		t.Errorf("status not healthy:\n%s", status)
	}

	// edit via $EDITOR: change title.
	run(t, dir, map[string]string{"MOCK_EDIT_FIND": "Bugfix", "MOCK_EDIT_REPLACE": "Renamed"}, "edit", "2")
	list = run(t, dir, nil, "list")
	if !strings.Contains(list, "Renamed thing") {
		t.Errorf("edit did not persist new title:\n%s", list)
	}

	run(t, dir, nil, "renum")
	if s := run(t, dir, nil, "status"); !strings.Contains(s, "healthy") {
		t.Errorf("post-renum status broken:\n%s", s)
	}
}

func TestCLIStatusDetectsDuplicateIDs(t *testing.T) {
	dir := newProject(t)
	run(t, dir, nil, "add", "A")
	run(t, dir, nil, "add", "B")
	// duplicate ID 1 manually
	dup := "id = 1\ntitle = 'Dup'\nstate = 'todo'\ncreated = 'c'\nupdated = 'u'\n"
	if err := os.WriteFile(filepath.Join(dir, ".pm", "items", "001-dup.toml"), []byte(dup), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runErr(dir, nil, "status")
	if err == nil {
		t.Fatalf("expected status to fail on duplicate IDs, got: %s", out)
	}
	if !strings.Contains(out, "CONFLICT") || !strings.Contains(out, "renum") {
		t.Errorf("expected conflict report pointing to renum, got:\n%s", out)
	}
}

func TestEditRejectsInvalidStateEndToEnd(t *testing.T) {
	dir := newProject(t)
	run(t, dir, nil, "add", "Guarded")
	before, err := os.ReadFile(filepath.Join(dir, ".pm", "items", "001-guarded.md"))
	if err != nil {
		t.Fatal(err)
	}
	out, err := runErr(dir, map[string]string{"MOCK_EDIT_FIND": "todo", "MOCK_EDIT_REPLACE": "banana"}, "edit", "1")
	if err == nil || !strings.Contains(out, "banana") {
		t.Fatalf("expected invalid-state error mentioning banana, got err=%v out=%s", err, out)
	}
	after, err := os.ReadFile(filepath.Join(dir, ".pm", "items", "001-guarded.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("rejected edit must leave the item file untouched")
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func errOutput(err error) []byte {
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.Stderr
	}
	return nil
}

func runErr(dir string, env map[string]string, args ...string) (string, error) {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "EDITOR="+mockEditBin)
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestTUIStateChangeSavesToDisk(t *testing.T) {
	dir := newProject(t)
	run(t, dir, nil, "add", "TUI target")
	s := store.NewStore(filepath.Join(dir, ".pm"))

	items, err := s.LoadItems()
	if err != nil {
		t.Fatal(err)
	}
	m := tui.NewModel(items, s)

	// Simulate the space key: cycles the selected item's state and saves.
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if _, quit := updated.(tui.Model); !quit { // still running is fine
		itemsAfter, err := s.LoadItems()
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, it := range itemsAfter {
			if it.Title == "TUI target" && it.State == "in progress" {
				found = true
			}
		}
		if !found {
			t.Errorf("space key cycle did not persist 'in progress' to disk: %+v", itemsAfter)
		}
	}
}

func TestAddInteractiveToEndToEnd(t *testing.T) {
	dir := newProject(t)

	// Test 1: pj add "Title" -i with mock editor appending description
	run(t, dir, map[string]string{
		"MOCK_EDIT_APPEND": "\nCustom detailed description\n# comment to ignore\n",
	}, "add", "Item with interactive desc", "-i")

	s := store.NewStore(filepath.Join(dir, ".pm"))
	item, err := s.FindItem(1)
	if err != nil {
		t.Fatalf("failed to find item 1: %v", err)
	}
	if item.Title != "Item with interactive desc" {
		t.Errorf("expected title 'Item with interactive desc', got %q", item.Title)
	}
	if item.Description != "Custom detailed description" {
		t.Errorf("expected description 'Custom detailed description', got %q", item.Description)
	}

	// Test 2: pj add -i with mock editor providing both title and description
	run(t, dir, map[string]string{
		"MOCK_EDIT_WRITE": "feat: full interactive task\n\nBody of full interactive task\n",
	}, "add", "-i")

	item2, err := s.FindItem(2)
	if err != nil {
		t.Fatalf("failed to find item 2: %v", err)
	}
	if item2.Title != "feat: full interactive task" {
		t.Errorf("expected title 'feat: full interactive task', got %q", item2.Title)
	}
	if item2.Type != "feat" {
		t.Errorf("expected auto-detected type 'feat', got %q", item2.Type)
	}
	if item2.Description != "Body of full interactive task" {
		t.Errorf("expected description 'Body of full interactive task', got %q", item2.Description)
	}
}

func TestShellCompletionToEndToEnd(t *testing.T) {
	dir := newProject(t)

	// Verify completion generation works through CLI
	for _, sh := range []string{"bash", "zsh", "fish"} {
		out := run(t, dir, nil, "completion", sh)
		if len(out) == 0 || !strings.Contains(out, "pj") {
			t.Errorf("completion for %s invalid or empty:\n%s", sh, out)
		}
	}
}



