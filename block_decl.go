// +build !amd64

package nabhash

//go:noescape
func block(d *digest, p []byte)

//go:noescape
func final(d *digest, p []byte)
