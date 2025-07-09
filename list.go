package main

import (
	"easyem/internal/store"
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List cached embeddings",
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := store.New()
		if err != nil {
			return err
		}
		metas, err := st.ReadManifest()
		if err != nil {
			return err
		}
		for _, m := range metas {
			fmt.Printf("%s\t%s\t%d\n", m.DocumentName, m.Model, m.ChunkSize)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
