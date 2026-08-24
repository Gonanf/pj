package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chaos/pj/internal/model"
)

type stubDir string

func (d stubDir) ItemsDir() string { return string(d) }

func sampleItems() []model.Item {
	return []model.Item{
		{ID: 1, Title: "Item 1", State: "done"},
		{ID: 2, Title: "Item 2", State: "todo"},
		{ID: 3, Title: "Item 3", State: "in progress"},
		{ID: 4, Title: "Item 4", State: "blocked"},
	}
}

func TestTUIModel_Init(t *testing.T) {
	m := New(sampleItems())
	cmd := m.Init()
	if cmd != nil {
		t.Errorf("expected Init() to return nil, got %v", cmd)
	}
}

func TestTUIModel_Navigation(t *testing.T) {
	items := sampleItems()
	m := New(items)

	if m.Cursor() != 0 {
		t.Fatalf("expected initial cursor 0, got %d", m.Cursor())
	}

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.Cursor() != 1 {
		t.Errorf("expected cursor 1 after down, got %d", m.Cursor())
	}

	// Move down with 'j'
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if m.Cursor() != 2 {
		t.Errorf("expected cursor 2 after j, got %d", m.Cursor())
	}

	// Move down past end
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.Cursor() != len(items)-1 {
		t.Errorf("expected cursor clamped at %d, got %d", len(items)-1, m.Cursor())
	}

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.Cursor() != len(items)-2 {
		t.Errorf("expected cursor %d after up, got %d", len(items)-2, m.Cursor())
	}

	// Move up with 'k'
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.Cursor() != len(items)-3 {
		t.Errorf("expected cursor %d after k, got %d", len(items)-3, m.Cursor())
	}

	// Move up past start
	for i := 0; i < 5; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
		m = updated.(Model)
	}
	if m.Cursor() != 0 {
		t.Errorf("expected cursor clamped at 0, got %d", m.Cursor())
	}
}

func TestTUIModel_MarkDoneAndSave(t *testing.T) {
	tempDir := filepath.Join(t.TempDir(), "items")
	if err := os.MkdirAll(tempDir, 0o755); err != nil {
		t.Fatal(err)
	}

	dir := stubDir(tempDir)
	items := []model.Item{
		{ID: 1, Title: "First Task", State: "todo"},
	}
	if err := items[0].Save(dir); err != nil {
		t.Fatal(err)
	}

	m := New(items, dir)

	// Press space: todo -> in progress (first step of the cycle)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)

	item := m.Items()[m.Cursor()]
	if item.State != "in progress" {
		t.Errorf("expected state 'in progress' after space, got '%s'", item.State)
	}

	// Keep pressing space until reaching done
	steps := 0
	for item.State != "done" && steps < len(model.ValidStates) {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
		m = updated.(Model)
		item = m.Items()[m.Cursor()]
		steps++
	}
	if item.State != "done" {
		t.Fatalf("expected to reach 'done' cycling, got '%s' after %d steps", item.State, steps)
	}

	// Verify saved file reflects 'done'
	savedPath := filepath.Join(tempDir, "001-first-task.md")
	data, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	if !strings.Contains(string(data), `state = 'done'`) && !strings.Contains(string(data), `state = "done"`) {
		t.Errorf("saved file does not contain state = done: %s", string(data))
	}

	// Verify NextState cycles correctly: len(ValidStates) presses return to the start
	start := "todo"
	s := start
	for i := 0; i < len(model.ValidStates); i++ {
		s = model.NextState(s)
	}
	if s != start {
		t.Errorf("expected NextState to cycle back to %q, got %q", start, s)
	}
}

func TestTUIModel_QuitKeys(t *testing.T) {
	m := New(sampleItems())

	for _, key := range []string{"q", "esc"} {
		var keyMsg tea.KeyMsg
		if key == "q" {
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
		} else {
			keyMsg = tea.KeyMsg{Type: tea.KeyEscape}
		}
		_, cmd := m.Update(keyMsg)
		if cmd == nil {
			t.Errorf("expected quit cmd for key %s, got nil", key)
		}
	}

	// Ctrl+C
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("expected quit cmd for ctrl+c, got nil")
	}
}

func TestTUIModel_View(t *testing.T) {
	items := sampleItems()
	m := New(items)

	view := m.View()

	// Check header & progress
	if !strings.Contains(view, "Project Journal") {
		t.Error("view missing header 'Project Journal'")
	}
	if !strings.Contains(view, "25%") {
		t.Errorf("expected progress bar to show 25%% (1/4 done), got:\n%s", view)
	}

	// Check checkmarks
	if !strings.Contains(view, "■") {
		t.Error("view missing done checkmark '■'")
	}
	if !strings.Contains(view, "□") {
		t.Error("view missing pending checkmark '□'")
	}

	// Check groups
	if !strings.Contains(view, "IN PROGRESS") {
		t.Error("view missing 'IN PROGRESS' group header")
	}
	if !strings.Contains(view, "TODO") {
		t.Error("view missing 'TODO' group header")
	}
	if !strings.Contains(view, "DONE") {
		t.Error("view missing 'DONE' group header")
	}

	// Check footer
	if !strings.Contains(view, "navigate") || !strings.Contains(view, "space cycle") || !strings.Contains(view, "q quit") {
		t.Errorf("view missing expected footer help, got:\n%s", view)
	}
}

func TestTUIModel_Empty(t *testing.T) {
	m := New([]model.Item{})
	view := m.View()
	if !strings.Contains(view, "0%") {
		t.Errorf("expected 0%% for empty items, got:\n%s", view)
	}

	// Navigation should not panic
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)
}

func TestRenderShowsTypeBracket(t *testing.T) {
	items := []model.Item{
		{ID: 1, Title: "Fix crash", State: "todo", Type: "fix"},
	}
	m := New(items)
	view := m.View()
	if !strings.Contains(view, "[fix]") {
		t.Errorf("view missing type bracket '[fix]', got:\n%s", view)
	}
	if !strings.Contains(view, "Fix crash") {
		t.Errorf("view missing item title, got:\n%s", view)
	}
}

func TestRenderWithoutTypeOmitsBracket(t *testing.T) {
	items := []model.Item{
		{ID: 1, Title: "Plain task", State: "todo"},
	}
	m := New(items)
	view := m.View()
	for _, typ := range model.ValidTypes {
		bracket := "[" + typ + "]"
		if strings.Contains(view, bracket) {
			t.Errorf("view should not contain type bracket %q when item has no type, got:\n%s", bracket, view)
		}
	}
}
