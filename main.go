package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Record struct {
	ID   int
	Text string
	Vec  []float64
}

type Store struct {
	Records []Record
	nextID  int
}

func (s *Store) Add(text string) error {
	vec := embed(text)
	s.Records = append(s.Records, Record{ID: s.nextID, Text: text, Vec: vec})
	s.nextID++
	return nil
}

func (s *Store) Query(text string, topk int) ([]Record, error) {
	if len(s.Records) == 0 {
		return nil, errors.New("empty store")
	}
	q := embed(text)
	type pair struct {
		rec   Record
		score float64
	}
	pairs := make([]pair, 0, len(s.Records))
	for _, r := range s.Records {
		score := cosine(q, r.Vec)
		pairs = append(pairs, pair{r, score})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].score > pairs[j].score })
	if topk > len(pairs) {
		topk = len(pairs)
	}
	results := make([]Record, topk)
	for i := 0; i < topk; i++ {
		results[i] = pairs[i].rec
	}
	return results, nil
}

func (s *Store) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(s)
}

func LoadStore(path string) (*Store, error) {
	s := &Store{}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Store{}, nil
		}
		return nil, err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(s); err != nil {
		return nil, err
	}
	return s, nil
}

func embed(text string) []float64 {
	const dim = 100
	vec := make([]float64, dim)
	words := strings.Fields(strings.ToLower(text))
	for _, w := range words {
		h := fnv.New32a()
		h.Write([]byte(w))
		idx := int(h.Sum32() % dim)
		vec[idx] += 1
	}
	norm := 0.0
	for _, v := range vec {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i, v := range vec {
			vec[i] = v / norm
		}
	}
	return vec
}

func cosine(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	dot := 0.0
	for i := range a {
		dot += a[i] * b[i]
	}
	if dot < 0 {
		return 0
	}
	return dot
}

func usage() {
	fmt.Println("Usage: easyem <store|query> <text> [--topk K]")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}
	cmd := os.Args[1]
	text := strings.Join(os.Args[2:], " ")

	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".easyem.gob")
	store, err := LoadStore(path)
	if err != nil {
		fmt.Println("error loading store:", err)
		return
	}

	switch cmd {
	case "store":
		if err := store.Add(text); err != nil {
			fmt.Println("error storing text:", err)
			return
		}
		if err := store.Save(path); err != nil {
			fmt.Println("error saving store:", err)
			return
		}
		fmt.Println("stored text with id", store.nextID-1)
	case "query":
		topk := 3
		for i, arg := range os.Args {
			if arg == "--topk" && i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &topk)
			}
		}
		results, err := store.Query(text, topk)
		if err != nil {
			fmt.Println("query error:", err)
			return
		}
		for _, r := range results {
			fmt.Printf("id %d: %s\n", r.ID, r.Text)
		}
	default:
		usage()
	}
}
