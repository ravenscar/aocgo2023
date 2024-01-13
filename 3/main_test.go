package main

import (
	"testing"
)

func TestLoadAll(t *testing.T) {
	lines := loadAll("./test.txt")

	if len(lines) != 10 {
		t.Fatalf("Expected %d lines but found %d", 10, len(lines))
	}

	if lines[0] != "467..114.." {
		t.Fatalf("bad first line %q", lines[0])
	}

	if lines[9] != ".664.598.." {
		t.Fatalf("bad last line %q", lines[9])
	}
}

func TestScanForNumbers(t *testing.T) {
	lines := loadAll("./test.txt")
	matches := scanForNumbers(lines)

	if len(matches) != 10 {
		t.Fatalf("Expected %d matches but found %d", 10, len(matches))
	}
}

func TestPart1(t *testing.T) {
	val := part1("./test.txt")

	if val != 4361 {
		t.Fatalf("Expected %d matches but found %d", 4361, val)
	}
}

func TestPart2(t *testing.T) {
	val := part2("./test.txt")

	if val != 4361 {
		t.Fatalf("Expected %d matches but found %d", 4361, val)
	}
}
