package nabhash

import "golang.org/x/sys/cpu"

func init() {
	if cpu.ARM64.HasAES {
		block = blockNEON
		final = finalNEON
	}
}

//go:noescape
func blockNEON(s *state, p []byte)

//go:noescape
func finalNEON(s *state, l uint64)
