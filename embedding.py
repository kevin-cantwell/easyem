#!/usr/bin/env python3
import sys, json
from sentence_transformers import SentenceTransformer

model = SentenceTransformer("BAAI/bge-small-en-v1.5")

text = sys.stdin.read().strip()
emb = model.encode([text])[0]
print(json.dumps({"embedding": emb.tolist()}))
