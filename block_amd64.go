package nabhash

import "golang.org/x/sys/cpu"

func init() {
	if cpu.X86.HasAES && cpu.X86.HasAVX {
		block = blockAESNI
		final = finalAESNI
	}
}

// go:noesacpe
func blockAESNI(s *state, p []byte)

// go:noesacpe
func finalAESNI(s *state, l uint64)
