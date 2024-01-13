package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func loadAll(path string) []string {
	strings := []string{}
	c := make(chan string)

	go readlines(path, c)

	for line := range c {
		strings = append(strings, line)
	}
	scanForNumbers(strings)
	return strings
}

func isNumberByte(b byte) bool {
	if b >= 48 && b <= 57 {
		return true
	}
	return false
}

func isGear(b byte) bool {
	return b == 42
}

type number_match struct {
	value     string
	start_idx int
	end_idx   int
	line      int
	symbols   []byte
}

func scanForNumbers(lines []string) []number_match {
	w := len(lines[0])
	h := len(lines)

	matches := []number_match{}
	match_start := -1
	match_end := -1
	matching := false
	match_content := ""

	end_matching := func(y int) {
		if matching {
			match := number_match{
				value:     match_content,
				start_idx: match_start,
				end_idx:   match_end,
				line:      y,
			}

			symbols := getSurroundingSymbols(lines, match)
			match.symbols = symbols

			matches = append(matches, match)
			matching = false
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if isNumberByte(lines[y][x]) {
				if !matching {
					match_content = string(lines[y][x])
					match_start = x
					match_end = x
					matching = true
				} else {
					match_content = match_content + string(lines[y][x])
					match_end = x
				}
			} else {
				end_matching(y)
			}
		}
		end_matching(y)
	}

	// fmt.Println(matches)
	return matches
}

func getSymbol(b byte) *byte {
	if isNumberByte(b) || b == 46 {
		return nil
	}
	return &b
}

func getSurroundingSymbols(lines []string, match number_match) []byte {
	w := len(lines[0])
	h := len(lines)

	x1 := max(match.start_idx-1, 0)
	x2 := min(match.end_idx+1, w-1)
	y1 := max(match.line-1, 0)
	y2 := min(match.line+1, h-1)

	found := []byte{}

	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			b := getSymbol(lines[y][x])
			if b != nil {
				found = append(found, *b)
			}
		}
	}
	return found
}

type xy struct {
	x int
	y int
}

func getGearPositions(lines []string) []xy {
	w := len(lines[0])
	h := len(lines)
	positions := []xy{}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if isGear(lines[y][x]) {
				positions = append(positions, xy{x: x, y: y})
			}
		}
	}

	return positions
}

func getAdjacent(gear xy, matches []number_match) []number_match {
	adjacent := []number_match{}

	for _, m := range matches {
		if gear.x >= m.start_idx-1 && gear.x <= m.end_idx+1 {
			if gear.y >= m.line-1 && gear.y <= m.line+1 {
				adjacent = append(adjacent, m)
			}
		}
	}

	return adjacent
}

func part1(filepath string) int {
	acc := 0

	lines := loadAll(filepath)
	matches := scanForNumbers(lines)

	for _, m := range matches {
		if len(m.symbols) > 0 {
			ival, err := strconv.Atoi(m.value)
			if err != nil {
				panic(err)
			}
			acc = acc + ival
		}
	}

	return acc
}

func part2(filepath string) int {
	acc := 0

	lines := loadAll(filepath)
	matches := scanForNumbers(lines)
	gears := getGearPositions(lines)

	for _, g := range gears {
		adjacent := getAdjacent(g, matches)
		if len(adjacent) == 2 {
			v1, err := strconv.Atoi(adjacent[0].value)
			if err != nil {
				panic(err)
			}
			v2, err := strconv.Atoi(adjacent[1].value)
			if err != nil {
				panic(err)
			}
			acc = acc + v1*v2
		}
	}

	return acc
}

func main() {
	fmt.Println("Part 1: ", part1("./data.txt"))
	fmt.Println("Part 2: ", part2("./data.txt"))
}
