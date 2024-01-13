package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

func parseGame(input string) (int, []map[string]int) {
	game_re := regexp.MustCompile("^Game\\s*([0-9]+):\\s(.*)$")
	pair_re := regexp.MustCompile("^\\s*([0-9]+)\\s*(.*)\\s*$")

	matches := game_re.FindSubmatch([]byte(input))

	if len(matches) == 0 {
		panic(fmt.Sprintf("could not match %q", input))
	}

	game_no, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		panic(fmt.Sprintf("failed to parse game number %q", input))
	}

	sub_game_strings := strings.Split(string(matches[2]), ";")
	sub_games := []map[string]int{}
	for _, sg := range sub_game_strings {
		sub_game := map[string]int{}
		pairs := strings.Split(sg, ",")
		for _, p := range pairs {
			pair_match := pair_re.FindSubmatch([]byte(p))

			if len(pair_match) == 0 {
				panic(fmt.Sprintf("failed to parse pair %q", p))
			}

			pair_count, err := strconv.Atoi(string(pair_match[1]))
			if err != nil {
				panic(fmt.Sprintf("failed to parse pair number %q", p))
			}
			pair_colour := string(pair_match[2])
			sub_game[pair_colour] = pair_count
		}
		sub_games = append(sub_games, sub_game)
	}

	return game_no, sub_games
}

func testSubGame(input, threshold map[string]int) bool {
	for k, v := range threshold {
		if input[k] > v {
			return false
		}
	}
	return true
}

func testGame(input []map[string]int, threshold map[string]int) bool {
	for _, sg := range input {
		res := testSubGame(sg, threshold)
		if res == false {
			return false
		}
	}
	return true
}

func getGamePower(input []map[string]int) int {
	maxmap := map[string]int{}

	for _, sg := range input {
		for colour, val := range sg {
			if maxmap[colour] < val {
				maxmap[colour] = val
			}
		}
	}

	product := 1

	for _, val := range maxmap {
		product = product * val
	}

	return product
}

func part1(filepath string, threshold map[string]int) int {
	c := make(chan string)
	go readlines(filepath, c)

	acc := 0

	for line := range c {
		game_no, sub_games := parseGame(line)
		if testGame(sub_games, threshold) {
			acc = acc + game_no
		}
	}

	return acc
}

func part2(filepath string) int {
	c := make(chan string)
	go readlines(filepath, c)

	acc := 0

	for line := range c {
		_, sub_games := parseGame(line)
		pow := getGamePower(sub_games)
		acc = acc + pow
	}

	return acc
}

func main() {
	threshold := map[string]int{"red": 12, "green": 13, "blue": 14}
	part1_res := part1("./data.txt", threshold)
	fmt.Println("Part 1 result:", part1_res)
	part2_res := part2("./data.txt")
	fmt.Println("Part 2 result:", part2_res)
}
