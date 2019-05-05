// Package nabhash implements the 128-bit NonCrypto-AES-Based Hash Checksum.
// See https://nabhash.org for information.
package nabhash

import (
	"hash"
)

const (
	// Size of nabhash in bytes.
	Size = 16

	// BlockSize of nabhash in bytes.
	BlockSize = 64
)

var zeroData = make([]byte, BlockSize)

var initState = state{
	0x5A, 0x82, 0x79, 0x99, 0x6E, 0xD9, 0xEB, 0xA1,
	0x8F, 0x1B, 0xBC, 0xDC, 0xCA, 0x62, 0xC1, 0xD6,
	0x5A, 0x82, 0x79, 0x99, 0x6E, 0xD9, 0xEB, 0xA1,
	0x8F, 0x1B, 0xBC, 0xDC, 0xCA, 0x62, 0xC1, 0xD6,
	0x5A, 0x82, 0x79, 0x99, 0x6E, 0xD9, 0xEB, 0xA1,
	0x8F, 0x1B, 0xBC, 0xDC, 0xCA, 0x62, 0xC1, 0xD6,
	0x5A, 0x82, 0x79, 0x99, 0x6E, 0xD9, 0xEB, 0xA1,
	0x8F, 0x1B, 0xBC, 0xDC, 0xCA, 0x62, 0xC1, 0xD6,
}

var block = blockGeneric
var final = finalGeneric

type state [BlockSize]byte

type digest struct {
	h      state
	buf    state
	length uint64
	remain int
}

// New return a new hash.Hash computing the nabhash checksum.
func New() hash.Hash {
	d := &digest{}
	d.Reset()
	return d
}

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	if d.remain > 0 {
		n := copy(d.buf[d.remain:], p)
		d.remain += n
		if d.remain == BlockSize {
			block(&d.h, d.buf[:])
			d.remain -= BlockSize
		}
		p = p[n:]
	}

	if len(p) >= BlockSize {
		n := len(p) &^ (BlockSize - 1)
		block(&d.h, p[:n])
		p = p[n:]
	}

	if len(p) > 0 {
		d.remain = copy(d.buf[:], p)
	}
	d.length += uint64(nn)
	return
}

func (d *digest) Sum(b []byte) []byte {
	hash := d.checkSum()
	return append(b, hash[:]...)
}

func (d *digest) checkSum() (digest [Size]byte) {
	l := d.length
	if l%BlockSize != 0 {
		d.Write(zeroData[l%BlockSize:])
	}

	final(&d.h, l)
	copy(digest[:], d.h[:])

	return
}

func (d *digest) Reset() {
	d.remain = 0
	d.length = 0
	copy(d.h[:], initState[:])
	d.buf = [BlockSize]byte{}
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }
