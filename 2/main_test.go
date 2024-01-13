package main

import (
	"testing"
)

func TestReadlines(t *testing.T) {
	c := make(chan string)
	go readlines("./test.txt", c)

	firstline := <-c
	if firstline != "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green" {
		t.Fatalf("bad first line %q", firstline)
	}
	counter := 1
	lastline := ""
	for line := range c {
		counter = counter + 1
		lastline = line
	}

	if lastline != "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green" {
		t.Fatalf("bad last line %q", firstline)
	}

	if counter != 5 {
		t.Fatalf("Expected %d lines but found %d", 5, counter)
	}
}

func TestParseGame(t *testing.T) {
	game_no, sub_games := parseGame("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green")
	if game_no != 1 {
		t.Fatalf("Expected game number of %d but found %d", 1, game_no)
	}

	if len(sub_games) != 3 {
		t.Fatalf("Expected %d sub games but found %d", 3, len(sub_games))
	}

	test_sub_game := func(idx int, game_counts map[string]int) {
		length := len(game_counts)
		if len(sub_games[idx]) != length {
			t.Fatalf("Expected sub game %d len to be %d but found %d", idx, length, len(sub_games[idx]))
		}

		for k, v := range game_counts {
			if sub_games[idx][k] != v {
				t.Fatalf("Expected sub game %d %q count to be %d but found %d", idx, k, v, sub_games[idx][k])
			}
		}
	}

	sub_game_counts := []map[string]int{
		{"blue": 3, "red": 4},
		{"red": 1, "green": 2, "blue": 6},
		{"green": 2},
	}

	test_sub_game(0, sub_game_counts[0])
	test_sub_game(1, sub_game_counts[1])
	test_sub_game(2, sub_game_counts[2])
}

func TestTestSubGame(t *testing.T) {
	sub_game_counts := map[string]int{"red": 1, "green": 2, "blue": 6}

	run_test := func(input, threshold map[string]int, expected bool) {
		res := testSubGame(input, threshold)

		if res != expected {
			t.Fatalf("Expected test to be %t but it was %t based on %v vs %v", expected, res, input, threshold)
		}
	}

	run_test(sub_game_counts, sub_game_counts, true)
	run_test(sub_game_counts, map[string]int{"red": 100}, true)
	run_test(sub_game_counts, map[string]int{"red": 1}, true)
	run_test(sub_game_counts, map[string]int{"blue": 7, "red": 1, "green": 1}, false)
	run_test(sub_game_counts, map[string]int{"purple": 10}, true)
	run_test(sub_game_counts, map[string]int{"purple": 0}, true)
}

func TestTestGame(t *testing.T) {
	sub_game_counts := []map[string]int{
		{"blue": 3, "red": 4},
		{"red": 1, "green": 2, "blue": 6},
		{"green": 2},
	}

	run_test := func(input []map[string]int, threshold map[string]int, expected bool) {
		res := testGame(input, threshold)

		if res != expected {
			t.Fatalf("Expected test to be %t but it was %t based on %v vs %v", expected, res, input, threshold)
		}
	}

	run_test(sub_game_counts, map[string]int{"red": 100}, true)
	run_test(sub_game_counts, map[string]int{"red": 1}, false)
	run_test(sub_game_counts, map[string]int{"blue": 7, "red": 1, "green": 1}, false)
	run_test(sub_game_counts, map[string]int{"purple": 10}, true)
	run_test(sub_game_counts, map[string]int{"purple": 0}, true)
}

func TestPart1(t *testing.T) {
	res := part1("./test.txt", map[string]int{"red": 12, "green": 13, "blue": 14})

	if res != 8 {
		t.Fatalf("Expected %d but it was %d", 8, res)
	}
}

func TestGamePower(t *testing.T) {
	sub_game_counts := []map[string]int{
		{"blue": 3, "red": 4},
		{"red": 1, "green": 2, "blue": 6},
		{"green": 2},
	}
	pow := getGamePower(sub_game_counts)

	if pow != 48 {
		t.Fatalf("Expected %d but it was %d", 48, pow)
	}
}

func TestPart2(t *testing.T) {
	res := part2("./test.txt")

	if res != 2286 {
		t.Fatalf("Expected %d but it was %d", 2286, res)
	}
}
