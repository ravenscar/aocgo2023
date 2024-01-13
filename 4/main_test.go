package main

import (
	"testing"
)

func TestContains(t *testing.T) {
	arr := []int{41, 48, 83, 86, 17}

	if !contains(41, arr) {
		t.Fatalf("Expected to contain %d", 41)
	}

	if contains(42, arr) {
		t.Fatalf("Expected to not contain %d", 42)
	}
}

func TestParseLine(t *testing.T) {
	line_info := parseLine("Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53")
	if line_info.card_no != 1 {
		t.Fatalf("Bad card %v", line_info)
	}

	expected_winning := []int{41, 48, 83, 86, 17}
	expected_line := []int{83, 86, 6, 31, 17, 9, 48, 53}

	for _, i := range expected_winning {
		if !contains(i, line_info.winning_numbers) {
			t.Fatalf("Bad card %v expected winning %d", line_info, i)
		}
	}

	for _, i := range expected_line {
		if !contains(i, line_info.line_numbers) {
			t.Fatalf("Bad card %v expected line %d", line_info, i)
		}
	}

	if line_info.card_no != 1 {
		t.Fatalf("Bad card %v", line_info)
	}
}

func TestGetLineValue(t *testing.T) {
	line_info := parseLine("Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53")
	val := getLineValue(line_info)
	if val != 8 {
		t.Fatalf("expoected %d got %d", 8, val)
	}
}

func TestPart1(t *testing.T) {
	val := part1("./test.txt")

	if val != 13 {
		t.Fatalf("expoected %d got %d", 13, val)
	}
}

func TestPart2(t *testing.T) {
	val := part2("./test.txt")

	if val != 30 {
		t.Fatalf("expoected %d got %d", 30, val)
	}
}
