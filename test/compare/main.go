package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed static
var staticFS embed.FS

var fixturesDir = filepath.Join("..", "fixtures")

func main() {
	mux := http.NewServeMux()

	sub, _ := fs.Sub(staticFS, "static")
	mux.Handle("/", http.FileServer(http.FS(sub)))

	mux.HandleFunc("/api/fixtures", handleFixtures)
	mux.HandleFunc("/api/fixture-html/", handleFixtureHTML)
	mux.HandleFunc("/api/render/", handleRender)
	mux.HandleFunc("/api/diff/", handleDiff)

	fmt.Println("Compare tool running at http://localhost:4444")
	if err := http.ListenAndServe(":4444", mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func handleFixtures(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var names []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".html") {
			names = append(names, strings.TrimSuffix(e.Name(), ".html"))
		}
	}
	sort.Strings(names)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func handleFixtureHTML(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/fixture-html/")
	data, err := os.ReadFile(filepath.Join(fixturesDir, name+".html"))
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

func handleRender(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/render/")
	if name == "" {
		http.Error(w, "fixture name required", 400)
		return
	}

	html, err := os.ReadFile(filepath.Join(fixturesDir, name+".html"))
	if err != nil {
		http.Error(w, "fixture not found: "+name, 404)
		return
	}

	results := renderAll(string(html))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleDiff(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/diff/")
	if name == "" {
		http.Error(w, "fixture name required", 400)
		return
	}

	html, err := os.ReadFile(filepath.Join(fixturesDir, name+".html"))
	if err != nil {
		http.Error(w, "fixture not found: "+name, 404)
		return
	}

	results := renderAll(string(html))

	type diffPair struct {
		Label string      `json:"label"`
		Diff  *DiffResult `json:"diff,omitempty"`
		Error string      `json:"error,omitempty"`
	}

	var diffs []diffPair

	pairs := [][3]string{
		{"Ogre vs Satori", "ogre", "satori"},
		{"Ogre vs Takumi", "ogre", "takumi"},
		{"Satori vs Takumi", "satori", "takumi"},
	}

	getPNG := func(name string) ([]byte, error) {
		var b64 string
		switch name {
		case "ogre":
			if results.Ogre == nil || results.Ogre.PNG == "" {
				return nil, fmt.Errorf("no ogre PNG")
			}
			b64 = results.Ogre.PNG
		case "satori":
			if results.Satori == nil || results.Satori.PNG == "" {
				return nil, fmt.Errorf("no satori PNG")
			}
			b64 = results.Satori.PNG
		case "takumi":
			if results.Takumi == nil || results.Takumi.PNG == "" {
				return nil, fmt.Errorf("no takumi PNG")
			}
			b64 = results.Takumi.PNG
		}
		return base64.StdEncoding.DecodeString(b64)
	}

	for _, p := range pairs {
		a, errA := getPNG(p[1])
		b, errB := getPNG(p[2])
		if errA != nil || errB != nil {
			msg := ""
			if errA != nil {
				msg = errA.Error()
			} else {
				msg = errB.Error()
			}
			diffs = append(diffs, diffPair{Label: p[0], Error: msg})
			continue
		}
		d, err := pixelDiff(a, b)
		if err != nil {
			diffs = append(diffs, diffPair{Label: p[0], Error: err.Error()})
			continue
		}
		diffs = append(diffs, diffPair{Label: p[0], Diff: d})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diffs)
}
