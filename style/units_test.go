package style

import (
	"math"
	"testing"
)

func TestParseValuePx(t *testing.T) {
	v := ParseValue("16px")
	if v.Raw != 16 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueEm(t *testing.T) {
	v := ParseValue("1.5em")
	if v.Raw != 1.5 || v.Unit != UnitEm {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueRem(t *testing.T) {
	v := ParseValue("2rem")
	if v.Raw != 2 || v.Unit != UnitRem {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValuePercent(t *testing.T) {
	v := ParseValue("50%")
	if v.Raw != 50 || v.Unit != UnitPercent {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueVw(t *testing.T) {
	v := ParseValue("10vw")
	if v.Raw != 10 || v.Unit != UnitVw {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueVh(t *testing.T) {
	v := ParseValue("100vh")
	if v.Raw != 100 || v.Unit != UnitVh {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueAuto(t *testing.T) {
	v := ParseValue("auto")
	if !v.IsAuto() {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueNone(t *testing.T) {
	v := ParseValue("none")
	if !v.IsNone() {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValuePlainZero(t *testing.T) {
	v := ParseValue("0")
	if v.Raw != 0 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValuePlainNumber(t *testing.T) {
	v := ParseValue("1.5")
	if v.Raw != 1.5 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueNegative(t *testing.T) {
	v := ParseValue("-10px")
	if v.Raw != -10 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueNegativeEm(t *testing.T) {
	v := ParseValue("-0.5em")
	if v.Raw != -0.5 || v.Unit != UnitEm {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueDecimalPx(t *testing.T) {
	v := ParseValue("0.75px")
	if v.Raw != 0.75 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueWhitespace(t *testing.T) {
	v := ParseValue("  16px  ")
	if v.Raw != 16 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueEmpty(t *testing.T) {
	v := ParseValue("")
	if v.Raw != 0 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func TestParseValueInvalid(t *testing.T) {
	v := ParseValue("abc")
	if v.Raw != 0 || v.Unit != UnitPx {
		t.Fatalf("got %+v", v)
	}
}

func approx(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}

func TestResolvePx(t *testing.T) {
	v := Value{Raw: 16, Unit: UnitPx}
	if got := v.Resolve(ResolveContext{}); got != 16 {
		t.Fatalf("got %f", got)
	}
}

func TestResolveEm(t *testing.T) {
	v := Value{Raw: 1.5, Unit: UnitEm}
	got := v.Resolve(ResolveContext{ParentFontSize: 16})
	if !approx(got, 24) {
		t.Fatalf("got %f", got)
	}
}

func TestResolveRem(t *testing.T) {
	v := Value{Raw: 2, Unit: UnitRem}
	got := v.Resolve(ResolveContext{RootFontSize: 16})
	if !approx(got, 32) {
		t.Fatalf("got %f", got)
	}
}

func TestResolvePercent(t *testing.T) {
	v := Value{Raw: 50, Unit: UnitPercent}
	got := v.Resolve(ResolveContext{ContainerSize: 800})
	if !approx(got, 400) {
		t.Fatalf("got %f", got)
	}
}

func TestResolveVw(t *testing.T) {
	v := Value{Raw: 10, Unit: UnitVw}
	got := v.Resolve(ResolveContext{ViewportWidth: 1200})
	if !approx(got, 120) {
		t.Fatalf("got %f", got)
	}
}

func TestResolveVh(t *testing.T) {
	v := Value{Raw: 100, Unit: UnitVh}
	got := v.Resolve(ResolveContext{ViewportHeight: 630})
	if !approx(got, 630) {
		t.Fatalf("got %f", got)
	}
}

func TestResolveAuto(t *testing.T) {
	v := Value{Unit: UnitAuto}
	if got := v.Resolve(ResolveContext{}); got != 0 {
		t.Fatalf("got %f", got)
	}
}

func TestResolveNone(t *testing.T) {
	v := Value{Unit: UnitNone}
	if got := v.Resolve(ResolveContext{}); got != 0 {
		t.Fatalf("got %f", got)
	}
}

func TestIsZero(t *testing.T) {
	cases := []struct {
		v    Value
		want bool
	}{
		{Value{Raw: 0, Unit: UnitPx}, true},
		{Value{Raw: 0, Unit: UnitEm}, true},
		{Value{Raw: 1, Unit: UnitPx}, false},
		{Value{Unit: UnitAuto}, false},
		{Value{Unit: UnitNone}, false},
	}
	for _, c := range cases {
		if got := c.v.IsZero(); got != c.want {
			t.Errorf("IsZero(%+v) = %v, want %v", c.v, got, c.want)
		}
	}
}

func TestResolveNegative(t *testing.T) {
	v := Value{Raw: -10, Unit: UnitPx}
	if got := v.Resolve(ResolveContext{}); got != -10 {
		t.Fatalf("got %f", got)
	}
}

func TestResolveNegativePercent(t *testing.T) {
	v := Value{Raw: -25, Unit: UnitPercent}
	got := v.Resolve(ResolveContext{ContainerSize: 400})
	if !approx(got, -100) {
		t.Fatalf("got %f", got)
	}
}
