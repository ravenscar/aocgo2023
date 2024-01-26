package main

import (
	"bufio"
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

func getPoints(input string) []int {
  parts := strings.Split(input, " ")
  points := make([]int, len(parts))

  for i, s := range parts {
    v, err := strconv.Atoi(s)
    if (err != nil) {
      panic(err)
    }
    points[i] = v
  }

  return points
}
