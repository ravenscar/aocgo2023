package main

import (
	"testing"
)

func make_expect(t *testing.T) func(i, j int) {
	return func(expected, i int) {
		if i != expected {
			t.Fatalf("Expected %d but got %d", expected, i)
		}
	}
}

func Test_parse_file(t *testing.T) {
	expect := make_expect(t)
	race_data := parse_file("test.txt")
	expect(3, len(race_data))
	expect(7, race_data[0].time)
	expect(9, race_data[0].dist)
	expect(15, race_data[1].time)
	expect(40, race_data[1].dist)
	expect(30, race_data[2].time)
	expect(200, race_data[2].dist)
}

func Test_calculate_speeds(t *testing.T) {
	expect := make_expect(t)
	expect(4, calculate_speeds(race_data{7, 9}))
	expect(8, calculate_speeds(race_data{15, 40}))
	expect(9, calculate_speeds(race_data{30, 200}))
}

func Test_multiply(t *testing.T) {
	expect := make_expect(t)
	expect(288, multiply([]int{4, 8, 9}))
}

func Test_part1(t *testing.T) {
	expect := make_expect(t)
	expect(288, part1("test.txt"))
}

func Test_part2(t *testing.T) {
	expect := make_expect(t)
	expect(71503, part2("test.txt"))
}
