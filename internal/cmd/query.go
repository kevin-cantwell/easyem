package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query embeddings",
	RunE: func(cmd *cobra.Command, args []string) error {
		text, _ := cmd.Flags().GetString("text")
		if text == "" {
			return cmd.Help()
		}
		storePath, _ := cmd.Flags().GetString("store")
		topk, _ := cmd.Flags().GetInt("topk")
		emb, err := embed(text)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(storePath)
		if err != nil {
			return err
		}
		var items []Item
		if err := json.Unmarshal(b, &items); err != nil {
			return err
		}
		type result struct {
			Item
			Score float32
		}
		var results []result
		for _, it := range items {
			score := cosineSimilarity(emb, it.Embedding)
			results = append(results, result{Item: it, Score: score})
		}
		sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
		if topk > len(results) {
			topk = len(results)
		}
		for i := 0; i < topk; i++ {
			fmt.Printf("%s\t%f\n", results[i].Text, results[i].Score)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringP("text", "t", "", "query text")
	queryCmd.Flags().IntP("topk", "k", 5, "number of results to return")
}

func cosineSimilarity(a, b []float32) float32 {
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
