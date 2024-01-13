package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readlines(path string, c chan string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		c <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	close(c)
}

func fixWordNumbers(text string, search_strings []string) int {
	index := len(text)
	var first_number_pos int

	for p, s := range search_strings {
		idx := strings.Index(text, s)
		if idx != -1 && idx < index {
			index = idx
			first_number_pos = p
		}
	}

	index = -1
	var last_number_pos int

	for p, s := range search_strings {
		idx := strings.LastIndex(text, s)
		if idx != -1 && idx > index {
			index = idx
			last_number_pos = p
		}
	}

	v1 := first_number_pos%9 + 1
	v2 := last_number_pos%9 + 1

	return v1*10 + v2
}

func sum(vals []int) int {
	acc := 0

	for _, n := range vals {
		acc = acc + n
	}

	return acc
}

func part1() int {
	lines := make(chan string)

	go readlines("./data.txt", lines)

	vals := []int{}

	for text := range lines {
		search_strings := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
		num := fixWordNumbers(text, search_strings)

		vals = append(vals, num)
	}

	return sum(vals)
}

func part2() int {
	lines := make(chan string)

	go readlines("./data.txt", lines)

	vals := []int{}

	for text := range lines {
		search_strings := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
		num := fixWordNumbers(text, search_strings)

		vals = append(vals, num)
	}

	return sum(vals)
}

func main() {
	fmt.Println("part 1", part1())
	fmt.Println("part 2", part2())
}
