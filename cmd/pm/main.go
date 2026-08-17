package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/chaos/pj/internal/store"
	"github.com/chaos/pj/internal/tui"
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
		item := s.NewItem(title, desc)
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
			fmt.Printf("%s[%d] %s — %s\033[0m\n", c, it.ID, it.Title, it.State)
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
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doneCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(renumCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
