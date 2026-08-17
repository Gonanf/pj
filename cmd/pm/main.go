package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/chaos/pj/internal/store"
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
		if len(args) == 1 {
			title = args[0]
		} else {
			return fmt.Errorf("title required (TUI not yet implemented)")
		}
		item := s.NewItem(title, "")
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

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}