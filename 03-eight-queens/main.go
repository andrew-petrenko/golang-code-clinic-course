package main

import "fmt"

// permutation example
func main() {
	options := []string{"A", "B", "C", "D"}

	for _, first := range options {
		for _, second := range options {
			if second == first {
				continue
			}

			for _, third := range options {
				if third == second || third == first {
					continue
				}

				for _, fourth := range options {
					if fourth == third || fourth == second || fourth == first {
						continue
					}
					fmt.Println(first, second, third, fourth)
				}
			}
		}
	}
}
