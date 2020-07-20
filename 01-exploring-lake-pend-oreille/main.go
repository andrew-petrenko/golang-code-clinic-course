package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	file, err := os.Open("01-exploring-lake-pend-oreille/lake_data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	rdr.Comma = '\t'
	rdr.TrimLeadingSpace = true
	rows, err := rdr.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("Total Records: ", len(rows)-1)
	fmt.Println("Air Temp:\t", mean(rows, 1), median(rows, 1))
	fmt.Println("Barometric:\t", mean(rows, 2), median(rows, 2))
	fmt.Println("Wind Speed:\t", mean(rows, 7), median(rows, 7))
}

func mean(rows [][]string, idx int) float64 {
	var total float64
	for i, row := range rows {
		if i != 0 {
			val, _ := strconv.ParseFloat(row[idx], 64)
			total += val
		}
	}
	return total / float64(len(rows)-1)
}

func median(rows [][]string, idx int) float64 {
	var sorted []float64

	for i, row := range rows {
		if i != 0 {
			val, _ := strconv.ParseFloat(row[idx], 64)
			sorted = append(sorted, val)
		}
	}

	sort.Float64s(sorted)

	if len(sorted)%2 == 0 {
		middle := len(sorted) / 2
		higher := sorted[middle]
		lower := sorted[middle-1]

		return (higher + lower) / 2
	}

	middle := len(sorted) / 2

	return sorted[middle]
}
