package nabhash

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
)

const hashSize = Size * 8

func randBytes(r *rand.Rand, b []byte) {
	for i := range b {
		b[i] = byte(r.Uint32())
	}
}

type BytesKey struct {
	b []byte
}

func (k *BytesKey) clear() {
	for i := range k.b {
		k.b[i] = 0
	}
}
func (k *BytesKey) random(r *rand.Rand) {
	randBytes(r, k.b)
}
func (k *BytesKey) bits() int {
	return len(k.b) * 8
}
func (k *BytesKey) flipBit(i int) {
	k.b[i>>3] ^= byte(1 << uint(i&7))
}
func (k *BytesKey) hash() []byte {
	h := New()
	h.Write(k.b)
	return h.Sum(nil)
}
func (k *BytesKey) name() string {
	return fmt.Sprintf("bytes%d", len(k.b))
}

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
	expected := float64(pairs) / math.Pow(2.0, float64(Size*8))
	stddev := math.Sqrt(expected)
	if float64(collisions) > expected+SLOP*(3*stddev+1) {
		t.Errorf("unexpected number of collisions: got=%d mean=%f stddev=%f", collisions, expected, stddev)
	}
}

func (s *HashSet) add(x []byte) {
	s.m[string(s.d.Sum(x))] = struct{}{}
	s.d.Reset()
	s.n++
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

// a string plus adding zeros must make distinct hashes
func TestSmhasherAppendedZeros(t *testing.T) {
	s := "hello" + strings.Repeat("\x00", 256)
	h := newHashSet()
	for i := 0; i <= len(s); i++ {
		h.add([]byte(s[:i]))
	}
	h.check(t)
}

// Different length strings of all zeros have distinct hashes.
func TestSmhasherZeros(t *testing.T) {
	N := 128 * 1024
	if testing.Short() {
		N = 1024
	}
	h := newHashSet()
	b := make([]byte, N)
	for i := 0; i <= N; i++ {
		h.add(b[:i])
	}
	h.check(t)
}

// Strings with up to two nonzero bytes all have distinct hashes.
func TestSmhasherTwoNonzero(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}
	h := newHashSet()
	for n := 2; n <= 16; n++ {
		twoNonZero(h, n)
	}
	h.check(t)
}
func twoNonZero(h *HashSet, n int) {
	b := make([]byte, n)

	// all zero
	h.add(b)

	// one non-zero byte
	for i := 0; i < n; i++ {
		for x := 1; x < 256; x++ {
			b[i] = byte(x)
			h.add(b)
			b[i] = 0
		}
	}

	// two non-zero bytes
	for i := 0; i < n; i++ {
		for x := 1; x < 256; x++ {
			b[i] = byte(x)
			for j := i + 1; j < n; j++ {
				for y := 1; y < 256; y++ {
					b[j] = byte(y)
					h.add(b)
					b[j] = 0
				}
			}
			b[i] = 0
		}
	}
}

// Test strings with repeats, like "abcdabcdabcdabcd..."
func TestSmhasherCyclic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}
	r := rand.New(rand.NewSource(1234))
	const REPEAT = 8
	const N = 1000000
	for n := 4; n <= 12; n++ {
		h := newHashSet()
		b := make([]byte, REPEAT*n)
		for i := 0; i < N; i++ {
			b[0] = byte(i * 79 % 97)
			b[1] = byte(i * 43 % 137)
			b[2] = byte(i * 151 % 197)
			b[3] = byte(i * 199 % 251)
			randBytes(r, b[4:n])
			for j := n; j < n*REPEAT; j++ {
				b[j] = b[j-n]
			}
			h.add(b)
		}
		h.check(t)
	}
}

// Sanity checks.
// hash should not depend on values outside key.
// hash should not depend on alignment.
func TestSmhasherSanity(t *testing.T) {
	r := rand.New(rand.NewSource(1234))
	const REP = 10
	const KEYMAX = 128
	const PAD = 16
	const OFFMAX = 16
	hb := New()
	hc := New()
	for k := 0; k < REP; k++ {
		for n := 0; n < KEYMAX; n++ {
			for i := 0; i < OFFMAX; i++ {
				var b [KEYMAX + OFFMAX + 2*PAD]byte
				var c [KEYMAX + OFFMAX + 2*PAD]byte
				randBytes(r, b[:])
				randBytes(r, c[:])
				copy(c[PAD+i:PAD+i+n], b[PAD:PAD+n])
				hb.Write(b[PAD : PAD+n])
				rb := hb.Sum(nil)
				hc.Write(c[PAD+i : PAD+i+n])
				rc := hc.Sum(nil)
				if !bytes.Equal(rc, rb) {
					t.Errorf("hash depends on bytes outside key")
				}
				hb.Reset()
				hc.Reset()
			}
		}
	}
}

// Test strings with only a few bits set
func TestSmhasherSparse(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}
	sparse(t, 32, 6)
	sparse(t, 40, 6)
	sparse(t, 48, 5)
	sparse(t, 56, 5)
	sparse(t, 64, 5)
	sparse(t, 96, 4)
	sparse(t, 256, 3)
	sparse(t, 2048, 2)
}
func sparse(t *testing.T, n int, k int) {
	b := make([]byte, n/8)
	h := newHashSet()
	setbits(h, b, 0, k)
	h.check(t)
}

// set up to k bits at index i and greater
func setbits(h *HashSet, b []byte, i int, k int) {
	h.add(b)
	if k == 0 {
		return
	}
	for j := i; j < len(b)*8; j++ {
		b[j/8] |= byte(1 << uint(j&7))
		setbits(h, b, j+1, k-1)
		b[j/8] &= byte(^(1 << uint(j&7)))
	}
}

// Test all possible combinations of n blocks from the set s.
// "permutation" is a bad name here, but it is what Smhasher uses.
func TestSmhasherPermutation(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping in short mode")
	}
	permutation(t, []uint32{0, 1, 2, 3, 4, 5, 6, 7}, 8)
	permutation(t, []uint32{0, 1 << 29, 2 << 29, 3 << 29, 4 << 29, 5 << 29, 6 << 29, 7 << 29}, 8)
	permutation(t, []uint32{0, 1}, 20)
	permutation(t, []uint32{0, 1 << 31}, 20)
	permutation(t, []uint32{0, 1, 2, 3, 4, 5, 6, 7, 1 << 29, 2 << 29, 3 << 29, 4 << 29, 5 << 29, 6 << 29, 7 << 29}, 6)
}
func permutation(t *testing.T, s []uint32, n int) {
	b := make([]byte, n*4)
	h := newHashSet()
	genPerm(h, b, s, 0)
	h.check(t)
}
func genPerm(h *HashSet, b []byte, s []uint32, n int) {
	h.add(b[:n])
	if n == len(b) {
		return
	}
	for _, v := range s {
		b[n] = byte(v)
		b[n+1] = byte(v >> 8)
		b[n+2] = byte(v >> 16)
		b[n+3] = byte(v >> 24)
		genPerm(h, b, s, n+4)
	}
}

// All keys of the form prefix + [A-Za-z0-9]*N + suffix.
func TestSmhasherText(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}
	text(t, "Foo", "Bar")
	text(t, "FooBar", "")
	text(t, "", "FooBar")
}
func text(t *testing.T, prefix, suffix string) {
	const N = 4
	const S = "ABCDEFGHIJKLMNOPQRSTabcdefghijklmnopqrst0123456789"
	const L = len(S)
	b := make([]byte, len(prefix)+N+len(suffix))
	copy(b, prefix)
	copy(b[len(prefix)+N:], suffix)
	h := newHashSet()
	c := b[len(prefix):]
	for i := 0; i < L; i++ {
		c[0] = S[i]
		for j := 0; j < L; j++ {
			c[1] = S[j]
			for k := 0; k < L; k++ {
				c[2] = S[k]
				for x := 0; x < L; x++ {
					c[3] = S[x]
					h.add(b)
				}
			}
		}
	}
	h.check(t)
}

// All bit rotations of a set of distinct keys
func TestSmhasherWindowed(t *testing.T) {
	windowed(t, &BytesKey{make([]byte, 128)})
}
func windowed(t *testing.T, k *BytesKey) {

	if testing.Short() {
		t.Skip("Skipping in short mode")
	}
	const BITS = 16

	for r := 0; r < k.bits(); r++ {
		h := newHashSet()
		for i := 0; i < 1<<BITS; i++ {
			k.clear()
			for j := 0; j < BITS; j++ {
				if i>>uint(j)&1 != 0 {
					k.flipBit((j + r) % k.bits())
				}
			}
			h.add(k.hash())
		}
		h.check(t)
	}
}
