package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func get_points(input string) []int {
	parts := strings.Split(input, " ")
	points := make([]int, len(parts))

	for i, s := range parts {
		v, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		points[i] = v
	}

	return points
}

func iterate_points(points []int) []int {
	if len(points) < 2 {
		panic("less than two points")
	}

	next := make([]int, len(points)-1)

	for i := 1; i < len(points); i++ {
		next[i-1] = points[i] - points[i-1]
	}

	return next
}

func end_state(points []int) bool {
	for _, p := range points {
		if p != 0 {
			return false
		}
	}

	return true
}

func get_point_arrays(input string) [][]int {
	point_arrays := [][]int{}

	current_points := get_points(input)
	point_arrays = append(point_arrays, current_points)

	for !end_state(current_points) {
		current_points = iterate_points(current_points)
		point_arrays = append(point_arrays, current_points)
	}

	return point_arrays
}

func get_next_point(arrs [][]int) int {
	acc := 0

	for i := len(arrs) - 2; i >= 0; i-- {
		acc = acc + arrs[i][len(arrs[i])-1]
	}

	return acc
}

func get_prev_point(arrs [][]int) int {
	acc := 0

	for i := len(arrs) - 2; i >= 0; i-- {
		acc = arrs[i][0] - acc
	}

	return acc
}

func part1(filepath string) int {
	c := make(chan string)
	go readlines(filepath, c)
	acc := 0

	for line := range c {
		arrs := get_point_arrays(line)
		next := get_next_point(arrs)
		acc = acc + next
	}

	return acc
}

func part2(filepath string) int {
	c := make(chan string)
	go readlines(filepath, c)
	acc := 0

	for line := range c {
		arrs := get_point_arrays(line)
		next := get_prev_point(arrs)
		acc = acc + next
	}

	return acc
}

func main() {
	fmt.Printf("part1: %d\n", part1("data.txt"))
	fmt.Printf("part2: %d\n", part2("data.txt"))
}
