package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mengzhuo/nabhash"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("nabsum <filename>")
	}

	for _, fn := range args[1:] {
		if stat, err := os.Stat(fn); os.IsNotExist(err) || stat.IsDir() {
			continue
		}
		f, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		h := nabhash.New()
		io.Copy(h, f)
		fmt.Printf("%x\t%s\n", h.Sum(nil), fn)
	}
}
