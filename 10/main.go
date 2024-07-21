package main

import (
	"bufio"
	"fmt"
	"os"
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

type pipe struct {
	n bool
	s bool
	e bool
	w bool
}

func classify(c rune) pipe {
	var p pipe

	switch c {
	case '|':
		{
			p.n = true
			p.s = true
		}
	case '-':
		{
			p.w = true
			p.e = true
		}
	case 'L':
		{
			p.n = true
			p.e = true
		}
	case 'J':
		{
			p.n = true
			p.w = true
		}
	case '7':
		{
			p.w = true
			p.s = true
		}
	case 'F':
		{
			p.s = true
			p.e = true
		}
	case 'S':
		{
			p.n = true
			p.s = true
			p.w = true
			p.e = true
		}
	}

	return p
}

func part1(filepath string) {
	lines := load(filepath)
	newLines := prune(lines)
	dfs(newLines)
}

func part2(filepath string) {
	lines := load(filepath)
	newLines := prune(lines)
	costs := dfs(newLines)
	printLines(lines)
	printCost(costs)
	exploded := explode(costs, lines)
	filled := fill(exploded)
	imploded := implode(filled)
	printCost(imploded)
	inside := calculateInterior(imploded)
	glyphs := getGlyphys(imploded, newLines)
	printGlyphs(glyphs)
	fmt.Println("inside: ", inside)
	printCost(filled)
}

func calculateInterior(costs [][]int) int {
	tiles := 0

	for _, l := range costs {
		for _, v := range l {
			if v == 0 {
				tiles++
			}
		}
	}

	return tiles
}

func load(filepath string) []string {
	lines := []string{}
	c := make(chan string)

	go readlines(filepath, c)

	for s := range c {
		lines = append(lines, s)
	}

	return lines
}

func getPipeRune(x, y int, lines []string) rune {
	if x < 0 || y < 0 || y >= len(lines) || x >= len(lines[y]) {
		return '.'
	}

	return rune(lines[y][x])
}

// y x
type pos struct {
	x int
	y int
}

func findStart(lines []string) pos {
	for y, l := range lines {
		for x := range l {
			if string(lines[y][x]) == "S" {
				return pos{x, y}
			}
		}
	}

	panic("Could not find start")
}

func dfs(lines []string) [][]int {
	costs := make([][]int, len(lines))

	for i := range costs {
		costs[i] = make([]int, len(lines[0]))
	}

	for y, l := range lines {
		for x := range l {
			costs[y][x] = 0
		}
	}

	currentCost := 1
	positions := []pos{}
	positions = append(positions, findStart(lines))

	for len(positions) > 0 {
		nextPositions := []pos{}

		addIfPipe := func(x, y int, w rune) {
			r := getPipeRune(x, y, lines)
			if r != '.' && costs[y][x] == 0 {
				other := classify(r)
				if w == 'w' && !other.w {
					return
				}
				if w == 'e' && !other.e {
					return
				}
				if w == 'n' && !other.n {
					return
				}
				if w == 's' && !other.s {
					return
				}
				nextPositions = append(nextPositions, pos{x, y})
			}
		}
		for _, p := range positions {
			costs[p.y][p.x] = currentCost
			pipe := classify(getPipeRune(p.x, p.y, lines))

			if pipe.w {
				addIfPipe(p.x-1, p.y, 'e')
			}
			if pipe.e {
				addIfPipe(p.x+1, p.y, 'w')
			}
			if pipe.n {
				addIfPipe(p.x, p.y-1, 's')
			}
			if pipe.s {
				addIfPipe(p.x, p.y+1, 'n')
			}
		}
		positions = nextPositions
		currentCost++
	}

	fmt.Printf("\nmax cost: %d\n", currentCost-2)

	return costs
}

func printLines(lines []string) {
	for _, l := range lines {
		fmt.Printf(l)
		fmt.Println()
	}
}

func printCost(cost [][]int) {
	fmt.Println()
	for y, l := range cost {
		fmt.Printf("%3d", y)
		for _, v := range l {
			var s string
			if v == 0 {
				s = fmt.Sprintf("%1s", ".")
			} else if v == -1 {
				s = fmt.Sprintf(" ")
			} else {
				s = fmt.Sprintf("%1s", "*")
			}
			fmt.Printf(s)
		}
		fmt.Println()
	}
}

func explode(costsIn [][]int, lines []string) [][]int {
	costsOut := make([][]int, len(costsIn)*2)

	for i := range costsOut {
		costsOut[i] = make([]int, len(costsIn[0])*2)
	}

	for y, l := range costsOut {
		for x := range l {
			costsOut[y][x] = 0
		}
	}

	for y := 0; y < len(costsIn); y++ {
		for x := 0; x < len(costsIn[y]); x++ {
			v := costsIn[y][x]
			costsOut[y*2][x*2] = v
			if y < len(costsIn)-1 && connectedNS(y, x, lines) {
				costsOut[y*2+1][x*2] = v
			}
			if x < len(costsIn[0])-1 && connectedWE(y, x, lines) {
				costsOut[y*2][x*2+1] = v
			}
		}
	}
	return costsOut
}

func implode(costsIn [][]int) [][]int {
	costsOut := make([][]int, len(costsIn)/2)

	for i := range costsOut {
		costsOut[i] = make([]int, len(costsIn[0])/2)
	}

	for y := 0; y < len(costsOut); y++ {
		for x := 0; x < len(costsOut[y]); x++ {
			costsOut[y][x] = costsIn[y*2][x*2]
		}
	}
	return costsOut
}

func getGlyphys(costsIn [][]int, lines []string) [][]rune {
	glyphs := make([][]rune, len(costsIn))

	for i := range glyphs {
		glyphs[i] = make([]rune, len(costsIn[0]))
	}

	for y, l := range glyphs {
		for x := range l {
			glyphs[y][x] = ' '
		}
	}

	for y := 0; y < len(glyphs); y++ {
		for x := 0; x < len(glyphs[y]); x++ {
			v := costsIn[y][x]

			if v == 0 {
				glyphs[y][x] = '.'
			} else if v > 0 {
				letter := lines[y][x]
				glyphs[y][x] = convLetter(rune(letter))
			}
		}
	}

	return glyphs
}

func convLetter(letter rune) rune {
	switch letter {
	case '|':
		return '│'
	case '-':
		return '─'
	case 'L':
		return '╰'
	case 'J':
		return '╯'
	case '7':
		return '╮'
	case 'F':
		return '╭'
	case 'S':
		return '*'
	default:
		return '?'
	}
}

/*
╭─[1]─Status───────────────────────────────────────────────────────────────────────────────────╮
│✓ aocgo2023 → master                                                                          │
╰──────────────────────────────────────────────────────────────────────────────────────────────╯
*/

func printGlyphs(glyphs [][]rune) {
	for y, l := range glyphs {
		fmt.Printf("%3d", y)
		fmt.Println(string(l))
	}
}

func fill(costsIn [][]int) [][]int {
	// costsOut has a border of 1 which we will remove before returning
	costsOut := make([][]int, len(costsIn)+2)

	for i := range costsOut {
		costsOut[i] = make([]int, len(costsIn[0])+2)
	}

	for y, l := range costsIn {
		for x := range l {
			costsOut[y+1][x+1] = costsIn[y][x]
		}
	}

	yMax := len(costsIn) - 1
	xMax := len(costsIn[0]) - 1

	// fill borders
	for y := 0; y <= yMax+2; y++ {
		costsOut[y][0] = -1
		costsOut[y][xMax+2] = -1
	}

	for x := 0; x <= xMax+2; x++ {
		costsOut[0][x] = -1
		costsOut[yMax+2][x] = -1
	}

	// bruteforce progressive fill
	changed := 1

	for changed > 0 {
		changed = 0

		for y := 1; y < yMax+2; y++ {
			for x := 1; x < xMax+2; x++ {
				if costsOut[y][x] == 0 && canFill(y, x, costsOut) {
					costsOut[y][x] = -1
					changed++
				}
			}
		}
	}

	costsOut = costsOut[1 : yMax+2]

	for y, l := range costsOut {
		costsOut[y] = l[1 : xMax+2]
	}

	return costsOut
}

func canFill(y, x int, costs [][]int) bool {
	if costs[y-1][x] == -1 {
		return true
	}

	if costs[y+1][x] == -1 {
		return true
	}

	if costs[y][x-1] == -1 {
		return true
	}

	if costs[y][x+1] == -1 {
		return true
	}

	if costs[y+1][x+1] == -1 {
		return true
	}

	if costs[y+1][x-1] == -1 {
		return true
	}

	if costs[y-1][x+1] == -1 {
		return true
	}

	if costs[y-1][x-1] == -1 {
		return true
	}

	return false
}

func connectedWE(y, x int, lines []string) bool {
	p1 := classify(getPipeRune(x, y, lines))
	p2 := classify(getPipeRune(x+1, y, lines))

	return p1.e && p2.w
}

func connectedNS(y, x int, lines []string) bool {
	p1 := classify(getPipeRune(x, y, lines))
	p2 := classify(getPipeRune(x, y+1, lines))

	return p1.s && p2.n
}

func survives(x, y int, lines []string) bool {
	p := classify(getPipeRune(x, y, lines))
	hits := 0

	if p.n {
		p2 := classify(getPipeRune(x, y-1, lines))
		if p2.s {
			hits++
		}
	}
	if p.s {
		p2 := classify(getPipeRune(x, y+1, lines))
		if p2.n {
			hits++
		}
	}
	if p.e {
		p2 := classify(getPipeRune(x+1, y, lines))
		if p2.w {
			hits++
		}
	}
	if p.w {
		p2 := classify(getPipeRune(x-1, y, lines))
		if p2.e {
			hits++
		}
	}

	return hits > 1
}

func prune(lines []string) []string {
	var newLines []string

	for {
		newLines = []string{}
		changed := 0

		for y, l := range lines {
			nl := ""
			for x := range l {
				ok := survives(x, y, lines)
				if ok {
					nl += string(lines[y][x])
				} else {
					nl += "."
					if string(lines[y][x]) != "." {
						changed++
					}
				}
			}
			newLines = append(newLines, nl)
		}

		if changed == 0 {
			break
		}
		lines = newLines
	}

	return newLines
}

func main() {
	part2("test.txt")
	part2("test2.txt")
	part2("test3.txt")
	part2("test4.txt")
	part2("test5.txt")
	part2("test6.txt")
	part2("data.txt")
}
