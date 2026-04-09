package font

import (
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

func (m *Manager) LoadDefaults() error {
	sources := []FontSource{
		{Name: "sans-serif", Weight: 400, Style: "normal", Data: goregular.TTF},
		{Name: "sans-serif", Weight: 700, Style: "normal", Data: gobold.TTF},
		{Name: "default", Weight: 400, Style: "normal", Data: goregular.TTF},
		{Name: "default", Weight: 700, Style: "normal", Data: gobold.TTF},
	}
	for _, src := range sources {
		if err := m.LoadFont(src); err != nil {
			return err
		}
	}
	return nil
}
