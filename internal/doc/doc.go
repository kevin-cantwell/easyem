package doc

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Chunk struct {
	ID   string
	Text string
}

func ExtractText(path string, chunkSize int) ([]Chunk, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(&buf, b); err != nil {
		return nil, err
	}
	words := strings.Fields(buf.String())
	if chunkSize <= 0 {
		return nil, errors.New("chunk size must be positive")
	}
	var chunks []Chunk
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		text := strings.Join(words[i:end], " ")
		chunks = append(chunks, Chunk{Text: text})
	}
	return chunks, nil
}
