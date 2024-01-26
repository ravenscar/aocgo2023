package main

import (
	"fmt"
	"os"
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
	stzero   = iota
	stnaming = iota
	sttimes  = iota
	stdists  = iota
)

type race_data struct {
	time int
	dist int
}

func parse_file(filepath string) []race_data {
	times := []int{}
	dists := []int{}
	state := stzero
	name := ""

	c := make(chan string)
	go tokenize(filepath, c)

	for s := range c {
		switch state {
		case stzero:
			state = stnaming
			name = s
		case stnaming:
			if s == ":" {
				if name == "Time" {
					state = sttimes
				} else if name == "Distance" {
					state = stdists
				}
			} else {
				name = name + s
			}
		case sttimes:
			v, err := strconv.Atoi(s)
			if err != nil {
				state = stnaming
				name = s
			} else {
				times = append(times, v)
			}
		case stdists:
			v, err := strconv.Atoi(s)
			if err != nil {
				state = stnaming
				name = s
			} else {
				dists = append(dists, v)
			}
		}
	}

	if len(dists) != len(times) {
		panic("expected dists and times to have same length")
	}

	data := make([]race_data, len(dists))

	for i := range dists {
		data[i] = race_data{
			time: times[i],
			dist: dists[i],
		}
	}

	return data
}

func calculate_speeds(d race_data) int {
	var min, max int

	for i := 1; i < d.time; i++ {
		dis := i * (d.time - i)
		if dis > d.dist {
			min = i
			break
		}
	}

	for i := d.time - 1; i > 0; i-- {
		dis := i * (d.time - i)
		if dis > d.dist {
			max = i
			break
		}
	}

	if max < min {
		return 0
	}
	return (max - min) + 1
}

func multiply(vs []int) int {
	acc := 1

	for _, v := range vs {
		acc = acc * v
	}

	return acc
}

func part1(filepath string) int {
	data := parse_file(filepath)
	vals := make([]int, len(data))

	for i, d := range data {
		vals[i] = calculate_speeds(d)
	}

	return multiply(vals)
}

func combine_race_datas(ds []race_data) race_data {
	cdist := ""
	ctime := ""

	for _, d := range ds {
		cdist = cdist + strconv.Itoa(d.dist)
		ctime = ctime + strconv.Itoa(d.time)
	}

	dist, err := strconv.Atoi(cdist)
	if err != nil {
		panic(err)
	}

	time, err := strconv.Atoi(ctime)
	if err != nil {
		panic(err)
	}

	return race_data{time, dist}
}

func part2(filepath string) int {
	ds := parse_file(filepath)
	d := combine_race_datas(ds)
	return calculate_speeds(d)
}

func main() {
	fmt.Printf("part1: %d\n", part1("data.txt"))
	fmt.Printf("part2: %d\n", part2("data.txt"))
}
