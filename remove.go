package main

import (
	"easyem/internal/store"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [files...]",
	Short: "Remove embeddings for documents",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := store.New()
		if err != nil {
			return err
		}
		return st.Delete(args)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
