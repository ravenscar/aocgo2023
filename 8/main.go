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

type node struct {
	L string
	R string
}

type data struct {
	directions []string
	lookup     map[string]node
}

func loadAll(path string) data {
	var directions []string
	c := make(chan string)

	go readlines(path, c)

	lookup := map[string]node{}

	for line := range c {
		if len(strings.Trim(line, " ")) == 0 {
			continue
		}
		if len(directions) == 0 {
			directions = strings.Split(line, "")
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			panic("couldn't split line on =")
		}
		source := strings.Trim(parts[0], " ")
		dest := strings.Trim(parts[1], " ()")
		parts = strings.Split(dest, ",")
		if len(parts) != 2 {
			panic("couldn't split line on ,")
		}
		left := strings.Trim(parts[0], " ")
		right := strings.Trim(parts[1], " ")

		lookup[source] = node{left, right}
	}

	return data{directions, lookup}
}

func part1(filepath string) int {
	acc := 0
	current_node := "AAA"
	data := loadAll(filepath)
	direction_pos := 0

	for current_node != "ZZZ" {
		acc = acc + 1
		if direction_pos == len(data.directions) {
			direction_pos = 0
		}
		direction := data.directions[direction_pos]
		if direction == "L" {
			current_node = data.lookup[current_node].L
		} else {
			current_node = data.lookup[current_node].R
		}

		direction_pos = direction_pos + 1
	}

	return acc
}

func part2lcm(filepath string) int {
	acc := 0
	data := loadAll(filepath)
	current_nodes := []string{}
	for k := range data.lookup {
		if k[len(k)-1] == 'A' {
			current_nodes = append(current_nodes, k)
		}
	}

	hits := make([]int, len(current_nodes))

	end_state := func() bool {
		for _, n := range hits {
			if n == 0 {
				return false
			}
		}

		return true
	}

	direction_pos := 0

	for !end_state() {
		acc = acc + 1
		if direction_pos == len(data.directions) {
			direction_pos = 0
		}
		direction := data.directions[direction_pos]

		for i, n := range current_nodes {
			if direction == "L" {
				current_nodes[i] = data.lookup[n].L
			} else {
				current_nodes[i] = data.lookup[n].R
			}
			n = current_nodes[i]
			if hits[i] == 0 && n[len(n)-1] == 'Z' {
				hits[i] = acc
			}
		}

		direction_pos = direction_pos + 1
	}

	fmt.Printf("hits: %v", hits)
	rest := hits[2:]
	return LCM(hits[0], hits[1], rest...)
}

// from https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func part2(filepath string) int {
	acc := 0
	data := loadAll(filepath)
	current_nodes := []string{}
	for k := range data.lookup {
		if k[len(k)-1] == 'A' {
			current_nodes = append(current_nodes, k)
		}
	}

	end_state := func() bool {
		for _, n := range current_nodes {
			if n[len(n)-1] != 'Z' {
				return false
			}
		}

		return true
	}

	direction_pos := 0

	for !end_state() {
		acc = acc + 1
		if direction_pos == len(data.directions) {
			direction_pos = 0
		}
		direction := data.directions[direction_pos]

		for i, n := range current_nodes {
			if direction == "L" {
				current_nodes[i] = data.lookup[n].L
			} else {
				current_nodes[i] = data.lookup[n].R
			}
		}

		direction_pos = direction_pos + 1
	}

	return acc
}

func main() {
	fmt.Printf("Part1: %d\n", part1("data.txt"))
	fmt.Printf("Part2: %d\n", part2lcm("data.txt"))
}
