package main

import (
	"fmt"
)

type mazeData struct {
	data  []string
	costs [][]int
}

type pipe rune

type connections struct {
	n bool
	s bool
	e bool
	w bool
}

func (maze mazeData) classify(y, x int) connections {
	var conns connections
	p := maze.getPipe(y, x)

	switch p {
	case '|':
		{
			conns.n = true
			conns.s = true
		}
	case '-':
		{
			conns.w = true
			conns.e = true
		}
	case 'L':
		{
			conns.n = true
			conns.e = true
		}
	case 'J':
		{
			conns.n = true
			conns.w = true
		}
	case '7':
		{
			conns.w = true
			conns.s = true
		}
	case 'F':
		{
			conns.s = true
			conns.e = true
		}
	case 'S':
		{
			conns.n = true
			conns.s = true
			conns.w = true
			conns.e = true
		}
	}

	return conns
}

func part1(filepath string) {
	var maze mazeData

	maze.load(filepath)
	maze.prune()
	maze.dfs_costs()
}

func part2(filepath string) {
	var maze mazeData

	fmt.Println(filepath)
	maze.load(filepath)
	fmt.Println("Raw:")
	maze.printLines()
	fmt.Println("Loaded:")
	maze.printGlyphs()
	maze.prune()
	fmt.Println("Pruned:")
	maze.printGlyphs()
	maze.clearOutside()
	fmt.Println("Cleared:")
	maze.printGlyphs()
	inside := maze.calculateInterior()
	fmt.Println("Inside count: ", inside)
	maze.dfs_costs()
}

func (maze mazeData) calculateInterior() int {
	tiles := 0

	for _, l := range maze.data {
		for _, v := range l {
			if v == '.' {
				tiles++
			}
		}
	}

	return tiles
}

func (maze *mazeData) load(filepath string) {
	lines := []string{}
	c := make(chan string)

	go readlines(filepath, c)

	for s := range c {
		lines = append(lines, s)
	}

	maze.data = lines
}

func (maze mazeData) getPipe(y, x int) pipe {
	if x < 0 || y < 0 || y >= len(maze.data) || x >= len(maze.data[y]) {
		return '.'
	}

	return pipe(maze.data[y][x])
}

// y x
type pos struct {
	y int
	x int
}

func (maze mazeData) findStart() pos {
	for y, l := range maze.data {
		for x := range l {
			if string(l[x]) == "S" {
				return pos{y, x}
			}
		}
	}

	maze.printLines()
	panic("Could not find start")
}

func (maze *mazeData) dfs_costs() {
	costs := make([][]int, len(maze.data))

	for i := range costs {
		costs[i] = make([]int, len(maze.data[0]))
	}

	for y, l := range maze.data {
		for x := range l {
			costs[y][x] = 0
		}
	}

	currentCost := 1
	positions := []pos{}
	positions = append(positions, maze.findStart())

	for len(positions) > 0 {
		nextPositions := []pos{}

		addIfPipe := func(y, x int, w rune) {
			p := maze.getPipe(y, x)
			if p != '.' && costs[y][x] == 0 {
				other := maze.classify(y, x)
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
				nextPositions = append(nextPositions, pos{y, x})
			}
		}
		for _, p := range positions {
			costs[p.y][p.x] = currentCost
			pipe := maze.classify(p.y, p.x)

			if pipe.w {
				addIfPipe(p.y, p.x-1, 'e')
			}
			if pipe.e {
				addIfPipe(p.y, p.x+1, 'w')
			}
			if pipe.n {
				addIfPipe(p.y-1, p.x, 's')
			}
			if pipe.s {
				addIfPipe(p.y+1, p.x, 'n')
			}
		}
		positions = nextPositions
		currentCost++
	}

	fmt.Printf("\nmax cost: %d\n", currentCost-2)

	maze.costs = costs
}

func (maze mazeData) printLines() {
	for x, l := range maze.data {
		fmt.Print(x)
		fmt.Printf(l)
		fmt.Println()
	}
}

func (maze mazeData) printCost() {
	fmt.Println()
	for y, l := range maze.costs {
		fmt.Printf("%3d", y)
		for _, v := range l {
			var s string
			if v == 0 {
				s = fmt.Sprintf("%3s", ".")
			} else if v == -1 {
				s = fmt.Sprintf(" ")
			} else {
				s = fmt.Sprintf("%3d", v-1)
			}
			fmt.Printf(s)
		}
		fmt.Println()
	}
}

func explode(maze mazeData) [][]pipe {
	exploded := make([][]pipe, len(maze.data)*2)

	for i := range exploded {
		exploded[i] = make([]pipe, len(maze.data[0])*2)
	}

	for y, l := range exploded {
		for x := range l {
			exploded[y][x] = '.'
		}
	}

	for y := 0; y < len(maze.data); y++ {
		for x := 0; x < len(maze.data[y]); x++ {
			v := maze.getPipe(y, x)
			exploded[y*2][x*2] = v
			if y < len(maze.data)-1 && maze.connectedNS(y, x) {
				exploded[y*2+1][x*2] = '|'
			}
			if x < len(maze.data[0])-1 && maze.connectedWE(y, x) {
				exploded[y*2][x*2+1] = '-'
			}
		}
	}
	return exploded
}

func implode(costsIn [][]pipe) [][]pipe {
	costsOut := make([][]pipe, len(costsIn)/2)

	for i := range costsOut {
		costsOut[i] = make([]pipe, len(costsIn[0])/2)
	}

	for y := 0; y < len(costsOut); y++ {
		for x := 0; x < len(costsOut[y]); x++ {
			costsOut[y][x] = costsIn[y*2][x*2]
		}
	}
	return costsOut
}

func (maze mazeData) printGlyphs() {
	glyphs := make([][]rune, len(maze.data))

	for y := range glyphs {
		glyphs[y] = make([]rune, len(maze.data[y]))
		for x := range glyphs[y] {
			v := rune(maze.data[y][x])
			glyphs[y][x] = convLetter(v)
		}
	}

	for _, l := range glyphs {
		fmt.Println(string(l))
	}
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
		return letter
	}
}

func (maze *mazeData) clearOutside() {
	exploded := explode(*maze)
	// filled has a border of 1 which we will remove before returning
	filled := make([][]pipe, len(exploded)+2)

	for i := range filled {
		filled[i] = make([]pipe, len(exploded[0])+2)
	}

	for y, l := range exploded {
		for x := range l {
			filled[y+1][x+1] = exploded[y][x]
		}
	}

	yMax := len(exploded) - 1
	xMax := len(exploded[0]) - 1

	// fill borders
	for y := 0; y <= yMax+2; y++ {
		filled[y][0] = ' '
		filled[y][xMax+2] = ' '
	}

	for x := 0; x <= xMax+2; x++ {
		filled[0][x] = ' '
		filled[yMax+2][x] = ' '
	}

	// bruteforce progressive fill
	changed := 1

	for changed > 0 {
		changed = 0

		for y := 1; y < yMax+2; y++ {
			for x := 1; x < xMax+2; x++ {
				if filled[y][x] == '.' && canFill(y, x, filled) {
					filled[y][x] = ' '
					changed++
				}
			}
		}
	}

	filled = filled[1 : yMax+2]

	for y, l := range filled {
		filled[y] = l[1 : xMax+2]
	}

	imploded := implode(filled)

	lines := make([]string, len(imploded))

	for x, pipes := range imploded {
		lines[x] = string(pipes)
	}

	maze.data = lines
}

func canFill(y, x int, costs [][]pipe) bool {
	if costs[y-1][x] == ' ' {
		return true
	}

	if costs[y+1][x] == ' ' {
		return true
	}

	if costs[y][x-1] == ' ' {
		return true
	}

	if costs[y][x+1] == ' ' {
		return true
	}

	if costs[y+1][x+1] == ' ' {
		return true
	}

	if costs[y+1][x-1] == ' ' {
		return true
	}

	if costs[y-1][x+1] == ' ' {
		return true
	}

	if costs[y-1][x-1] == ' ' {
		return true
	}

	return false
}

func (maze mazeData) connectedWE(y, x int) bool {
	p1 := maze.classify(y, x)
	p2 := maze.classify(y, x+1)

	return p1.e && p2.w
}

func (maze mazeData) connectedNS(y, x int) bool {
	p1 := maze.classify(y, x)
	p2 := maze.classify(y+1, x)

	return p1.s && p2.n
}

func (maze mazeData) survives(y, x int) bool {
	p := maze.classify(y, x)
	hits := 0

	if p.n {
		p2 := maze.classify(y-1, x)
		if p2.s {
			hits++
		}
	}
	if p.s {
		p2 := maze.classify(y+1, x)
		if p2.n {
			hits++
		}
	}
	if p.e {
		p2 := maze.classify(y, x+1)
		if p2.w {
			hits++
		}
	}
	if p.w {
		p2 := maze.classify(y, x-1)
		if p2.e {
			hits++
		}
	}

	return hits > 1
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func (maze *mazeData) prune() {
	for {
		changed := 0

		for y, l := range maze.data {
			for x := range l {
				survives := maze.survives(y, x)
				if !survives && maze.data[y][x] != '.' {
					maze.data[y] = replaceAtIndex(maze.data[y], '.', x)
					changed++
				}
			}
		}

		if changed == 0 {
			break
		}
	}
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
