# nabhash

[![GoDoc Reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/mengzhuo/nabhash)
[![Build Status](https://travis-ci.org/mengzhuo/nabhash.svg?branch=master)](https://travis-ci.org/mengzhuo/nabhash/builds)

NABHash is an extremely fast Non-crypto-safe AES Based Hash algorithm for Big Data.

See https://nabhash.org for more.

## Benchmark
```
MacBook Pro (15-inch, 2016)
2.6 GHz Intel Core i7
16 GB 2133 MHz LPDDR3
```

```
goos: darwin
goarch: amd64
pkg: github.com/mengzhuo/nabhash
BenchmarkNABHash/8         	10000000	       116 ns/op	  68.41 MB/s
BenchmarkNABHash/16        	20000000	       111 ns/op	 143.67 MB/s
BenchmarkNABHash/32        	20000000	       111 ns/op	 287.51 MB/s
BenchmarkNABHash/64        	20000000	        99.7 ns/op	 641.66 MB/s
BenchmarkNABHash/128       	20000000	       100 ns/op	1267.72 MB/s
BenchmarkNABHash/256       	20000000	       104 ns/op	2453.07 MB/s
BenchmarkNABHash/512       	20000000	       106 ns/op	4812.62 MB/s
BenchmarkNABHash/1024      	20000000	       112 ns/op	9080.42 MB/s
BenchmarkNABHash/2048      	10000000	       129 ns/op	15858.82 MB/s
BenchmarkNABHash/4096      	10000000	       170 ns/op	24003.08 MB/s
BenchmarkNABHash/8192      	 5000000	       256 ns/op	31961.66 MB/s
BenchmarkNABHash/16384     	 3000000	       419 ns/op	39094.42 MB/s
BenchmarkNABHash/32768     	 2000000	       739 ns/op	44323.84 MB/s
BenchmarkNABHash/65536     	 1000000	      1391 ns/op	47086.78 MB/s
```

## Acknowledgements

NABHash makes use of the following open source projects:

* [Go](https://golang.org)
* [Meow](https://github.com/mmcloughlin/meow)
