package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/macawls/ogre"
)

type RenderResult struct {
	SVG   string  `json:"svg,omitempty"`
	PNG   string  `json:"png,omitempty"`
	Ms    float64 `json:"ms"`
	Error string  `json:"error,omitempty"`
}

type AllResults struct {
	Ogre   *RenderResult `json:"ogre"`
	Satori *RenderResult `json:"satori"`
	Takumi *RenderResult `json:"takumi"`
}

var ogreRenderer = ogre.NewRenderer()

func renderOgre(html string) *RenderResult {
	start := time.Now()
	svgResult, err := ogreRenderer.Render(html, ogre.Options{Width: 1200, Height: 630, Format: ogre.FormatSVG})
	if err != nil {
		return &RenderResult{Error: err.Error()}
	}
	pngResult, err := ogreRenderer.Render(html, ogre.Options{Width: 1200, Height: 630, Format: ogre.FormatPNG})
	if err != nil {
		return &RenderResult{Error: err.Error()}
	}
	ms := float64(time.Since(start).Microseconds()) / 1000

	return &RenderResult{
		SVG: base64.StdEncoding.EncodeToString(svgResult.Data),
		PNG: base64.StdEncoding.EncodeToString(pngResult.Data),
		Ms:  ms,
	}
}

type jsResult struct {
	Satori *RenderResult `json:"satori"`
	Takumi *RenderResult `json:"takumi"`
}

func renderSatoriTakumi(html string) (*RenderResult, *RenderResult) {
	scriptPath := filepath.Join("..", "satori-reference", "render-one.ts")
	cmd := exec.Command("bun", "run", scriptPath)
	cmd.Stdin = bytes.NewBufferString(html)
	cmd.Dir = filepath.Join("..", "satori-reference")

	out, err := cmd.Output()
	if err != nil {
		errMsg := fmt.Sprintf("bun error: %v", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			errMsg += " stderr: " + string(exitErr.Stderr)
		}
		return &RenderResult{Error: errMsg}, &RenderResult{Error: errMsg}
	}

	var jr jsResult
	if err := json.Unmarshal(out, &jr); err != nil {
		return &RenderResult{Error: "json parse: " + err.Error()}, &RenderResult{Error: "json parse: " + err.Error()}
	}

	satori := jr.Satori
	takumi := jr.Takumi
	if satori == nil {
		satori = &RenderResult{Error: "not available"}
	}
	if takumi == nil {
		takumi = &RenderResult{Error: "not available"}
	}
	return satori, takumi
}

func renderAll(html string) *AllResults {
	ogreResult := renderOgre(html)

	satori, takumi := renderSatoriTakumi(html)

	return &AllResults{
		Ogre:   ogreResult,
		Satori: satori,
		Takumi: takumi,
	}
}
