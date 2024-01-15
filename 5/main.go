package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"text/scanner"
)

func tokenize(filepath string, c chan string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	var s scanner.Scanner
	s.Init(file)
	s.Filename = "example"
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		c <- s.TokenText()
	}

	close(c)
}

const (
	zero   = iota
	naming = iota
	seeds  = iota
	drs    = iota
	srs    = iota
	rl     = iota
)

type almanac_conversion struct {
	dest_range_start   int
	source_range_start int
	range_length       int
}

func parse_file(filepath string) ([]int, []string, map[string][]almanac_conversion) {
	seed_codes := []int{}
	var map_name string
	var dest_range_start int
	var source_range_start int
	var range_length int
	state := zero
	almanac_entries := []string{}
	almanac_lookup := map[string][]almanac_conversion{}

	c := make(chan string)
	go tokenize(filepath, c)

	for s := range c {
		switch state {
		case zero:
			state = naming
			map_name = s
		case seeds:
			v, err := strconv.Atoi(s)
			if err != nil {
				state = naming
				map_name = s
			} else {
				seed_codes = append(seed_codes, v)
			}
		case naming:
			if s == ":" {
				if map_name == "seeds" {
					state = seeds
				} else {
					almanac_entries = append(almanac_entries, map_name)
					state = drs
				}
			} else {
				if s != "map" {
					map_name = map_name + s
				}
			}
		case drs:
			v, err := strconv.Atoi(s)
			if err != nil {
				state = naming
				map_name = s
			} else {
				dest_range_start = v
				state = srs
			}
		case srs:
			v, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			} else {
				source_range_start = v
				state = rl
			}
		case rl:
			v, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			} else {
				range_length = v
				state = drs
				almanac_lookup[map_name] = append(almanac_lookup[map_name], almanac_conversion{
					source_range_start: source_range_start,
					dest_range_start:   dest_range_start,
					range_length:       range_length,
				})
			}
		}
	}

	return seed_codes, almanac_entries, almanac_lookup
}

func convert_number(n int, convs []almanac_conversion) int {
	for _, conv := range convs {
		if n >= conv.source_range_start && n < conv.source_range_start+conv.range_length {
			delta := conv.dest_range_start - conv.source_range_start
			return n + delta
		}
	}

	return n
}

func convert_number_mult(n int, order []string, lookup map[string][]almanac_conversion) int {
	for _, key := range order {
		n = convert_number(n, lookup[key])
	}

	return n
}

func part1(filepath string) int {
	seed_codes, almanac_entries, almanac_lookup := parse_file(filepath)
	converted := []int{}

	for _, code := range seed_codes {
		converted = append(converted, convert_number_mult(code, almanac_entries, almanac_lookup))
	}

	return slices.Min(converted)
}

func part2(filepath string) int {
	seed_ranges, almanac_entries, almanac_lookup := parse_file(filepath)

	min_found := convert_number_mult(seed_ranges[0], almanac_entries, almanac_lookup)

	for i := range seed_ranges {
		if i%2 == 0 {
			for j := seed_ranges[i]; j < seed_ranges[i]+seed_ranges[i+1]; j++ {
				converted := convert_number_mult(j, almanac_entries, almanac_lookup)
				if converted < min_found {
					min_found = converted
				}
			}
		}
	}

	return min_found
}

func main() {
	fmt.Printf("Part 1: %d\n", part1("data.txt"))
	fmt.Printf("Part 2: %d\n", part2("data.txt"))
}
