package main

import (
	"flag"
	"fmt"
	"github.com/lock14/collections/bitset"
	"iter"
	"math"
	"os"
)

var (
	firstK   = flag.Int("first-k", -1, "output the first k primes")
	lessThan = flag.Int("less-than", -1, "output primes less than the given value")
)

func main() {
	flag.Parse()
	if *firstK == -1 && *lessThan == -1 {
		flag.PrintDefaults()
		os.Exit(1)
	} else if *firstK != -1 && *lessThan != -1 {
		flag.PrintDefaults()
		os.Exit(1)
	} else if *firstK != -1 {
		for n := range firstKPrimes(*firstK) {
			fmt.Println(n)
		}
	} else if *lessThan != -1 {
		for n := range primesLessThan(*lessThan) {
			fmt.Println(n)
		}
	}

}

func firstKPrimes(k int) iter.Seq[int] {
	if k < 1 {
		return func(yield func(int) bool) {}
	}
	bound := int(piInverse(float64(k))) + 4
	primes := primesLessThan(bound)
	return func(yield func(int) bool) {
		count := 0
		for p := range primes {
			if count == k {
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
