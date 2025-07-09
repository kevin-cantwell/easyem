package main

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"easyem/internal/config"
	"easyem/internal/doc"
	"easyem/internal/model"
	"easyem/internal/store"
)

var (
	modelName string
	chunkSize int
)

var createCmd = &cobra.Command{
	Use:   "create [files...]",
	Short: "Create document embeddings",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		st, err := store.New()
		if err != nil {
			return err
		}
		for _, file := range args {
			abs, _ := filepath.Abs(file)
			exists, err := st.Exist(abs, modelName, chunkSize)
			if err != nil {
				return err
			}
			if exists {
				fmt.Printf("Embedding already exists for %s\n", file)
				continue
			}
			chunks, err := doc.ExtractText(file, chunkSize)
			if err != nil {
				return err
			}
			var mChunks []model.Chunk
			for _, c := range chunks {
				mChunks = append(mChunks, model.Chunk{ID: uuid.NewString(), Text: c.Text})
			}
			mChunks = model.ComputeEmbeddings(mChunks)
			meta := store.EmbeddingMeta{
				ID:            uuid.NewString(),
				Model:         modelName,
				ChunkSize:     chunkSize,
				DocumentPath:  abs,
				DocumentName:  filepath.Base(file),
				EmbeddingName: fmt.Sprintf("%s_%s_%d.json", filepath.Base(file), modelName, chunkSize),
			}
			data := store.EmbeddingFile{EmbeddingMeta: meta}
			for _, c := range mChunks {
				data.Embeddings = append(data.Embeddings, store.Chunk{ID: c.ID, Text: c.Text, Embedding: c.Vector})
			}
			if err := st.Save(meta.EmbeddingName, meta, data); err != nil {
				return err
			}
			fmt.Printf("Created embeddings for %s\n", file)
		}
		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&modelName, "model", "m", config.DefaultModel, "model name")
	createCmd.Flags().IntVarP(&chunkSize, "chunks", "n", config.DefaultChunkSize, "chunk size")
	rootCmd.AddCommand(createCmd)
}
