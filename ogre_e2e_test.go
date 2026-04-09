package ogre

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestGoldenFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/fixtures/*.html")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no fixture files found")
	}

	for _, f := range files {
		name := filepath.Base(f)
		t.Run(name, func(t *testing.T) {
			html, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}

			result, err := Render(string(html), Options{Width: 1200, Height: 630})
			if err != nil {
				t.Fatal(err)
			}

			goldenPath := filepath.Join("testdata", "golden", strings.TrimSuffix(name, ".html")+".svg")

			if *update {
				if err := os.WriteFile(goldenPath, result.Data, 0644); err != nil {
					t.Fatal(err)
				}
				return
			}

			expected, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatal("golden file missing, run with -update")
			}

			if string(result.Data) != string(expected) {
				t.Errorf("output differs from golden file %s", goldenPath)
			}
		})
	}
}
