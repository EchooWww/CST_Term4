package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	frequencies := make(map[string]int) // or : map[string]int{}
	for scanner.Scan() {
		frequencies[scanner.Text()]++
	}
	//fmt.Println(frequencies)
	words, max := []string{}, 0
	for w, c := range frequencies {
		if c > max {
			max = c
			words = []string{w}
		} else if c == max {
			words = append(words, w)
		}
	}
	fmt.Println(words, max)
}