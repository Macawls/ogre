package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	ogre "github.com/macawls/ogre/v3"
	"github.com/macawls/ogre/v3/server"
)

func main() {
	flag.Usage = usage

	serve := flag.Bool("serve", false, "start HTTP server mode")
	port := flag.Int("port", 3000, "server port")
	render := flag.String("render", "", "path to HTML file to render")
	html := flag.String("html", "", "inline HTML to render")
	output := flag.String("output", "", "output file path (default stdout for SVG, required for PNG)")
	width := flag.Int("width", 1200, "canvas width")
	height := flag.Int("height", 630, "canvas height")
	format := flag.String("format", "svg", "output format: svg, png, or jpeg")
	twVersion := flag.String("tailwind-version", "v3", "Tailwind CSS version: v3 or v4")

	flag.Parse()

	switch {
	case *serve:
		runServer(*port)
	case *render != "":
		runRender(*render, *output, *width, *height, *format, *twVersion)
	case *html != "":
		runHTML(*html, *output, *width, *height, *format, *twVersion)
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func runServer(port int) {
	cfg := server.ConfigFromEnv()
	if port != 3000 {
		cfg.Addr = fmt.Sprintf(":%d", port)
	}
	srv := server.New(cfg)
	fmt.Fprintf(os.Stderr, "listening on %s\n", cfg.Addr)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func runRender(path, output string, width, height int, format, twVersion string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	renderAndWrite(string(data), output, width, height, format, twVersion)
}

func runHTML(html, output string, width, height int, format, twVersion string) {
	renderAndWrite(html, output, width, height, format, twVersion)
}

func renderAndWrite(html, output string, width, height int, format, twVersion string) {
	f := ogre.Format(strings.ToLower(format))
	if f != ogre.FormatSVG && f != ogre.FormatPNG && f != ogre.FormatJPEG {
		fmt.Fprintf(os.Stderr, "unsupported format: %s\n", format)
		os.Exit(1)
	}

	if (f == ogre.FormatPNG || f == ogre.FormatJPEG) && output == "" {
		fmt.Fprintf(os.Stderr, "error: --output is required for %s format\n", format)
		os.Exit(1)
	}

	tv := ogre.TailwindVersion(strings.ToLower(twVersion))
	if tv != "" && tv != ogre.TailwindV3 && tv != ogre.TailwindV4 {
		fmt.Fprintf(os.Stderr, "unsupported tailwind version: %s (want v3 or v4)\n", twVersion)
		os.Exit(1)
	}

	result, err := ogre.Render(html, ogre.Options{
		Width:           width,
		Height:          height,
		Format:          f,
		TailwindVersion: tv,
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

const helpText = `ogre — HTML/CSS to image renderer (SVG, PNG, JPEG). Pure Go, no headless browser.

Accepts the same HTML/CSS subset as Vercel Satori. Inline styles + Tailwind
utility classes only — no external stylesheets, no <style> blocks.

USAGE
  ogre --html '<div style="…">…</div>' [flags]         # inline HTML
  ogre --render page.html [flags]                       # file input
  ogre --serve [--port N]                               # HTTP server

FLAGS
  Input (choose one):
    --html    string  inline HTML source
    --render  string  path to an HTML file
    --serve           start HTTP server

  Output:
    --output  string  output file (required for png/jpeg; stdout for svg)
    --format  string  svg | png | jpeg               (default svg)
    --width   int     canvas width in px             (default 1200)
    --height  int     canvas height in px            (default 630)

  Tailwind:
    --tailwind-version string   v3 | v4              (default v3)

  Server:
    --port    int     server port (with --serve)    (default 3000)

EXAMPLES
  # Inline HTML → PNG
  ogre --html '<div class="flex w-full h-full items-center justify-center bg-blue-500 text-white text-6xl font-bold">Hello</div>' \
       --format png --output card.png

  # File → SVG on stdout
  ogre --render card.html > card.svg

  # HTTP server on port 8080
  ogre --serve --port 8080

WRITING HTML FOR OGRE (read this before generating markup)

  Styling:
    * Use inline style="…" attributes and Tailwind class="…" utilities.
    * DO NOT emit <link rel="stylesheet" …>, <style> blocks, or a Tailwind
      CDN <script> — Ogre parses inline attributes only. External CSS is
      never fetched; <style>/<script> tag contents may render as text.
    * No @media, :hover / :focus, ::before / ::after, animations,
      transitions, or calc(). Values are resolved once, statically.

  Layout — flexbox only (W3C spec):
    * <div> defaults to display: block, rendered internally as
      flex-direction: column with align-items: stretch.
    * DO NOT use display: grid, grid-template-*, float, or columns —
      grid falls through to the default (block/column-flex) and your
      layout will drift.
    * position: static | relative | absolute only. fixed and sticky are
      not supported (they degrade silently).

  Images (<img src> and background-image: url(…)):
    * data: URIs (self-contained) — preferred for deterministic renders.
    * http(s):// URLs — fetched with a 5s timeout, cached in-process.
    * DO NOT use local file paths — they render a gray placeholder.
    * No <picture>, srcset, or sizes — only <img src>.

  Fonts:
    * Ships with Go's default sans (regular + bold).
    * Any Google Fonts family works by name — fetched on first use.
    * DO NOT use @font-face — ignored. Pass custom fonts via the Go
      Options.Fonts field or the HTTP "fonts" array instead.
    * WOFF2 is not supported (use TTF, OTF, or WOFF).

  Supported that agents often skip:
    transform, transform-origin, filter (blur / grayscale / brightness),
    box-shadow, text-shadow, clip-path, opacity, border-radius,
    background-image linear-gradient / radial-gradient, object-fit,
    line-clamp, emoji, RTL / bidi (Arabic, Hebrew, etc.).

  Tailwind arbitrary values work: w-[123px], bg-[#ff5500], rotate-[15deg].
  In v4, bg-(--brand) expands to background-color: var(--brand).

FULL LLM CONTEXT
  https://ogre.macawls.dev/llms-full.txt        (canonical digest)
  https://ogre.macawls.dev/reference/css/       (property table)
  https://ogre.macawls.dev/reference/tailwind/  (utility list)
`

func usage() {
	fmt.Fprint(os.Stderr, helpText)
}
