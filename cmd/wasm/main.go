//go:build js && wasm

package main

import (
	"fmt"
	"io"
	"net/http"
	"syscall/js"

	"github.com/macawls/ogre"
)

func main() {
	obj := js.Global().Get("Object").New()
	obj.Set("render", js.FuncOf(renderOneShot))
	obj.Set("createRenderer", js.FuncOf(createRenderer))
	js.Global().Set("Ogre", obj)
	js.Global().Set("ogreRender", js.FuncOf(legacyRender))
	select {}
}

func errResult(msg string) any {
	return js.ValueOf(map[string]any{"error": msg})
}

func renderOneShot(_ js.Value, args []js.Value) any {
	if len(args) < 1 || args[0].Type() != js.TypeString {
		return errResult("render(html, opts?): html string required")
	}
	var opts js.Value
	if len(args) >= 2 {
		opts = args[1]
	}
	return doRender(ogre.NewRenderer(), args[0].String(), opts)
}

func legacyRender(_ js.Value, args []js.Value) any {
	if len(args) < 1 || args[0].Type() != js.TypeString {
		return errResult("ogreRender(html, w?, h?): html required")
	}
	opts := js.Global().Get("Object").New()
	if len(args) > 1 && args[1].Type() == js.TypeNumber {
		opts.Set("width", args[1])
	}
	if len(args) > 2 && args[2].Type() == js.TypeNumber {
		opts.Set("height", args[2])
	}
	res := doRender(ogre.NewRenderer(), args[0].String(), opts)
	v, ok := res.(js.Value)
	if !ok {
		return res
	}
	if !v.Get("error").IsUndefined() {
		return v
	}
	return js.ValueOf(map[string]any{
		"svg":    v.Get("svg").String(),
		"width":  v.Get("width").Int(),
		"height": v.Get("height").Int(),
	})
}

func createRenderer(_ js.Value, _ []js.Value) any {
	r := ogre.NewRenderer()
	obj := js.Global().Get("Object").New()

	var funcs []js.Func
	loadFn := js.FuncOf(func(_ js.Value, a []js.Value) any {
		if len(a) < 1 || a[0].Type() != js.TypeObject {
			return errResult("loadFont(source): source object required")
		}
		src, err := parseFontSource(a[0])
		if err != nil {
			return errResult(err.Error())
		}
		if err := r.LoadFont(src); err != nil {
			return errResult(err.Error())
		}
		return js.ValueOf(map[string]any{"ok": true})
	})
	renderFn := js.FuncOf(func(_ js.Value, a []js.Value) any {
		if len(a) < 1 || a[0].Type() != js.TypeString {
			return errResult("render(html, opts?): html string required")
		}
		var o js.Value
		if len(a) >= 2 {
			o = a[1]
		}
		return doRender(r, a[0].String(), o)
	})
	freeFn := js.FuncOf(func(_ js.Value, _ []js.Value) any {
		for _, f := range funcs {
			f.Release()
		}
		return nil
	})
	funcs = []js.Func{loadFn, renderFn, freeFn}

	obj.Set("loadFont", loadFn)
	obj.Set("render", renderFn)
	obj.Set("free", freeFn)
	return obj
}

func doRender(r *ogre.Renderer, html string, opts js.Value) any {
	o, err := parseOptions(opts)
	if err != nil {
		return errResult(err.Error())
	}
	result, err := r.Render(html, o)
	if err != nil {
		return errResult(err.Error())
	}
	out := map[string]any{
		"format":      string(o.Format),
		"contentType": result.ContentType,
		"width":       result.Width,
		"height":      result.Height,
	}
	if o.Format == ogre.FormatSVG {
		out["svg"] = string(result.Data)
	} else {
		u8 := js.Global().Get("Uint8Array").New(len(result.Data))
		js.CopyBytesToJS(u8, result.Data)
		out["bytes"] = u8
	}
	return js.ValueOf(out)
}

func parseOptions(v js.Value) (ogre.Options, error) {
	var o ogre.Options
	if v.IsUndefined() || v.IsNull() || v.Type() != js.TypeObject {
		o.Format = ogre.FormatSVG
		return o, nil
	}
	if w := v.Get("width"); w.Type() == js.TypeNumber {
		o.Width = w.Int()
	}
	if h := v.Get("height"); h.Type() == js.TypeNumber {
		o.Height = h.Int()
	}
	if f := v.Get("format"); f.Type() == js.TypeString {
		o.Format = ogre.Format(f.String())
	}
	if q := v.Get("quality"); q.Type() == js.TypeNumber {
		o.Quality = q.Int()
	}
	if d := v.Get("debug"); d.Type() == js.TypeBoolean {
		o.Debug = d.Bool()
	}
	if e := v.Get("emojiProvider"); e.Type() == js.TypeString {
		o.EmojiProvider = e.String()
	}
	if m := v.Get("maxElements"); m.Type() == js.TypeNumber {
		o.MaxElements = m.Int()
	}
	if t := v.Get("tailwindVersion"); t.Type() == js.TypeString {
		if t.String() == "v4" {
			o.TailwindVersion = ogre.TailwindV4
		} else {
			o.TailwindVersion = ogre.TailwindV3
		}
	}
	if fonts := v.Get("fonts"); fonts.Type() == js.TypeObject {
		lenVal := fonts.Get("length")
		if lenVal.Type() == js.TypeNumber {
			for i := 0; i < lenVal.Int(); i++ {
				src, err := parseFontSource(fonts.Index(i))
				if err != nil {
					return o, fmt.Errorf("fonts[%d]: %w", i, err)
				}
				o.Fonts = append(o.Fonts, src)
			}
		}
	}
	if o.Format == "" {
		o.Format = ogre.FormatSVG
	}
	return o, nil
}

func parseFontSource(v js.Value) (ogre.FontSource, error) {
	var s ogre.FontSource
	if v.Type() != js.TypeObject {
		return s, fmt.Errorf("font source must be object")
	}
	if name := v.Get("name"); name.Type() == js.TypeString {
		s.Name = name.String()
	}
	if s.Name == "" {
		return s, fmt.Errorf("font source requires name")
	}
	if w := v.Get("weight"); w.Type() == js.TypeNumber {
		s.Weight = w.Int()
	}
	if st := v.Get("style"); st.Type() == js.TypeString {
		s.Style = st.String()
	}
	if d := v.Get("data"); !d.IsUndefined() && !d.IsNull() {
		length := d.Get("length")
		if length.Type() != js.TypeNumber {
			return s, fmt.Errorf("font data must be a Uint8Array")
		}
		buf := make([]byte, length.Int())
		js.CopyBytesToGo(buf, d)
		s.Data = buf
	}
	if u := v.Get("url"); u.Type() == js.TypeString {
		s.URL = u.String()
	}
	if len(s.Data) == 0 && s.URL != "" {
		data, err := fetchURL(s.URL)
		if err != nil {
			return s, fmt.Errorf("fetch %s: %w", s.URL, err)
		}
		s.Data = data
		s.URL = ""
	}
	if len(s.Data) == 0 {
		return s, fmt.Errorf("font source requires data or url")
	}
	return s, nil
}

func fetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
