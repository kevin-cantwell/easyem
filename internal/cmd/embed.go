package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func embed(text string) ([]float32, error) {
	exe, _ := os.Executable()
	script := filepath.Join(filepath.Dir(exe), "embedding.py")
	cmd := exec.Command("python3", script)
	cmd.Stdin = strings.NewReader(text)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var res struct {
		Embedding []float32 `json:"embedding"`
	}
	err = json.Unmarshal(out, &res)
	return res.Embedding, err
}
