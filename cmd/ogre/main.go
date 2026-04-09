package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	ogre "github.com/macawls/ogre"
	"github.com/macawls/ogre/server"
)

func main() {
	serve := flag.Bool("serve", false, "start HTTP server mode")
	port := flag.Int("port", 3000, "server port")
	render := flag.String("render", "", "path to HTML file to render")
	html := flag.String("html", "", "inline HTML to render")
	output := flag.String("output", "", "output file path (default stdout for SVG, required for PNG)")
	width := flag.Int("width", 1200, "canvas width")
	height := flag.Int("height", 630, "canvas height")
	format := flag.String("format", "svg", "output format: svg, png, or jpeg")

	flag.Parse()

	switch {
	case *serve:
		runServer(*port)
	case *render != "":
		runRender(*render, *output, *width, *height, *format)
	case *html != "":
		runHTML(*html, *output, *width, *height, *format)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func runServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	srv := server.New(server.Config{Addr: addr})
	fmt.Fprintf(os.Stderr, "listening on %s\n", addr)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func runRender(path, output string, width, height int, format string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	renderAndWrite(string(data), output, width, height, format)
}

func runHTML(html, output string, width, height int, format string) {
	renderAndWrite(html, output, width, height, format)
}

func renderAndWrite(html, output string, width, height int, format string) {
	f := ogre.Format(strings.ToLower(format))
	if f != ogre.FormatSVG && f != ogre.FormatPNG && f != ogre.FormatJPEG {
		fmt.Fprintf(os.Stderr, "unsupported format: %s\n", format)
		os.Exit(1)
	}

	if (f == ogre.FormatPNG || f == ogre.FormatJPEG) && output == "" {
		fmt.Fprintf(os.Stderr, "error: --output is required for %s format\n", format)
		os.Exit(1)
	}

	result, err := ogre.Render(html, ogre.Options{
		Width:  width,
		Height: height,
		Format: f,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "render error: %v\n", err)
		os.Exit(1)
	}

	if output == "" {
		if _, err := os.Stdout.Write(result.Data); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := os.WriteFile(output, result.Data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		os.Exit(1)
	}
}
