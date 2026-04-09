package style

import (
	"strconv"
	"strings"
)

// ValueUnit identifies the unit of a CSS value (px, em, rem, %, etc.).
type ValueUnit int

const (
	UnitPx ValueUnit = iota
	UnitEm
	UnitRem
	UnitPercent
	UnitVw
	UnitVh
	UnitAuto
	UnitNone
)

// Value is a numeric CSS value with a unit.
type Value struct {
	Raw  float64
	Unit ValueUnit
}

// ParseValue parses a CSS length or keyword (e.g. "12px", "auto") into a Value.
// ParseValue parses a CSS length value like "16px", "1.5em", or "auto".
func ParseValue(s string) Value {
	s = strings.TrimSpace(s)
	if s == "" {
		return Value{}
	}

	if s == "auto" {
		return Value{Unit: UnitAuto}
	}
	if s == "none" {
		return Value{Unit: UnitNone}
	}

	for _, suffix := range []struct {
		s string
		u ValueUnit
	}{
		{"rem", UnitRem},
		{"em", UnitEm},
		{"px", UnitPx},
		{"vw", UnitVw},
		{"vh", UnitVh},
		{"%", UnitPercent},
	} {
		if strings.HasSuffix(s, suffix.s) {
			num := strings.TrimSpace(s[:len(s)-len(suffix.s)])
			v, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return Value{}
			}
			return Value{Raw: v, Unit: suffix.u}
		}
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Value{}
	}
	return Value{Raw: v, Unit: UnitPx}
}

// ResolveContext provides the dimensions needed to resolve relative CSS units.
type ResolveContext struct {
	ParentFontSize float64
	RootFontSize   float64
	ViewportWidth  float64
	ViewportHeight float64
	ContainerSize  float64
}

func (v Value) Resolve(ctx ResolveContext) float64 {
	switch v.Unit {
	case UnitPx:
		return v.Raw
	case UnitEm:
		return v.Raw * ctx.ParentFontSize
	case UnitRem:
		return v.Raw * ctx.RootFontSize
	case UnitPercent:
		return v.Raw * ctx.ContainerSize / 100
	case UnitVw:
		return v.Raw * ctx.ViewportWidth / 100
	case UnitVh:
		return v.Raw * ctx.ViewportHeight / 100
	case UnitAuto, UnitNone:
		return 0
	default:
		return v.Raw
	}
}

func (v Value) Px() float64 {
	return v.Raw
}

func (v Value) IsAuto() bool {
	return v.Unit == UnitAuto
}

func (v Value) IsNone() bool {
	return v.Unit == UnitNone
}

func (v Value) IsZero() bool {
	return v.Raw == 0 && v.Unit != UnitAuto && v.Unit != UnitNone
}
