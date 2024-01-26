package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

type bet_data struct {
	hand   string
	bet    int
	lookup map[string]int
}

func parse_file(filepath string) []bet_data {
	c := make(chan string)
	go readlines(filepath, c)
	hands := []bet_data{}

	for s := range c {
		if len(strings.Trim(s, " ")) == 0 {
			continue
		}

		split := strings.Split(s, " ")
		if len(split) != 2 {
			panic(fmt.Sprintf("expected two strings: %s", s))
		}
		hand := split[0]
		bet, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}

		lookup := build_hand_lookup(hand)

		hands = append(hands, bet_data{hand, bet, lookup})
	}

	return hands
}

const (
	high_card  = iota
	pair       = iota
	two_pair   = iota
	three_oaK  = iota
	full_house = iota
	four_oaK   = iota
	five_oaK   = iota
)

func build_hand_lookup(h string) map[string]int {
	cards := strings.Split(h, "")
	lookup := map[string]int{}

	for _, c := range cards {
		lookup[c] = lookup[c] + 1
	}

	return lookup
}

func get_hand_value(h bet_data, wild rune) int {
	inv_lookup := map[int]int{}

	wild_count := h.lookup[string(wild)]

	for k, v := range h.lookup {
		if k != string(wild) {
			inv_lookup[v] = inv_lookup[v] + 1
		}
	}

	if inv_lookup[5] > 0 {
		return five_oaK
	}

	if inv_lookup[4] > 0 {
		if wild_count > 0 {
			return five_oaK
		}
		return four_oaK
	}

	if inv_lookup[3] > 0 && inv_lookup[2] > 0 {
		return full_house
	}

	if inv_lookup[3] > 0 {
		if wild_count > 1 {
			return five_oaK
		}
		if wild_count > 0 {
			return four_oaK
		}
		return three_oaK
	}

	if inv_lookup[2] > 1 {
		if wild_count > 0 {
			return full_house
		}
		return two_pair
	}

	if inv_lookup[2] > 0 {
		if wild_count > 2 {
			return five_oaK
		}
		if wild_count > 1 {
			return four_oaK
		}
		if wild_count > 0 {
			return three_oaK
		}
		return pair
	}

	if wild_count > 3 {
		return five_oaK
	}
	if wild_count > 2 {
		return four_oaK
	}
	if wild_count > 1 {
		return three_oaK
	}
	if wild_count > 0 {
		return pair
	}
	return high_card
}

const (
	card_rank  = "23456789TJQKA"
	card_rank2 = "J23456789TQKA"
)

func compare_hand_position(h1, h2, rank string) int {
	i1 := 0
	i2 := 0
	i := 0

	for i1 == i2 {
		if i >= len(h1) {
			break
		}
		i1 = strings.Index(rank, string(h1[i]))
		if i1 == -1 {
			panic("not a card")
		}
		i2 = strings.Index(rank, string(h2[i]))
		if i2 == -1 {
			panic("not a card")
		}
		i = i + 1
	}

	return i1 - i2
}

func make_compare_hand(rank string, wild rune) func(bet_data, bet_data) int {
	comp_func := func(h1, h2 bet_data) int {
		i1 := get_hand_value(h1, wild)
		i2 := get_hand_value(h2, wild)

		if i1 == i2 {
			return compare_hand_position(h1.hand, h2.hand, rank)
		}

		return i1 - i2
	}
	return comp_func
}

func part1(filepath string) int {
	acc := 0
	bets := parse_file(filepath)
	slices.SortFunc(bets, make_compare_hand(card_rank, '*'))

	for i, b := range bets {
		v := (i + 1) * b.bet
		acc = acc + v
	}

	return acc
}

func part2(filepath string) int {
	acc := 0
	bets := parse_file(filepath)
	slices.SortFunc(bets, make_compare_hand(card_rank2, 'J'))

	for i, b := range bets {
		v := (i + 1) * b.bet
		acc = acc + v
	}

	return acc
}

func main() {
	fmt.Printf("part1: %d\n", part1("data.txt"))
	fmt.Printf("part2: %d\n", part2("data.txt"))
}
