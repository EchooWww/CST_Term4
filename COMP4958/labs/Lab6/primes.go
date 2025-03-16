package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func primes(m int, n int) []int {
	if n < 2 {
		return []int{}
	}

	isPrime := make([]bool, n)
	for i := 2; i < n; i++ {
		isPrime[i] = true
	}

	for i := 2; i*i < n; i++ {
		if isPrime[i] {
			for j := i * i; j < n; j += i {
				isPrime[j] = false
			}
		}
	}

	primes := []int{}
	for i := m; i < n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func getSortedDigits(n int) string {
	s := strconv.Itoa(n)
	chars := strings.Split(s, "")
	sort.Strings(chars)
	return strings.Join(chars, "")
}

func findLargestPermutationGroup(primes []int) []int {
	groups := make(map[string][]int)

	for _, prime := range primes {
		key := getSortedDigits(prime)
		groups[key] = append(groups[key], prime)
	}

	largestGroup := []int{}
	maxSize := 0

	for _, group := range groups {
		if len(group) > maxSize {
			largestGroup = group
			maxSize = len(group)
		}
	}

	return largestGroup
}

func main() {
	primesNumbers := primes(100000,1000000)
	largestGroup := findLargestPermutationGroup(primesNumbers)
	fmt.Println(len(largestGroup))
}