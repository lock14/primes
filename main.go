package main

import (
	"fmt"
	"github.com/lock14/collections/bitset"
	"iter"
	"math"
	"os"
	"strconv"
)

var primeGenFuncs = map[string]func(int) iter.Seq[int]{
	"less-than": primesLessThan,
	"first-n":   firstNPrimes,
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: primes [less-than N|first-n N]")
		os.Exit(1)
	}
	cmd := os.Args[1]
	if primeGenFunc, ok := primeGenFuncs[cmd]; ok {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			if _, ioErr := fmt.Fprintf(os.Stderr, "not a valid number: %s\n", os.Args[2]); ioErr != nil {
				panic(ioErr)
			}
			os.Exit(1)
		}
		for p := range primeGenFunc(n) {
			fmt.Println(p)
		}
	} else {
		fmt.Printf("unknown command %q\n", cmd)
	}
}

func firstNPrimes(n int) iter.Seq[int] {
	if n < 1 {
		return func(yield func(int) bool) {}
	}
	bound := int(piInverse(float64(n))) + 4
	primes := primesLessThan(bound)
	return func(yield func(int) bool) {
		count := 0
		for p := range primes {
			if count == n {
				break
			}
			yield(p)
			count++
		}
	}
}

func primesLessThan(n int) iter.Seq[int] {
	b := bitset.New(bitset.NumBits(n))
	if n > 2 {
		b.Set(0)
		b.Set(1)
		for i := 4; i < n; i += 2 {
			b.Set(i)
		}
		for i := 3; (i*i) > i && (i*i) < n; i += 2 {
			if !b.Get(i) {
				// i is prime
				for j := i * i; j > i && j < n; j += i {
					b.Set(j)
				}
			}
		}
		b.FlipRange(0, n)
	}
	return b.SetBits()
}

func piInverse(y float64) float64 {
	var xOld float64
	x := y * math.Log(y)
	maxDiff := 0.5
	for x-xOld > maxDiff {
		xOld = x
		x = y * math.Log(x)
	}
	return x
}
