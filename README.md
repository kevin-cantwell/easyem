# easyem

`easyem` is a minimal CLI tool for storing text blobs with embeddings and performing semantic search. It uses a local embedding model via a small Python helper script.

## Requirements

- Go 1.21+
- Python 3.8+
- Create a virtual environment and install dependencies:

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install sentence-transformers==2.2.2
```

The Python script downloads the `BAAI/bge-small-en-v1.5` model on first run. It is about ~80MB and works well on a laptop.

## Building

```bash
go build ./cmd/easyem
```

## Usage

### Store text

```bash
./easyem store -t "your text here"
```

Embeddings are stored in `embeddings.json` by default. Use `-s` to change the location.

### Query

```bash
./easyem query -t "search text" -k 3
```

This prints the top results with their cosine similarity.

## External vector DB

For larger datasets you can import the stored embeddings into a vector database such as `faiss` or `pinecone`. The file `embeddings.json` contains raw vectors that can be loaded into any system supporting float32 embeddings.
