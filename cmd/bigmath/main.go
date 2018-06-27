package main

import (
	"fmt"
	"strconv"
)

func main() {
	x := "2.1222"
	y := "3.9"

	var base = x
	var addr = y
	if len(y) > len(x) {
		base = y
		addr = x
	}

	reassembled := make([]byte, len(base))
	var prevRemainder int

	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '.' {
			reassembled[i] = '.'
			continue
		}

		var howmuch int
		slot, _ := strconv.Atoi(string(base[i]))

		if i < len(addr) {
			maths, _ := strconv.Atoi(string(addr[i]))
			howmuch = slot + maths + prevRemainder
		} else {
			howmuch = slot + prevRemainder
		}

		var carry bool
		if howmuch >= 10 {
			howmuch = howmuch - 10
			carry = true
		}

		reassembled[i] = byte(howmuch + '0')

		if carry {
			prevRemainder = 1
		} else {
			prevRemainder = 0
		}
	}

	fmt.Printf("reassembled: %+v\n", string(reassembled))

}
