# easyem

`easyem` is a minimal command line tool for storing text embeddings and searching them.
It relies on scikit-learn's `HashingVectorizer` to generate embeddings locally and
a lightweight SQLite database for storage.

## Installation

Ensure Python 3.12+ is available and install dependencies:

```bash
pip install scikit-learn numpy
```

## Usage

Store a new text snippet:

```bash
python -m easyem store "Some text to remember"
```

Query existing snippets:

```bash
python -m easyem query "search phrase" --topk 3
```
