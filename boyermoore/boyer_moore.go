package boyermoore

// BoyerMoore represents a pattern matcher using the Boyer-Moore algorithm.
// It contains the pattern, case sensitivity option, and precomputed
// bad character & good suffix shift tables.
type BoyerMoore struct {
	pat        []byte   // pattern (converted to lowercase if ignoreCase is true)
	ignoreCase bool     // case insensitivity flag
	bcShift    [256]int // bad character shift table
	gsShift    []int    // good suffix shift table
}

// New creates a new BoyerMoore matcher for the given pattern.
// If ignoreCase is true, the search will be case-insensitive.
func New(pattern string, ignoreCase bool) *BoyerMoore {
	if len(pattern) == 0 {
		return &BoyerMoore{
			pat:        make([]byte, 0),
			ignoreCase: ignoreCase,
			bcShift:    [256]int{},
			gsShift:    make([]int, 0),
		}
	}
	p := []byte(pattern)

	// Convert pattern to lowercase if case-insensitive search is requested
	if ignoreCase {
		for i := 0; i < len(p); i++ {
			c := p[i]
			// Consider only ASCII range ('A'~'Z')
			if c >= 'A' && c <= 'Z' {
				p[i] = c + ('a' - 'A')
			}
		}
	}

	bm := &BoyerMoore{
		pat:        p,
		ignoreCase: ignoreCase,
		gsShift:    make([]int, len(p)),
	}

	bm.buildBadCharShift()
	bm.buildGoodSuffixShift()

	return bm
}

// FindAll returns all starting indices where the pattern matches in the text.
// Returns an empty slice if no matches are found.
func (bm *BoyerMoore) FindAll(txt string) []int {
	return bm._findAll([]byte(txt))
}

// FindAllBytes returns all starting indices where the pattern matches in the byte slice.
// Returns an empty slice if no matches are found.
func (bm *BoyerMoore) FindAllBytes(data []byte) []int {
	return bm._findAll(data)
}

// FindFirst returns the index of the first occurrence of the pattern in the text.
// Returns -1 if the pattern is not found.
func (bm *BoyerMoore) FindFirst(txt string) int {
	res := bm.FindAll(txt)
	if len(res) > 0 {
		return res[0]
	}
	return -1
}

// FindFirstBytes returns the index of the first occurrence of the pattern in the byte slice.
// Returns -1 if the pattern is not found.
func (bm *BoyerMoore) FindFirstBytes(data []byte) int {
	res := bm.FindAllBytes(data)
	if len(res) > 0 {
		return res[0]
	}
	return -1
}

// Contains reports whether the pattern appears in the text.
func (bm *BoyerMoore) Contains(txt string) bool {
	return bm.FindFirst(txt) != -1
}

// ContainsBytes reports whether the pattern appears in the byte slice.
func (bm *BoyerMoore) ContainsBytes(data []byte) bool {
	return bm.FindFirstBytes(data) != -1
}

// Count returns the number of non-overlapping occurrences of the pattern in the text.
func (bm *BoyerMoore) Count(txt string) int {
	return len(bm.FindAll(txt))
}

// CountBytes returns the number of non-overlapping occurrences of the pattern in the byte slice.
func (bm *BoyerMoore) CountBytes(data []byte) int {
	return len(bm.FindAllBytes(data))
}

// _findAll is an internal method that implements the Boyer-Moore search algorithm.
// It returns all indices where the pattern matches in the given byte slice.
func (bm *BoyerMoore) _findAll(data []byte) []int {
	var results []int
	m := len(bm.pat)
	n := len(data)
	if m == 0 || n == 0 || m > n {
		return results
	}

	s := 0 // current text position
	for s <= n-m {
		j := m - 1
		// Check pattern match from right to left
		for j >= 0 && bm.pat[j] == bm.normChar(data[s+j]) {
			j--
		}

		if j < 0 {
			// Pattern fully matched
			results = append(results, s)
			// Use bad character shift
			if s+m < n {
				s += m - bm.bcShift[bm.normChar(data[s+m])]
			} else {
				s++
			}
		} else {
			// Mismatch occurred
			badCharShift := j - bm.bcShift[bm.normChar(data[s+j])]
			goodSuffixShift := bm.gsShift[j]
			if badCharShift < 1 {
				badCharShift = 1
			}
			if badCharShift > goodSuffixShift {
				s += badCharShift
			} else {
				s += goodSuffixShift
			}
		}
	}
	return results
}

// normChar normalizes a byte for case-insensitive comparison.
// If ignoreCase is true, converts ASCII uppercase letters to lowercase.
func (bm *BoyerMoore) normChar(c byte) byte {
	if bm.ignoreCase && c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}

// buildBadCharShift constructs the bad character shift table for the pattern.
func (bm *BoyerMoore) buildBadCharShift() {
	// Initialize with -1
	for i := range bm.bcShift {
		bm.bcShift[i] = -1
	}
	// Calculate shifts based on pattern (already in lowercase if ignoreCase is true)
	for i := 0; i < len(bm.pat); i++ {
		bm.bcShift[bm.pat[i]] = i
	}
}

// buildGoodSuffixShift constructs the good suffix shift table for the pattern.
func (bm *BoyerMoore) buildGoodSuffixShift() {
	m := len(bm.pat)
	bm.gsShift = make([]int, m)
	suffix := make([]int, m)
	suffix[m-1] = m
	g := m - 1
	f := m - 1

	// Calculate suffix array
	for i := m - 2; i >= 0; i-- {
		if i > g && suffix[i+m-1-f] < i-g {
			suffix[i] = suffix[i+m-1-f]
		} else {
			g = i
			f = i
			for g >= 0 && bm.pat[g] == bm.pat[g+m-1-f] {
				g--
			}
			suffix[i] = f - g
		}
	}

	// Initialize good suffix table
	for i := 0; i < m; i++ {
		bm.gsShift[i] = m
	}

	j := 0
	for i := m - 1; i >= 0; i-- {
		if suffix[i] == i+1 {
			for j < m-1-i {
				if bm.gsShift[j] == m {
					bm.gsShift[j] = m - 1 - i
				}
				j++
			}
		}
	}
	for i := 0; i < m-1; i++ {
		bm.gsShift[m-1-suffix[i]] = m - 1 - i
	}
}
