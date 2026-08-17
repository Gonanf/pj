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

func init() {
	rootCmd.AddCommand(initCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}