package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chaos/pj/internal/model"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))

	progressBarFilledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("2")) // Green

	progressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")) // Gray

	groupHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("6")) // Cyan

	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6")) // Cyan arrow

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // Gray

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true)

	stateStyles = map[string]lipgloss.Style{
		"done":        lipgloss.NewStyle().Foreground(lipgloss.Color("2")),  // Verde
		"todo":        lipgloss.NewStyle().Foreground(lipgloss.Color("15")), // Blanco
		"in progress": lipgloss.NewStyle().Foreground(lipgloss.Color("3")),  // Amarillo
		"blocked":     lipgloss.NewStyle().Foreground(lipgloss.Color("1")),  // Rojo
		"testing":     lipgloss.NewStyle().Foreground(lipgloss.Color("5")),  // Magenta
		"discarded":   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),  // Gris
	}
	defaultItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)

type tuiModel struct {
	items  []model.Item
	cursor int
	dir    model.Dir
	width  int
	height int
}

type Model = tuiModel

func statePriority(state string) int {
	switch state {
	case "in progress":
		return 1
	case "todo":
		return 2
	case "testing":
		return 3
	case "blocked":
		return 4
	case "done":
		return 5
	case "discarded":
		return 6
	default:
		return 7
	}
}

func sortItems(items []model.Item) {
	sort.SliceStable(items, func(i, j int) bool {
		pi, pj := statePriority(items[i].State), statePriority(items[j].State)
		if pi != pj {
			return pi < pj
		}
		return items[i].ID < items[j].ID
	})
}

// New initializes a new Bubbletea TUI Model.
func New(items []model.Item, dir ...model.Dir) Model {
	copied := make([]model.Item, len(items))
	copy(copied, items)
	sortItems(copied)

	var d model.Dir
	if len(dir) > 0 {
		d = dir[0]
	}

	return tuiModel{
		items:  copied,
		cursor: 0,
		dir:    d,
	}
}

// NewModel is an alias for New.
func NewModel(items []model.Item, dir ...model.Dir) Model {
	return New(items, dir...)
}

func (m Model) Cursor() int {
	return m.cursor
}

func (m Model) Items() []model.Item {
	return m.items
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "enter":
			if len(m.items) > 0 && m.cursor >= 0 && m.cursor < len(m.items) {
				if m.items[m.cursor].State == "done" {
					m.items[m.cursor].State = "todo"
				} else {
					m.items[m.cursor].State = "done"
				}
				m.items[m.cursor].Updated = time.Now().Format(time.RFC3339)
				if m.dir != nil {
					_ = m.items[m.cursor].Save(m.dir)
				}
				sortItems(m.items)
				if m.cursor >= len(m.items) {
					m.cursor = len(m.items) - 1
				}
				if m.cursor < 0 {
					m.cursor = 0
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("Project Journal") + "\n")

	// Progress bar calculation
	total := len(m.items)
	doneCount := 0
	for _, it := range m.items {
		if it.State == "done" {
			doneCount++
		}
	}

	percent := 0
	if total > 0 {
		percent = (doneCount * 100) / total
	}

	barWidth := 20
	filledWidth := 0
	if total > 0 {
		filledWidth = (doneCount * barWidth) / total
	}
	emptyWidth := barWidth - filledWidth

	filledBar := progressBarFilledStyle.Render(strings.Repeat("█", filledWidth))
	emptyBar := progressBarEmptyStyle.Render(strings.Repeat("░", emptyWidth))

	progressBar := fmt.Sprintf("[%s%s] %d%% (%d/%d done)", filledBar, emptyBar, percent, doneCount, total)
	b.WriteString(progressBar + "\n\n")

	// Item list grouped by state
	if total == 0 {
		b.WriteString(emptyStyle.Render("No items yet. Add an item with: pj add \"Title\"") + "\n\n")
	} else {
		currentGroup := ""
		for i, it := range m.items {
			if it.State != currentGroup {
				if currentGroup != "" {
					b.WriteString("\n")
				}
				currentGroup = it.State
				groupTitle := strings.ToUpper(currentGroup)
				b.WriteString(groupHeaderStyle.Render(groupTitle) + "\n")
			}

			checkmark := "□"
			if it.State == "done" {
				checkmark = "■"
			}

			cursorStr := "  "
			if i == m.cursor {
				cursorStr = cursorStyle.Render("> ")
			}

			st, ok := stateStyles[it.State]
			if !ok {
				st = defaultItemStyle
			}

			itemText := fmt.Sprintf("%s [#%d] %s", checkmark, it.ID, it.Title)
			if i == m.cursor {
				itemText = selectedItemStyle.Render(itemText)
			}
			b.WriteString(cursorStr + st.Render(itemText) + "\n")
		}
		b.WriteString("\n")
	}

	// Footer
	b.WriteString(footerStyle.Render("↑/↓ navigate | space done | q quit"))

	return b.String()
}
