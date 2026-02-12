package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	assetMap        = map[string]string{} // "style.css" → "style-a1b2c3d4.css"
	reverseAssetMap = map[string]string{} // "style-a1b2c3d4.css" → "style.css"
)

func buildAssetMap() error {
	entries, err := os.ReadDir("public")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		data, err := os.ReadFile(filepath.Join("public", name))
		if err != nil {
			return err
		}

		hash := fmt.Sprintf("%x", sha256.Sum256(data))[:12]
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		digested := base + "-" + hash + ext

		assetMap[name] = digested
		reverseAssetMap[digested] = name
	}

	return nil
}

func assetPath(name string) string {
	if digested, ok := assetMap[name]; ok {
		return "/public/" + digested
	}
	return "/public/" + name
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	// Strip "/public/" prefix to get the digested filename
	digested := strings.TrimPrefix(r.URL.Path, "/public/")

	original, ok := reverseAssetMap[digested]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeFile(w, r, filepath.Join("public", original))
}
