package main

import (
	"bufio"
	"fmt"
	"math"
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

func toInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return i
}

type line_info struct {
	card_no         int
	winning_numbers []int
	line_numbers    []int
}

func parseLine(line string) line_info {
	var parts []string

	parts = strings.Split(line, ":")
	card := parts[0]
	rest := parts[1]

	parts = strings.Split(card, " ")
	card_no := toInt(parts[len(parts)-1])
	parts = strings.Split(rest, "|")
	winning_numbers := []int{}
	for _, v := range strings.Split(parts[0], " ") {
		v = strings.Trim(v, " ")
		if len(v) > 0 {
			winning_numbers = append(winning_numbers, toInt(v))
		}
	}
	line_numbers := []int{}
	for _, v := range strings.Split(parts[1], " ") {
		v = strings.Trim(v, " ")
		if len(v) > 0 {
			line_numbers = append(line_numbers, toInt(v))
		}
	}

	return line_info{
		card_no,
		winning_numbers,
		line_numbers,
	}
}

func contains(v int, c []int) bool {
	for _, i := range c {
		if i == v {
			return true
		}
	}

	return false
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func getLineValue(line line_info) int {
	count := getLineCount(line)

	if count > 0 {
		count = powInt(2, count-1)
	}

	return count
}

func getLineCount(line line_info) int {
	count := 0

	for _, n := range line.line_numbers {
		if contains(n, line.winning_numbers) {
			count = count + 1
		}
	}

	return count
}

func part1(filepath string) int {
	acc := 0
	c := make(chan string)
	go readlines(filepath, c)

	for line := range c {
		acc = acc + getLineValue(parseLine(line))
	}

	return acc
}

func part2(filepath string) int {
	acc := 0
	c := make(chan string)
	go readlines(filepath, c)

	remaining_games := []int{}
	infos := []line_info{}

	for line := range c {
		info := parseLine(line)
		count := getLineCount(info)
		for i := 1; i <= count; i = i + 1 {
			remaining_games = append(remaining_games, info.card_no+i)
		}
		acc = acc + 1
		infos = append(infos, info)
	}

	for i := 0; i < len(remaining_games); i = i + 1 {
		game_no := remaining_games[i]
		if game_no > len(infos) {
			continue
		}
		info := infos[game_no-1]
		count := getLineCount(info)
		for i := 1; i <= count; i = i + 1 {
			remaining_games = append(remaining_games, info.card_no+i)
		}
		acc = acc + 1
	}

	return acc
}

func main() {
	fmt.Println("Part 1: ", part1("data.txt"))
	fmt.Println("Part 2: ", part2("data.txt"))
}
