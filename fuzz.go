// +build gofuzz,amd64

package nabhash

import (
	"bytes"
	"fmt"
)

func Fuzz(data []byte) int {
	h := New()
	h.Write(data)
	asm := h.Sum(nil)

	block = blockGeneric
	final = finalGeneric

	h = New()
	h.Write(data)
	gen := h.Sum(nil)

	if !bytes.Equal(gen, asm) {
		src := fmt.Sprintf("mismatch on %x\n%x\n%x", data, asm, gen)
		panic(src)
	}
	return 0
}
