// Package server provides an HTTP API for rendering HTML to images.
package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/macawls/ogre"
)

const (
	maxBodySize    = 10 << 20
	maxFontsPerReq = 5
	maxFontSize    = 5 << 20
)

// Server is an HTTP server that renders HTML to images via a JSON API.
type Server struct {
	cache    *Cache
	mux      *http.ServeMux
	addr     string
	fonts    []ogre.FontSource
	renderer *ogre.Renderer
	logger   *slog.Logger
	limiter  *rateLimiter
	metrics  *metrics
	cfg      Config
}

// Config holds settings for the render server.
type Config struct {
	Addr          string
	CacheBytes    int64
	Fonts         []ogre.FontSource
	RateLimit     float64
	RenderTimeout time.Duration
	MaxElements   int
}

type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64
	max     float64
}

type bucket struct {
	tokens   float64
	lastFill time.Time
}

func newRateLimiter(rate float64) *rateLimiter {
	return &rateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		max:     rate * 2,
	}
}

func (rl *rateLimiter) Allow(ip string) bool {
	if rl.rate <= 0 {
		return true
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[ip]
	if !ok {
		b = &bucket{tokens: rl.max, lastFill: now}
		rl.buckets[ip] = b
	}

	elapsed := now.Sub(b.lastFill).Seconds()
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.max {
		b.tokens = rl.max
	}
	b.lastFill = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

type metrics struct {
	mu              sync.Mutex
	renderTotal     int64
	renderErrors    int64
	cacheHits       int64
	cacheMisses     int64
	totalDurationMs int64
}

type fontRequest struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Style  string `json:"style"`
	URL    string `json:"url"`
	Data   string `json:"data"`
}

type renderRequest struct {
	HTML     string         `json:"html"`
	Width    int            `json:"width,omitempty"`
	Height   int            `json:"height,omitempty"`
	Format   string         `json:"format,omitempty"`
	Quality  int            `json:"quality,omitempty"`
	Template string         `json:"template,omitempty"`
	Data     map[string]any `json:"data,omitempty"`
	Fonts    []fontRequest  `json:"fonts,omitempty"`
}

type templateRequest struct {
	Template string         `json:"template"`
	Data     map[string]any `json:"data"`
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Format   string         `json:"format"`
	Quality  int            `json:"quality"`
}

// New creates a Server with the given configuration and registers routes.
func New(cfg Config) *Server {
	if cfg.CacheBytes <= 0 {
		cfg.CacheBytes = 64 << 20
	}
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.RenderTimeout <= 0 {
		cfg.RenderTimeout = 10 * time.Second
	}
	if cfg.MaxElements <= 0 {
		cfg.MaxElements = 1000
	}

	renderer := ogre.NewRenderer()
	for _, f := range cfg.Fonts {
		_ = renderer.LoadFont(f)
	}

	s := &Server{
		cache:    NewCache(cfg.CacheBytes),
		mux:      http.NewServeMux(),
		addr:     cfg.Addr,
		fonts:    cfg.Fonts,
		renderer: renderer,
		logger:   slog.Default(),
		limiter:  newRateLimiter(cfg.RateLimit),
		metrics:  &metrics{},
		cfg:      cfg,
	}

	s.mux.HandleFunc("POST /render", s.withRateLimit(s.handleRender))
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /render/template", s.withRateLimit(s.handleRenderTemplate))
	s.mux.HandleFunc("GET /metrics", s.handleMetrics)

	return s
}

// Handler returns the server's HTTP handler with CORS middleware.
func (s *Server) Handler() http.Handler {
	return corsMiddleware(s.mux)
}

// Start listens on the configured address and serves until interrupted.
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.Handler(),
	}

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(ln)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		_ = sig
		ctx, cancel := context.WithTimeout(context.Background(), 5e9)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-errCh:
		return err
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) withRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr
		}
		if !s.limiter.Allow(ip) {
			httpError(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, `{"status":"ok"}`)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	s.metrics.mu.Lock()
	data := map[string]int64{
		"render_total":      s.metrics.renderTotal,
		"render_errors":     s.metrics.renderErrors,
		"cache_hits":        s.metrics.cacheHits,
		"cache_misses":      s.metrics.cacheMisses,
		"total_duration_ms": s.metrics.totalDurationMs,
	}
	s.metrics.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) handleRender(w http.ResponseWriter, r *http.Request) {
	var req renderRequest
	if err := decodeBody(r, &req); err != nil {
		httpError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.HTML == "" {
		httpError(w, "missing html field", http.StatusBadRequest)
		return
	}

	fonts, err := parseFontRequests(req.Fonts)
	if err != nil {
		httpError(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.renderAndRespond(w, r, req.HTML, req.Width, req.Height, req.Format, req.Quality, fonts)
}

func parseFontRequests(reqs []fontRequest) ([]ogre.FontSource, error) {
	if len(reqs) > maxFontsPerReq {
		return nil, fmt.Errorf("too many fonts: max %d", maxFontsPerReq)
	}
	var fonts []ogre.FontSource
	for _, fr := range reqs {
		src := ogre.FontSource{
			Name:   fr.Name,
			Weight: fr.Weight,
			Style:  fr.Style,
		}
		switch {
		case fr.URL != "":
			src.URL = fr.URL
		case fr.Data != "":
			data, err := base64.StdEncoding.DecodeString(fr.Data)
			if err != nil {
				return nil, fmt.Errorf("invalid base64 font data for %q: %w", fr.Name, err)
			}
			if len(data) > maxFontSize {
				return nil, fmt.Errorf("font %q exceeds %dMB limit", fr.Name, maxFontSize>>20)
			}
			src.Data = data
		default:
			return nil, fmt.Errorf("font %q: must provide url or data", fr.Name)
		}
		fonts = append(fonts, src)
	}
	return fonts, nil
}

func (s *Server) handleRenderTemplate(w http.ResponseWriter, r *http.Request) {
	var req templateRequest
	if err := decodeBody(r, &req); err != nil {
		httpError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Template == "" {
		httpError(w, "missing template field", http.StatusBadRequest)
		return
	}

	tmpl, err := template.New("og").Parse(req.Template)
	if err != nil {
		httpError(w, "invalid template: "+err.Error(), http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, req.Data); err != nil {
		httpError(w, "template execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.renderAndRespond(w, r, buf.String(), req.Width, req.Height, req.Format, req.Quality, nil)
}

func (s *Server) renderAndRespond(w http.ResponseWriter, r *http.Request, html string, width, height int, format string, quality int, fonts []ogre.FontSource) {
	start := time.Now()
	if format == "" {
		format = "svg"
	}

	key := cacheKey(html, width, height, format)

	if data, ok := s.cache.Get(key); ok {
		elapsed := time.Since(start)
		s.metrics.mu.Lock()
		s.metrics.renderTotal++
		s.metrics.cacheHits++
		s.metrics.totalDurationMs += elapsed.Milliseconds()
		s.metrics.mu.Unlock()
		s.logger.Info("render",
			"method", r.Method,
			"path", r.URL.Path,
			"format", format,
			"width", width,
			"height", height,
			"cache_hit", true,
			"duration_ms", elapsed.Milliseconds(),
			"status", http.StatusOK,
			"size_bytes", len(data),
		)
		writeImageResponse(w, data, format, key, true)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.RenderTimeout)
	defer cancel()

	type renderResult struct {
		result *ogre.Result
		err    error
	}
	ch := make(chan renderResult, 1)
	go func() {
		res, err := s.renderer.Render(html, ogre.Options{
			Width:       width,
			Height:      height,
			Format:      ogre.Format(format),
			Quality:     quality,
			Fonts:       fonts,
			MaxElements: s.cfg.MaxElements,
		})
		ch <- renderResult{res, err}
	}()

	select {
	case <-ctx.Done():
		elapsed := time.Since(start)
		s.metrics.mu.Lock()
		s.metrics.renderTotal++
		s.metrics.renderErrors++
		s.metrics.totalDurationMs += elapsed.Milliseconds()
		s.metrics.mu.Unlock()
		s.logger.Error("render failed", "error", "timeout", "html_length", len(html))
		httpError(w, "render timeout exceeded", http.StatusGatewayTimeout)
		return
	case rr := <-ch:
		elapsed := time.Since(start)
		if rr.err != nil {
			s.metrics.mu.Lock()
			s.metrics.renderTotal++
			s.metrics.renderErrors++
			s.metrics.totalDurationMs += elapsed.Milliseconds()
			s.metrics.mu.Unlock()
			s.logger.Error("render failed", "error", rr.err, "html_length", len(html))
			httpError(w, "render failed: "+rr.err.Error(), http.StatusInternalServerError)
			return
		}

		s.cache.Set(key, rr.result.Data)
		s.metrics.mu.Lock()
		s.metrics.renderTotal++
		s.metrics.cacheMisses++
		s.metrics.totalDurationMs += elapsed.Milliseconds()
		s.metrics.mu.Unlock()
		s.logger.Info("render",
			"method", r.Method,
			"path", r.URL.Path,
			"format", format,
			"width", width,
			"height", height,
			"cache_hit", false,
			"duration_ms", elapsed.Milliseconds(),
			"status", http.StatusOK,
			"size_bytes", len(rr.result.Data),
		)
		writeImageResponse(w, rr.result.Data, format, key, false)
	}
}

func cacheKey(html string, width, height int, format string) string {
	h := sha256.New()
	fmt.Fprintf(h, "%s\x00%d\x00%d\x00%s", html, width, height, format)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func writeImageResponse(w http.ResponseWriter, data []byte, format string, etag string, cacheHit bool) {
	switch format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	default:
		w.Header().Set("Content-Type", "image/svg+xml")
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.Header().Set("ETag", `"`+etag+`"`)
	if cacheHit {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("X-Cache", "HIT")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("X-Cache", "MISS")
	}
	_, _ = w.Write(data)
}

func decodeBody(r *http.Request, v any) error {
	body := http.MaxBytesReader(nil, r.Body, maxBodySize)
	defer body.Close()
	dec := json.NewDecoder(body)
	return dec.Decode(v)
}

func httpError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
