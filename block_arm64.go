package nabhash

import "golang.org/x/sys/cpu"

func init() {
	if cpu.ARM64.HasAES {
		block = blockNEON
		final = finalNEON
	}
}

// go:noesacpe
func blockNEON(s *state, p []byte)

// go:noesacpe
func finalNEON(s *state, l uint64)
