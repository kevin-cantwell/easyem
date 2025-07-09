# easyem

A minimal CLI tool for storing text embeddings and querying them semantically.

## Usage

Build the binary:

```
go build
```

Store a piece of text:

```
./easyem store "some text to remember"
```

Query the store:

```
./easyem query "search text" --topk 5
```

Embeddings are saved to `~/.easyem.gob`.
