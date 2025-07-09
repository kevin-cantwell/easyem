package model

import (
	"math"
	"strings"
)

type Vector map[string]float64

type Chunk struct {
	ID     string
	Text   string
	Vector Vector
}

func tokenize(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func embed(text string) Vector {
	vec := make(Vector)
	for _, tok := range tokenize(text) {
		vec[tok]++
	}
	return vec
}

func cosSim(a, b Vector) float64 {
	var dot, normA, normB float64
	for k, av := range a {
		dot += av * b[k]
		normA += av * av
	}
	for _, bv := range b {
		normB += bv * bv
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func ComputeEmbeddings(chunks []Chunk) []Chunk {
	for i := range chunks {
		chunks[i].Vector = embed(chunks[i].Text)
	}
	return chunks
}

func Query(query string, embeddings []Chunk) []struct {
	Chunk
	Score float64
} {
	qvec := embed(query)
	results := []struct {
		Chunk
		Score float64
	}{}
	for _, ch := range embeddings {
		score := cosSim(qvec, ch.Vector)
		if score > 0 {
			results = append(results, struct {
				Chunk
				Score float64
			}{ch, score})
		}
	}
	// sort by score
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	return results
}
