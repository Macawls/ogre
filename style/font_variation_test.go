package style

import (
	"reflect"
	"testing"
)

func TestParseFontVariationSettings(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []VariationAxis
	}{
		{"empty", "", nil},
		{"normal keyword", "normal", nil},
		{"NORMAL keyword case-insensitive", "Normal", nil},
		{"single wght", `'wght' 900`, []VariationAxis{{Tag: "wght", Value: 900}}},
		{"double-quoted", `"wght" 700`, []VariationAxis{{Tag: "wght", Value: 700}}},
		{"multiple axes", `'wght' 900, 'SOFT' 0, 'WONK' 1`, []VariationAxis{
			{Tag: "wght", Value: 900},
			{Tag: "SOFT", Value: 0},
			{Tag: "WONK", Value: 1},
		}},
		{"whitespace", `  'wght'   900  ,   'opsz'  72  `, []VariationAxis{
			{Tag: "wght", Value: 900},
			{Tag: "opsz", Value: 72},
		}},
		{"fractional value", `'wght' 425.5`, []VariationAxis{{Tag: "wght", Value: 425.5}}},
		{"negative value", `'slnt' -10`, []VariationAxis{{Tag: "slnt", Value: -10}}},

		{"missing quotes", `wght 900`, nil},
		{"tag not 4 chars", `'wgh' 900`, nil},
		{"tag 5 chars", `'wghtx' 900`, nil},
		{"missing value", `'wght'`, nil},
		{"non-numeric value", `'wght' bold`, nil},
		{"mix of good and bad drops bad", `'wght' 900, garbage, 'opsz' 12`, []VariationAxis{
			{Tag: "wght", Value: 900},
			{Tag: "opsz", Value: 12},
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ParseFontVariationSettings(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ParseFontVariationSettings(%q) = %+v, want %+v", tc.in, got, tc.want)
			}
		})
	}
}

func TestFormatFontVariationSettings(t *testing.T) {
	cases := []struct {
		name string
		in   []VariationAxis
		want string
	}{
		{"empty", nil, ""},
		{"single", []VariationAxis{{Tag: "wght", Value: 900}}, `'wght' 900`},
		{"multiple", []VariationAxis{
			{Tag: "wght", Value: 900},
			{Tag: "SOFT", Value: 0},
		}, `'wght' 900, 'SOFT' 0`},
		{"fractional", []VariationAxis{{Tag: "wght", Value: 425.5}}, `'wght' 425.5`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatFontVariationSettings(tc.in)
			if got != tc.want {
				t.Errorf("FormatFontVariationSettings(%+v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
