package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chaos/pj/internal/model"
	"github.com/chaos/pj/internal/store"
	"github.com/chaos/pj/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pj",
	Short: "pj — project journal",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize pj in current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		if err := s.Init(); err != nil {
			return err
		}
		fmt.Println("Initialized pj in", filepath.Join(cwd, ".pm"))
		return nil
	},
}

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new item",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		title := ""
		desc, _ := cmd.Flags().GetString("description")
		if len(args) == 1 {
			title = args[0]
		} else {
			return fmt.Errorf("title required")
		}
		typ, _ := cmd.Flags().GetString("type")
		if typ == "" {
			typ = model.DetectType(title)
		} else if !model.IsValidType(typ) {
			return fmt.Errorf("invalid type %q (valid types: %s)", typ, strings.Join(model.ValidTypes, ", "))
		}
		item := s.NewItem(title, desc)
		item.Type = typ
		if err := item.Save(s); err != nil {
			return err
		}
		fmt.Printf("Added [#%d] %s\n", item.ID, item.Title)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all items",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		items, err := s.LoadItems()
		if err != nil {
			return err
		}
		for _, it := range items {
			stateColor := map[string]string{
				"todo":        "\033[37m",
				"in progress": "\033[33m",
				"done":        "\033[32m",
				"blocked":     "\033[31m",
				"discarded":   "\033[30m",
				"testing":     "\033[35m",
			}
			c := stateColor[it.State]
			if c == "" {
				c = "\033[37m"
			}
			typeStr := ""
			if it.Type != "" {
				typeColor := map[string]string{
					"feat":  "\033[32m",
					"fix":   "\033[31m",
					"chore": "\033[33m",
					"docs":  "\033[34m",
				}
				tc, ok := typeColor[it.Type]
				if !ok {
					tc = "\033[37m"
				}
				typeStr = fmt.Sprintf("%s[%s]\033[0m%s ", tc, it.Type, c)
			}
			fmt.Printf("%s[%d] %s%s — %s\033[0m\n", c, it.ID, typeStr, it.Title, it.State)
		}
		return nil
	},
}

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Mark item as done",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		cwd, _ := os.Getwd()
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		items, _ := s.LoadItems()
		for i := range items {
			if items[i].ID == id {
				items[i].State = "done"
				items[i].Updated = time.Now().Format(time.RFC3339)
				if err := items[i].Save(s); err != nil {
					return err
				}
				fmt.Printf("[#%d] marked as done\n", id)
				return nil
			}
		}
		return fmt.Errorf("item %d not found", id)
	},
}

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Open an item in $EDITOR",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		cwd, _ := os.Getwd()
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		orig, err := s.FindItem(id)
		if err != nil {
			return err
		}

		data, err := model.MarshalItem(orig)
		if err != nil {
			return err
		}
		tmp, err := os.CreateTemp("", "pj-edit-*.md")
		if err != nil {
			return err
		}
		defer os.Remove(tmp.Name())
		if _, err := tmp.Write(data); err != nil {
			tmp.Close()
			return err
		}
		if err := tmp.Close(); err != nil {
			return err
		}

		// ponytail: EDITOR split on spaces; quoted paths with spaces break —
		// switch to shellwords if that ever matters.
		parts := strings.Fields(os.Getenv("EDITOR"))
		if len(parts) == 0 {
			parts = []string{"vi"}
		}
		e := exec.Command(parts[0], parts[1:]...)
		e.Args = append(e.Args, tmp.Name())
		e.Stdin, e.Stdout, e.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := e.Run(); err != nil {
			return fmt.Errorf("editor %q failed: %w", parts[0], err)
		}

		raw, err := os.ReadFile(tmp.Name())
		if err != nil {
			return err
		}
		edited, err := model.UnmarshalItem(raw)
		if err != nil {
			return fmt.Errorf("invalid item after editing: %w", err)
		}
		if err := s.EditItem(id, edited); err != nil {
			return err
		}
		fmt.Printf("[#%d] updated\n", id)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check project health",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		items, _ := s.LoadItems()
		seen := make(map[int]int)
		for _, it := range items {
			seen[it.ID]++
		}
		conflicts := 0
		for id, count := range seen {
			if count > 1 {
				fmt.Printf("CONFLICT: ID %d appears %d times\n", id, count)
				conflicts++
			}
		}
		if conflicts > 0 {
			fmt.Printf("\nRun `pj renum` to fix %d conflict(s)\n", conflicts)
			return fmt.Errorf("conflicts detected")
		}
		fmt.Println("Project healthy")
		return nil
	},
}

var renumCmd = &cobra.Command{
	Use:   "renum",
	Short: "Re-number items by creation time",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		items, _ := s.LoadItems()

		// Delete all existing item files first to avoid orphans
		itemsDir := s.ItemsDir()
		entries, _ := os.ReadDir(itemsDir)
		for _, e := range entries {
			if !e.IsDir() {
				os.Remove(filepath.Join(itemsDir, e.Name()))
			}
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Created < items[j].Created
		})
		for i := range items {
			items[i].ID = i + 1
			items[i].Save(s)
		}
		fmt.Printf("Renumbered %d items\n", len(items))
		return nil
	},
}

var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Print a closing summary of the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		save, _ := cmd.Flags().GetBool("save")
		return runFinish(save)
	},
}

// Contract from spec T3: finish is read-only over .pm and never mutates items.
func runFinish(save bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	pmDir := filepath.Join(cwd, ".pm")
	name, err := projectName(pmDir)
	if err != nil {
		return err
	}
	items, err := store.NewStore(pmDir).LoadItems()
	if err != nil {
		return err
	}
	summary := buildSummary(name, items, time.Now())
	fmt.Print(summary)
	if !save {
		return nil
	}
	return os.WriteFile(filepath.Join(pmDir, "SUMMARY.md"), []byte(summary), 0o644)
}

func projectName(pmDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(pmDir, "project.toml"))
	if err != nil {
		return "", fmt.Errorf("load project: %w", err)
	}
	var p struct {
		Name string `toml:"name"`
	}
	if err := toml.Unmarshal(data, &p); err != nil {
		return "", fmt.Errorf("parse project.toml: %w", err)
	}
	if p.Name == "" {
		p.Name = "unnamed"
	}
	return p.Name, nil
}

var openStates = map[string]bool{
	"todo":             true,
	"in progress":      true,
	"testing":          true,
	"blocked":          true,
	"in specification": true,
}

func buildSummary(name string, items []model.Item, today time.Time) string {
	counts := make(map[string]int)
	var oldest time.Time
	haveOldest := false
	open := []model.Item{}
	for _, it := range items {
		counts[it.State]++
		if ct, err := time.Parse(time.RFC3339, it.Created); err == nil && (!haveOldest || ct.Before(oldest)) {
			oldest, haveOldest = ct, true
		}
		if openStates[it.State] {
			open = append(open, it)
		}
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Project: %s\n\n", name))
	b.WriteString("Items by state:\n")
	for _, st := range model.ValidStates {
		b.WriteString(fmt.Sprintf("  %s: %d\n", st, counts[st]))
	}

	doneish := counts["done"] + counts["closed"]
	active := len(items) - counts["discarded"]
	pct := 0
	if active > 0 {
		pct = doneish * 100 / active
	}
	b.WriteString(fmt.Sprintf("\nCompletion: %d%% (%d of %d active)\n", pct, doneish, active))

	if !haveOldest {
		b.WriteString("\nDuration: n/a\n")
	} else {
		days := int(today.Sub(oldest).Hours() / 24)
		if days < 0 {
			days = 0
		}
		b.WriteString(fmt.Sprintf("\nDuration: %d days\n", days))
	}

	b.WriteString("\nOpen items (not finished):\n")
	if len(open) == 0 {
		b.WriteString("  (none)\n")
	}
	for _, it := range open {
		b.WriteString(fmt.Sprintf("  [%d] %s — %s\n", it.ID, it.Title, it.State))
	}
	return b.String()
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show interactive TUI with progress bar",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		s := store.NewStore(filepath.Join(cwd, ".pm"))
		items, err := s.LoadItems()
		if err != nil {
			return err
		}
		p := tea.NewProgram(tui.New(items, s))
		_, err = p.Run()
		return err
	},
}

func init() {
	addCmd.Flags().StringP("description", "d", "", "description of the item (use single quotes for $ and special chars)")
	addCmd.Flags().StringP("type", "t", "", "item type: feat, chore, fix, docs")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(renumCmd)
	finishCmd.Flags().Bool("save", false, "write the summary to .pm/SUMMARY.md")
	rootCmd.AddCommand(finishCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
