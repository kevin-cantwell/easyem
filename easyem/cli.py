import argparse
import sqlite3
import pickle
from typing import List, Tuple

import numpy as np
from sklearn.feature_extraction.text import HashingVectorizer

DB_PATH = "embeddings.db"


class EasyEm:
    def __init__(self, db_path: str = DB_PATH, n_features: int = 512):
        self.db_path = db_path
        self.vectorizer = HashingVectorizer(n_features=n_features, alternate_sign=False, norm=None)
        self._ensure_db()

    def _ensure_db(self) -> None:
        conn = sqlite3.connect(self.db_path)
        conn.execute(
            """
            CREATE TABLE IF NOT EXISTS embeddings (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                text TEXT NOT NULL,
                embedding BLOB NOT NULL
            )
            """
        )
        conn.commit()
        conn.close()

    def _embed(self, text: str) -> np.ndarray:
        vec = self.vectorizer.transform([text]).toarray()[0]
        return vec.astype(np.float32)

    def store(self, text: str) -> None:
        embedding = self._embed(text)
        conn = sqlite3.connect(self.db_path)
        conn.execute(
            "INSERT INTO embeddings (text, embedding) VALUES (?, ?)",
            (text, pickle.dumps(embedding)),
        )
        conn.commit()
        conn.close()

    def query(self, query: str, topk: int = 5) -> List[Tuple[float, int, str]]:
        q_vec = self._embed(query)
        q_norm = np.linalg.norm(q_vec)
        conn = sqlite3.connect(self.db_path)
        rows = conn.execute("SELECT id, text, embedding FROM embeddings").fetchall()
        conn.close()

        results = []
        for row in rows:
            emb = pickle.loads(row[2])
            score = float(np.dot(q_vec, emb) / (q_norm * np.linalg.norm(emb)))
            results.append((score, row[0], row[1]))
        results.sort(key=lambda x: x[0], reverse=True)
        return results[:topk]


def cli(argv: List[str] | None = None) -> None:
    parser = argparse.ArgumentParser(prog="easyem", description="Minimal embedding CLI")
    subparsers = parser.add_subparsers(dest="command", required=True)

    store_p = subparsers.add_parser("store", help="Store a text blob")
    store_p.add_argument("text", help="Text to embed and store")

    query_p = subparsers.add_parser("query", help="Query stored texts")
    query_p.add_argument("text", help="Search query")
    query_p.add_argument("--topk", type=int, default=5, help="Number of results")

    args = parser.parse_args(argv)
    app = EasyEm()

    if args.command == "store":
        app.store(args.text)
        print("stored")
    elif args.command == "query":
        results = app.query(args.text, args.topk)
        for score, idx, text in results:
            print(f"{score:.3f}\t{idx}\t{text}")


if __name__ == "__main__":
    cli()
