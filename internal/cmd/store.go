package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Item represents stored embedding
type Item struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding"`
}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store a text blob and its embedding",
	RunE: func(cmd *cobra.Command, args []string) error {
		text, _ := cmd.Flags().GetString("text")
		if text == "" {
			return cmd.Help()
		}
		storePath, _ := cmd.Flags().GetString("store")
		emb, err := embed(text)
		if err != nil {
			return err
		}
		item := Item{ID: uuid.New().String(), Text: text, Embedding: emb}
		var items []Item
		if _, err := os.Stat(storePath); err == nil {
			b, err := os.ReadFile(storePath)
			if err == nil {
				json.Unmarshal(b, &items)
			}
		}
		items = append(items, item)
		b, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return err
		}
		dir := filepath.Dir(storePath)
		os.MkdirAll(dir, 0755)
		return os.WriteFile(storePath, b, 0644)
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
	storeCmd.Flags().StringP("text", "t", "", "text to store")
}
