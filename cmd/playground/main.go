//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/macawls/ogre"
)

var renderer *ogre.Renderer

func render(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{"error": "missing html argument"})
	}

	html := args[0].String()
	width := 1200
	height := 630

	if len(args) > 1 && args[1].Type() == js.TypeNumber {
		width = args[1].Int()
	}
	if len(args) > 2 && args[2].Type() == js.TypeNumber {
		height = args[2].Int()
	}

	if renderer == nil {
		renderer = ogre.NewRenderer()
	}

	result, err := renderer.Render(html, ogre.Options{
		Width:  width,
		Height: height,
		Format: ogre.FormatSVG,
	})
	if err != nil {
		return js.ValueOf(map[string]any{"error": err.Error()})
	}

	return js.ValueOf(map[string]any{
		"svg":   string(result.Data),
		"width": result.Width,
		"height": result.Height,
	})
}

func main() {
	js.Global().Set("ogreRender", js.FuncOf(render))
	select {}
}
