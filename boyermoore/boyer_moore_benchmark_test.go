package boyermoore

import (
	"math/rand"
	"testing"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateBenchmarkData(patternLen, textLen int) (pattern, text string) {
	pattern = generateRandomString(patternLen)
	text = generateRandomString(textLen)
	return
}

func BenchmarkFindAll(b *testing.B) {
	benchmarks := []struct {
		name       string
		patternLen int
		textLen    int
		ignoreCase bool
	}{
		{"Short Pattern (5) in Short Text (100)", 5, 100, false},
		{"Short Pattern (5) in Long Text (1000)", 5, 1000, false},
		{"Medium Pattern (20) in Medium Text (500)", 20, 500, false},
		{"Long Pattern (50) in Long Text (2000)", 50, 2000, false},
		{"Case Insensitive Search", 10, 1000, true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			pattern, text := generateBenchmarkData(bm.patternLen, bm.textLen)
			matcher := New(pattern, bm.ignoreCase)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matcher.FindAll(text)
			}
		})
	}
}

func BenchmarkFindFirst(b *testing.B) {
	benchmarks := []struct {
		name       string
		patternLen int
		textLen    int
		ignoreCase bool
	}{
		{"Short Pattern (5) in Short Text (100)", 5, 100, false},
		{"Short Pattern (5) in Long Text (1000)", 5, 1000, false},
		{"Medium Pattern (20) in Medium Text (500)", 20, 500, false},
		{"Long Pattern (50) in Long Text (2000)", 50, 2000, false},
		{"Case Insensitive Search", 10, 1000, true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			pattern, text := generateBenchmarkData(bm.patternLen, bm.textLen)
			matcher := New(pattern, bm.ignoreCase)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matcher.FindFirst(text)
			}
		})
	}
}

func BenchmarkContains(b *testing.B) {
	benchmarks := []struct {
		name       string
		patternLen int
		textLen    int
		ignoreCase bool
	}{
		{"Short Pattern (5) in Short Text (100)", 5, 100, false},
		{"Short Pattern (5) in Long Text (1000)", 5, 1000, false},
		{"Medium Pattern (20) in Medium Text (500)", 20, 500, false},
		{"Long Pattern (50) in Long Text (2000)", 50, 2000, false},
		{"Case Insensitive Search", 10, 1000, true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			pattern, text := generateBenchmarkData(bm.patternLen, bm.textLen)
			matcher := New(pattern, bm.ignoreCase)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matcher.Contains(text)
			}
		})
	}
}

func BenchmarkCount(b *testing.B) {
	benchmarks := []struct {
		name       string
		patternLen int
		textLen    int
		ignoreCase bool
	}{
		{"Short Pattern (5) in Short Text (100)", 5, 100, false},
		{"Short Pattern (5) in Long Text (1000)", 5, 1000, false},
		{"Medium Pattern (20) in Medium Text (500)", 20, 500, false},
		{"Long Pattern (50) in Long Text (2000)", 50, 2000, false},
		{"Case Insensitive Search", 10, 1000, true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			pattern, text := generateBenchmarkData(bm.patternLen, bm.textLen)
			matcher := New(pattern, bm.ignoreCase)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				matcher.Count(text)
			}
		})
	}
}
