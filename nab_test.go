package nabhash

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

const ext = "nabsum"

func TestFromData(t *testing.T) {

	cur, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	filepath.Walk(filepath.Join(cur, "testdata"),
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				t.Error(err)
				return err
			}
			if info.IsDir() {
				return nil
			}

			base := filepath.Base(path)
			if filepath.Ext(base) != ext {
				return nil
			}

			t.Log(path)
			expect := base[:Size*2]
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			h := New()
			io.Copy(h, f)
			if got := fmt.Sprintf("%x", h.Sum(nil)); got != expect {
				t.Errorf("got=%s expect=%s", got, expect)
			}

			return nil
		})
}

func TestHexSmallSets(t *testing.T) {
	for _, gold := range []struct {
		target, hash string
	}{

		{"00", "741892770d81d644519a99a6bdfe072a"},
		{"0001", "95a85a5b60f29b9ec58569f92ebea60e"},
		{"000008", "4e3e54e8e3b0e245b50cb650721b3e47"},
		{"deadbeef", "9041d0d440fe93f69515b771619aa54c"},
		{"0000000FF1CE",
			"7e80664f949b119bdc9495cb0aebb335"},
		{"FEEDFACECAFEBEEF", "f250d662ca498c671d5fba0883137b60"},
		{"4d656e67205a68756f2069732074686520617574686f72206f66204e414248617368",
			"be7b4ad42ddd356e7abd8755902961dd"},
	} {
		target, err := hex.DecodeString(gold.target)
		if err != nil {
			t.Error(err)
			continue
		}
		h, err := hex.DecodeString(gold.hash)
		if err != nil {
			t.Error(err)
			continue
		}

		hasher := New()
		hasher.Write(target)
		got := hasher.Sum(nil)
		if !bytes.Equal(got, h) {
			t.Errorf("hash(%s) expect:%s got:%x", gold.target, gold.hash, got)
		}
	}
}

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
