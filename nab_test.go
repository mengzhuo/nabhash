package nabhash

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestNABHashAsm(t *testing.T) {

	if runtime.FuncForPC(reflect.ValueOf(block).Pointer()).Name() == "blockGeneric" {
		t.Skip("no asm implement on this platform")
	}

	for _, v := range []string{
		"4e414248617368",
	} {
		s, err := hex.DecodeString(v)
		if err != nil {
			t.Errorf("hex:%s", err)
			continue
		}
		testCmp(t, s)
	}
}

func testCmp(t *testing.T, s []byte) {
	ob := block
	of := final

	defer func() {
		block = ob
		final = of
	}()

	h := New()
	h.Write(s)
	asm := fmt.Sprintf("%x", h.Sum(nil))
	block = blockGeneric
	final = finalGeneric

	h = New()
	h.Write(s)
	gen := fmt.Sprintf("%x", h.Sum(nil))
	if gen != asm {
		t.Errorf("asm:%s != gen:%s", asm, gen)
	}
}

func BenchmarkNABHash(b *testing.B) {

	for i := 8; i <= 65536; i <<= 1 {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			buf := make([]byte, i)
			h := New()
			b.SetBytes(int64(i))
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				h.Write(buf)
				h.Sum(nil)
			}
			return
		})
	}

}
