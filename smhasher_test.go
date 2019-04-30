package nabhash

import (
	"math"
	"testing"
)

type HashSet struct {
	m map[string]struct{} // set of hashes added
	d *digest
	n int // number of hashes added
}

func newHashSet() *HashSet {
	return &HashSet{make(map[string]struct{}), New().(*digest), 0}
}

func (s *HashSet) check(t *testing.T) {
	const SLOP = 10.0
	collisions := s.n - len(s.m)
	t.Logf("%d/%d\n", len(s.m), s.n)
	pairs := int64(s.n) * int64(s.n-1) / 2
	expected := float64(pairs) / math.Pow(2.0, float64(Size))
	stddev := math.Sqrt(expected)
	if float64(collisions) > expected+SLOP*(3*stddev+1) {
		t.Errorf("unexpected number of collisions: got=%d mean=%f stddev=%f", collisions, expected, stddev)
	}
}

func (s *HashSet) add(x []byte) {
	s.m[string(s.d.Sum(x))] = struct{}{}
	s.n++
	s.d.Reset()
}

// All 0-3 byte strings have distinct hashes.
func TestSmhasherSmallKeys(t *testing.T) {
	h := newHashSet()
	var b [3]byte
	for i := 0; i < 256; i++ {
		b[0] = byte(i)
		h.add(b[:1])
		for j := 0; j < 256; j++ {
			b[1] = byte(j)
			h.add(b[:2])
			if !testing.Short() {
				for k := 0; k < 256; k++ {
					b[2] = byte(k)
					h.add(b[:3])
				}
			}
		}
	}
	h.check(t)
}
