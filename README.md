# Intra-Search Go Port

This project is a minimal Go port of [Intra-Search](https://github.com/monish-prabhu/Intra-Search). It allows creating very simple embeddings from PDF documents and querying them via a small HTTP API.

## Features

- Extract text from PDF files and split it into chunks
- Compute bag-of-words embeddings for each chunk
- Store embeddings locally under your cache directory
- Simple HTTP server to query embeddings
- Command line interface similar to the original project

## Commands

```
# create embeddings from PDFs
intrasearch create file1.pdf file2.pdf

# start the HTTP server
intrasearch start

# list cached embeddings
intrasearch list

# remove embeddings for files
intrasearch remove file1.pdf
```

This implementation uses a very naive embedding model and is meant as a demonstration only.

