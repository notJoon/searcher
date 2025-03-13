package ahocorasick

// ACMatch represents pattern matching information found in text
type ACMatch struct {
	PatternIndex int // which pattern in keywords
	Start        int // start index of the match
	End          int // end index of the match (inclusive)
}

// AhoCorasick is a struct that contains Aho-Corasick automaton for multiple pattern search
type AhoCorasick struct {
	keywords   [][]byte // patterns (may already be converted to lowercase)
	ignoreCase bool

	// trie nodes. node 0 is root.
	// ex: next[node][c] = transition
	// fail[node] = failure link
	// out[node] = list of pattern indices that end at this node
	next [][256]int
	fail []int
	out  [][]int
}

// New creates and returns an AhoCorasick struct with multiple patterns
func New(patterns []string, ignoreCase bool) *AhoCorasick {
	// Store keywords: if ignoreCase option is true, convert all to lowercase internally
	var kw [][]byte
	for _, p := range patterns {
		b := []byte(p)
		if ignoreCase {
			for i := range b {
				if b[i] >= 'A' && b[i] <= 'Z' {
					b[i] = b[i] + ('a' - 'A')
				}
			}
		}
		kw = append(kw, b)
	}

	ac := &AhoCorasick{
		keywords:   kw,
		ignoreCase: ignoreCase,
		// initially trie is empty, so allocate 1 node (root)
		next: make([][256]int, 1),
		fail: make([]int, 1),
		out:  make([][]int, 1),
	}

	ac.buildTrie()
	ac.buildFailureLinks()
	return ac
}

// FindAll finds all pattern matches (ACMatch) in text using Aho-Corasick
func (ac *AhoCorasick) FindAll(text string) []ACMatch {
	return ac._findAll([]byte(text))
}

// FindAllBytes finds all pattern matches (ACMatch) in byte slice using Aho-Corasick
func (ac *AhoCorasick) FindAllBytes(data []byte) []ACMatch {
	return ac._findAll(data)
}

// Contains returns whether any registered pattern matches in the text
func (ac *AhoCorasick) Contains(text string) bool {
	ms := ac.FindAll(text)
	return len(ms) > 0
}

// ContainsBytes returns whether any pattern matches in the byte slice
func (ac *AhoCorasick) ContainsBytes(data []byte) bool {
	ms := ac.FindAllBytes(data)
	return len(ms) > 0
}

// Count returns the number of **all** matches found in the text
func (ac *AhoCorasick) Count(text string) int {
	ms := ac.FindAll(text)
	return len(ms)
}

// CountBytes returns the number of all matches found in the byte slice
func (ac *AhoCorasick) CountBytes(data []byte) int {
	ms := ac.FindAllBytes(data)
	return len(ms)
}

// buildTrie inserts patterns from ac.keywords into the trie
func (ac *AhoCorasick) buildTrie() {
	for idx, k := range ac.keywords {
		node := 0 // start from root
		for _, c := range k {
			cc := c // (byte)
			if ac.next[node][cc] == 0 {
				// Create new node
				ac.next = append(ac.next, [256]int{})
				ac.fail = append(ac.fail, 0)
				ac.out = append(ac.out, []int{})
				ac.next[node][cc] = len(ac.next) - 1
			}
			node = ac.next[node][cc]
		}
		// Patterns ending at this node: idx
		ac.out[node] = append(ac.out[node], idx)
	}
}

// buildFailureLinks sets up failure links for each node using BFS method,
// and reflects the out information of nodes connected through fail links to the current node
func (ac *AhoCorasick) buildFailureLinks() {
	queue := []int{}

	// 1) Set up root(0)'s child nodes
	for c := 0; c < 256; c++ {
		nx := ac.next[0][c]
		if nx != 0 {
			// Set child's fail to 0(root)
			ac.fail[nx] = 0
			queue = append(queue, nx)
		} else {
			// Maintain fail[next[0][c]] = 0 even when not present
			ac.next[0][c] = 0
		}
	}

	// 2) Set failure links while running BFS
	for len(queue) > 0 {
		f := queue[0]
		queue = queue[1:]

		// for all edges c of f node
		for c := 0; c < 256; c++ {
			nx := ac.next[f][c]
			if nx != 0 {
				queue = append(queue, nx)
				// update failure link logic
				failTo := ac.fail[f]
				// follow c edge from failTo node
				ac.fail[nx] = ac.next[failTo][c]
				// inherit out information
				ac.out[nx] = append(ac.out[nx], ac.out[ac.fail[nx]]...)
			} else {
				// if no edge, follow fail[f] of current f node to the node connected by c edge
				ac.next[f][c] = ac.next[ac.fail[f]][c]
			}
		}
	}
}

// _findAll finds all matching patterns (ACMatch) in the byte slice data
func (ac *AhoCorasick) _findAll(data []byte) []ACMatch {
	var matches []ACMatch
	node := 0 // current node being searched in trie

	for i, c := range data {
		cc := c
		if ac.ignoreCase && cc >= 'A' && cc <= 'Z' {
			cc = cc + ('a' - 'A')
		}
		node = ac.next[node][cc]

		// Process all pattern indices in node(any node in trie)'s out
		if len(ac.out[node]) > 0 {
			for _, patIdx := range ac.out[node] {
				patLen := len(ac.keywords[patIdx])
				matches = append(matches, ACMatch{
					PatternIndex: patIdx,
					Start:        i - patLen + 1,
					End:          i,
				})
			}
		}
	}
	return matches
}
