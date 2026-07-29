package style

import (
	"testing"
)

// realisticCard is a typical Tailwind class list on a card component: mix of
// display, spacing, color, shadow, radius, and typography utilities.
var realisticCard = []string{
	"flex", "flex-col", "items-center", "justify-center",
	"gap-4", "px-6", "py-4",
	"bg-blue-500", "text-white", "text-lg", "font-semibold",
	"rounded-lg", "shadow-md", "border", "border-blue-600",
}

func BenchmarkResolveTailwind(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m := ResolveTailwind(realisticCard)
		if len(m) == 0 {
			b.Fatal("empty result")
		}
	}
}
