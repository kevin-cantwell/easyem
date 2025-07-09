package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"easyem/internal/config"
	"easyem/internal/model"
	"easyem/internal/store"
)

var port int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := store.New()
		if err != nil {
			return err
		}
		http.HandleFunc("/api/embeddings", func(w http.ResponseWriter, r *http.Request) {
			metas, err := st.ReadManifest()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(metas)
		})
		http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
			if len(parts) < 2 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if parts[0] == "doc" {
				id := parts[1]
				meta, err := st.GetMeta(id)
				if err != nil || meta == nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				http.ServeFile(w, r, meta.DocumentPath)
				return
			}
			id := parts[0]
			if parts[1] != "query" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			query := r.URL.Query().Get("query")
			meta, err := st.GetMeta(id)
			if err != nil || meta == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			data, err := st.Load(filepath.Join(st.Dir, meta.EmbeddingName))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var chunks []model.Chunk
			for _, c := range data.Embeddings {
				chunks = append(chunks, model.Chunk{ID: c.ID, Text: c.Text, Vector: c.Embedding})
			}
			results := model.Query(query, chunks)
			type res struct {
				ID    string  `json:"id"`
				Text  string  `json:"text"`
				Score float64 `json:"similarity"`
			}
			out := []res{}
			for _, r := range results {
				out = append(out, res{ID: r.ID, Text: r.Text, Score: r.Score})
			}
			json.NewEncoder(w).Encode(out)
		})
		fmt.Printf("Server running on port %d\n", port)
		return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	},
}

func init() {
	startCmd.Flags().IntVarP(&port, "port", "p", config.DefaultPort, "server port")
	rootCmd.AddCommand(startCmd)
}
