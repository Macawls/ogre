package ogre

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
)

type HandlerConfig struct {
	Width   int
	Height  int
	Format  Format
	Quality int
}

type renderPayload struct {
	HTML     string         `json:"html"`
	Template string         `json:"template"`
	Data     map[string]any `json:"data"`
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Format   string         `json:"format"`
	Quality  int            `json:"quality"`
}

func (r *Renderer) Handler(cfg HandlerConfig) http.Handler {
	if cfg.Width <= 0 {
		cfg.Width = 1200
	}
	if cfg.Height <= 0 {
		cfg.Height = 630
	}
	if cfg.Format == "" {
		cfg.Format = FormatPNG
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var p renderPayload
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			jsonError(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		html := p.HTML
		if html == "" && p.Template != "" {
			tmpl, err := template.New("og").Parse(p.Template)
			if err != nil {
				jsonError(w, "invalid template: "+err.Error(), http.StatusBadRequest)
				return
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, p.Data); err != nil {
				jsonError(w, "template execution failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			html = buf.String()
		}

		if html == "" {
			jsonError(w, "missing html or template field", http.StatusBadRequest)
			return
		}

		width := p.Width
		if width <= 0 {
			width = cfg.Width
		}
		height := p.Height
		if height <= 0 {
			height = cfg.Height
		}
		format := Format(p.Format)
		if format == "" {
			format = cfg.Format
		}
		quality := p.Quality
		if quality <= 0 {
			quality = cfg.Quality
		}

		result, err := r.Render(html, Options{
			Width:   width,
			Height:  height,
			Format:  format,
			Quality: quality,
		})
		if err != nil {
			jsonError(w, "render failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", result.ContentType)
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Write(result.Data)
	})
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
