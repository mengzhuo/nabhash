package nabhash

import "golang.org/x/sys/cpu"

func init() {
	if cpu.X86.HasAES && cpu.X86.HasAVX {
		block = blockAESNI
		final = finalAESNI
	}
}

//go:noescape
func blockAESNI(s *state, p []byte)

//go:noescape
func finalAESNI(s *state, l uint64)
