package ahocorasick

import (
	"reflect"
	"testing"
)

func TestAhoCorasickStringSearch(t *testing.T) {
	type testCase struct {
		name         string
		patterns     []string
		text         string
		ignoreCase   bool
		wantMatches  []ACMatch
		wantContains bool
		wantCount    int
	}
	tests := []testCase{
		{
			name:     "Basic multiple patterns",
			patterns: []string{"he", "she", "his", "hers"},
			// character index:
			// 'u' -> 0
			// 's' -> 1
			// 'h' -> 2
			// 'e' -> 3
			// 'r' -> 4
			// 's' -> 5
			text:       "ushers",
			ignoreCase: false,
			wantMatches: []ACMatch{
				{PatternIndex: 1, Start: 1, End: 3}, // "she" -> [1..3]
				{PatternIndex: 0, Start: 2, End: 3}, // "he"  -> [2..3]
				{PatternIndex: 3, Start: 2, End: 5}, // "hers"-> [2..5]
			},
			wantContains: true,
			wantCount:    3,
		},
		{
			name:         "No match",
			patterns:     []string{"cat", "dog"},
			text:         "mouse",
			ignoreCase:   false,
			wantMatches:  nil, // empty
			wantContains: false,
			wantCount:    0,
		},
		{
			name:       "Ignore case",
			patterns:   []string{"He", "She", "Hers"},
			text:       "USHERS",
			ignoreCase: true,
			wantMatches: []ACMatch{
				{PatternIndex: 1, Start: 1, End: 3},
				{PatternIndex: 0, Start: 2, End: 3},
				{PatternIndex: 2, Start: 2, End: 5},
			},
			wantContains: true,
			wantCount:    3,
		},
		{
			name:         "Empty text",
			patterns:     []string{"abc", "def"},
			text:         "",
			ignoreCase:   false,
			wantMatches:  nil,
			wantContains: false,
			wantCount:    0,
		},
		{
			name:         "Empty patterns",
			patterns:     []string{},
			text:         "some text",
			ignoreCase:   false,
			wantMatches:  nil,
			wantContains: false,
			wantCount:    0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ac := New(tc.patterns, tc.ignoreCase)

			gotMatches := ac.FindAll(tc.text)

			if len(gotMatches) != len(tc.wantMatches) {
				t.Fatalf("FindAll(%q) got %d matches, want %d matches",
					tc.text, len(gotMatches), len(tc.wantMatches))
			}

			for i := range gotMatches {
				got := gotMatches[i]
				want := tc.wantMatches[i]

				if got.PatternIndex != want.PatternIndex {
					t.Errorf("match[%d].PatternIndex = %d, want %d",
						i, got.PatternIndex, want.PatternIndex)
				}
				if got.Start != want.Start {
					t.Errorf("match[%d].Start = %d, want %d",
						i, got.Start, want.Start)
				}
				if got.End != want.End {
					t.Errorf("match[%d].End = %d, want %d",
						i, got.End, want.End)
				}
			}

			gotContains := ac.Contains(tc.text)
			if gotContains != tc.wantContains {
				t.Errorf("Contains(%q) got %v, want %v", tc.text, gotContains, tc.wantContains)
			}

			gotCount := ac.Count(tc.text)
			if gotCount != tc.wantCount {
				t.Errorf("Count(%q) got %d, want %d", tc.text, gotCount, tc.wantCount)
			}
		})
	}
}

func TestAhoCorasickByteSearch(t *testing.T) {
	type testCase struct {
		name         string
		patterns     []string
		data         []byte
		ignoreCase   bool
		wantMatches  []ACMatch
		wantContains bool
		wantCount    int
	}
	tests := []testCase{
		{
			name:       "Basic multiple patterns (bytes)",
			patterns:   []string{"he", "she", "his", "hers"},
			data:       []byte("ushers"),
			ignoreCase: false,
			wantMatches: []ACMatch{
				{PatternIndex: 0, Start: 3, End: 4}, // "he"
				{PatternIndex: 1, Start: 2, End: 4}, // "she"
				{PatternIndex: 3, Start: 2, End: 5}, // "hers"
			},
			wantContains: true,
			wantCount:    3,
		},
		{
			name:         "No match (bytes)",
			patterns:     []string{"cat", "dog"},
			data:         []byte("mouse"),
			ignoreCase:   false,
			wantMatches:  nil,
			wantContains: false,
			wantCount:    0,
		},
		{
			name:       "Ignore case (bytes)",
			patterns:   []string{"He", "She", "Hers"},
			data:       []byte("USHERS"),
			ignoreCase: true,
			wantMatches: []ACMatch{
				{PatternIndex: 0, Start: 3, End: 4},
				{PatternIndex: 1, Start: 2, End: 4},
				{PatternIndex: 2, Start: 2, End: 5},
			},
			wantContains: true,
			wantCount:    3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ac := New(tc.patterns, tc.ignoreCase)

			gotMatches := ac.FindAllBytes(tc.data)
			if !reflect.DeepEqual(gotMatches, tc.wantMatches) {
				t.Errorf("FindAllBytes(%q) got %v, want %v", string(tc.data), gotMatches, tc.wantMatches)
			}

			gotContains := ac.ContainsBytes(tc.data)
			if gotContains != tc.wantContains {
				t.Errorf("ContainsBytes(%q) got %v, want %v", string(tc.data), gotContains, tc.wantContains)
			}

			gotCount := ac.CountBytes(tc.data)
			if gotCount != tc.wantCount {
				t.Errorf("CountBytes(%q) got %d, want %d", string(tc.data), gotCount, tc.wantCount)
			}
		})
	}
}
