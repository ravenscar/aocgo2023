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

func Test_convert_number(t *testing.T) {
	// 0 15 37
	// 37 52 2
	// 39 0 15
	convs := []almanac_conversion{
		{
			dest_range_start:   0,
			source_range_start: 15,
			range_length:       37,
		},
		{
			dest_range_start:   37,
			source_range_start: 52,
			range_length:       2,
		},
		{
			dest_range_start:   39,
			source_range_start: 0,
			range_length:       15,
		},
	}

	expect := make_expect(t)

	expect(39, convert_number(0, convs))
	expect(44, convert_number(5, convs))
	expect(53, convert_number(14, convs))
	expect(0, convert_number(15, convs))
	expect(10, convert_number(25, convs))
	expect(36, convert_number(51, convs))
	expect(37, convert_number(52, convs))
	expect(38, convert_number(53, convs))
	expect(54, convert_number(54, convs))
}

func Test_parse_file(t *testing.T) {
	seed_codes, almanac_entries, almanac_lookup := parse_file("test.txt")

	expect := make_expect(t)

	expect(4, len(seed_codes))
	expect(7, len(almanac_entries))
	expect(7, len(almanac_lookup))
	expect(3, len(almanac_lookup["light-to-temperature"]))
	// 45 77 23
	// 81 45 19
	// 68 64 13
	expect(45, almanac_lookup["light-to-temperature"][0].dest_range_start)
	expect(77, almanac_lookup["light-to-temperature"][0].source_range_start)
	expect(23, almanac_lookup["light-to-temperature"][0].range_length)
	expect(81, almanac_lookup["light-to-temperature"][1].dest_range_start)
	expect(45, almanac_lookup["light-to-temperature"][1].source_range_start)
	expect(19, almanac_lookup["light-to-temperature"][1].range_length)
	expect(68, almanac_lookup["light-to-temperature"][2].dest_range_start)
	expect(64, almanac_lookup["light-to-temperature"][2].source_range_start)
	expect(13, almanac_lookup["light-to-temperature"][2].range_length)
}

func Test_parse_number_mult(t *testing.T) {
	_, almanac_entries, almanac_lookup := parse_file("test.txt")

	expect := make_expect(t)

	expect(82, convert_number_mult(79, almanac_entries, almanac_lookup))
	expect(43, convert_number_mult(14, almanac_entries, almanac_lookup))
	expect(86, convert_number_mult(55, almanac_entries, almanac_lookup))
	expect(35, convert_number_mult(13, almanac_entries, almanac_lookup))
}

func Test_part1(t *testing.T) {
	v := part1("test.txt")
	expect := make_expect(t)

	expect(35, v)
}
