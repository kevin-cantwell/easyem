package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"easyem/internal/config"
)

type EmbeddingMeta struct {
	ID            string `json:"id"`
	Model         string `json:"model"`
	ChunkSize     int    `json:"chunk_size"`
	DocumentPath  string `json:"document_path"`
	DocumentName  string `json:"document_name"`
	EmbeddingName string `json:"embedding_name"`
}

type Chunk struct {
	ID        string             `json:"id"`
	Text      string             `json:"text"`
	Embedding map[string]float64 `json:"embedding"`
}

type EmbeddingFile struct {
	EmbeddingMeta
	Embeddings []Chunk `json:"embeddings"`
}

type Store struct {
	Dir      string
	Manifest string
}

func New() (*Store, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	dir = filepath.Join(dir, config.AppName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	manifest := filepath.Join(dir, "manifest.json")
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		os.WriteFile(manifest, []byte("[]"), 0644)
	}
	return &Store{Dir: dir, Manifest: manifest}, nil
}

func (s *Store) ReadManifest() ([]EmbeddingMeta, error) {
	data, err := os.ReadFile(s.Manifest)
	if err != nil {
		return nil, err
	}
	var metas []EmbeddingMeta
	if err := json.Unmarshal(data, &metas); err != nil {
		return nil, err
	}
	return metas, nil
}

func (s *Store) appendManifest(meta EmbeddingMeta) error {
	metas, err := s.ReadManifest()
	if err != nil {
		return err
	}
	metas = append(metas, meta)
	data, _ := json.MarshalIndent(metas, "", "  ")
	return os.WriteFile(s.Manifest, data, 0644)
}

func (s *Store) Save(fileName string, meta EmbeddingMeta, data EmbeddingFile) error {
	path := filepath.Join(s.Dir, fileName)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return s.appendManifest(meta)
}

func (s *Store) Load(path string) (*EmbeddingFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data EmbeddingFile
	dec := json.NewDecoder(f)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Store) Exist(filePath, model string, chunkSize int) (bool, error) {
	metas, err := s.ReadManifest()
	if err != nil {
		return false, err
	}
	abs, _ := filepath.Abs(filePath)
	for _, m := range metas {
		if m.DocumentPath == abs && m.Model == model && m.ChunkSize == chunkSize {
			return true, nil
		}
	}
	return false, nil
}

func (s *Store) Delete(files []string) error {
	metas, err := s.ReadManifest()
	if err != nil {
		return err
	}
	newMeta := metas[:0]
	for _, meta := range metas {
		remove := false
		for _, f := range files {
			abs, _ := filepath.Abs(f)
			if meta.DocumentPath == abs {
				os.Remove(filepath.Join(s.Dir, meta.EmbeddingName))
				remove = true
			}
		}
		if !remove {
			newMeta = append(newMeta, meta)
		}
	}
	data, _ := json.MarshalIndent(newMeta, "", "  ")
	return os.WriteFile(s.Manifest, data, 0644)
}

func (s *Store) GetMeta(id string) (*EmbeddingMeta, error) {
	metas, err := s.ReadManifest()
	if err != nil {
		return nil, err
	}
	for _, m := range metas {
		if m.ID == id {
			return &m, nil
		}
	}
	return nil, nil
}
