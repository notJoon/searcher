package boyermoore

import (
	"testing"
)

func TestStringSearch(t *testing.T) {
	tests := []struct {
		name         string
		pattern      string
		text         string
		ignoreCase   bool
		wantAll      []int // 기대되는 모든 매칭 위치
		wantFirst    int   // 기대되는 첫 매칭 위치
		wantContains bool
		wantCount    int
	}{
		{
			name:         "Basic match",
			pattern:      "ABC",
			text:         "ZZZABCZZZ",
			ignoreCase:   false,
			wantAll:      []int{3},
			wantFirst:    3,
			wantContains: true,
			wantCount:    1,
		},
		{
			name:         "No match",
			pattern:      "ABC",
			text:         "ZZZABZ",
			ignoreCase:   false,
			wantAll:      []int{},
			wantFirst:    -1,
			wantContains: false,
			wantCount:    0,
		},
		{
			name:         "Multiple matches",
			pattern:      "AB",
			text:         "ABABAB",
			ignoreCase:   false,
			wantAll:      []int{0, 2, 4},
			wantFirst:    0,
			wantContains: true,
			wantCount:    3,
		},
		{
			name:         "Ignore case",
			pattern:      "AbC",
			text:         "zzZabcZZZAbCZZabcdZZ",
			ignoreCase:   true,
			wantAll:      []int{3, 9, 14}, // "abc" -> index 3, "AbC" -> index 9, "abcd" -> index 14 match "AbC" prefix
			wantFirst:    3,
			wantContains: true,
			wantCount:    3,
		},
		{
			name:         "Empty pattern",
			pattern:      "",
			text:         "ABC",
			ignoreCase:   false,
			wantAll:      []int{},
			wantFirst:    -1,
			wantContains: false,
			wantCount:    0,
		},
		{
			name:         "Pattern longer than text",
			pattern:      "ABCDEFG",
			text:         "ABC",
			ignoreCase:   false,
			wantAll:      []int{},
			wantFirst:    -1,
			wantContains: false,
			wantCount:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bm := New(tc.pattern, tc.ignoreCase)

			gotAll := bm.FindAll(tc.text)
			if !equalIntSlices(gotAll, tc.wantAll) {
				t.Errorf("FindAll(%q) = %v; want %v", tc.text, gotAll, tc.wantAll)
			}

			gotFirst := bm.FindFirst(tc.text)
			if gotFirst != tc.wantFirst {
				t.Errorf("FindFirst(%q) = %d; want %d", tc.text, gotFirst, tc.wantFirst)
			}

			gotContains := bm.Contains(tc.text)
			if gotContains != tc.wantContains {
				t.Errorf("Contains(%q) = %v; want %v", tc.text, gotContains, tc.wantContains)
			}

			gotCount := bm.Count(tc.text)
			if gotCount != tc.wantCount {
				t.Errorf("Count(%q) = %d; want %d", tc.text, gotCount, tc.wantCount)
			}
		})
	}
}

// 바이트 슬라이스 검색(Test for FindAllBytes, FindFirstBytes, ContainsBytes, CountBytes)
func TestByteSearch(t *testing.T) {
	tests := []struct {
		name         string
		pattern      string
		data         []byte
		ignoreCase   bool
		wantAll      []int
		wantFirst    int
		wantContains bool
		wantCount    int
	}{
		{
			name:         "Basic match (bytes)",
			pattern:      "ABC",
			data:         []byte("ZZZABCZZZ"),
			ignoreCase:   false,
			wantAll:      []int{3},
			wantFirst:    3,
			wantContains: true,
			wantCount:    1,
		},
		{
			name:         "Ignore case (bytes)",
			pattern:      "AbC",
			data:         []byte("ZZabcZZABCZZAbcdZZ"),
			ignoreCase:   true,
			wantAll:      []int{2, 7, 12},
			wantFirst:    2,
			wantContains: true,
			wantCount:    3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bm := New(tc.pattern, tc.ignoreCase)

			gotAll := bm.FindAllBytes(tc.data)
			if !equalIntSlices(gotAll, tc.wantAll) {
				t.Errorf("FindAllBytes(%q) = %v; want %v", string(tc.data), gotAll, tc.wantAll)
			}

			gotFirst := bm.FindFirstBytes(tc.data)
			if gotFirst != tc.wantFirst {
				t.Errorf("FindFirstBytes(%q) = %d; want %d", string(tc.data), gotFirst, tc.wantFirst)
			}

			gotContains := bm.ContainsBytes(tc.data)
			if gotContains != tc.wantContains {
				t.Errorf("ContainsBytes(%q) = %v; want %v", string(tc.data), gotContains, tc.wantContains)
			}

			gotCount := bm.CountBytes(tc.data)
			if gotCount != tc.wantCount {
				t.Errorf("CountBytes(%q) = %d; want %d", string(tc.data), gotCount, tc.wantCount)
			}
		})
	}
}

func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
