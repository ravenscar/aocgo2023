package main

import (
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
	state_hand   = iota
	state_bet = iota
)

type bet_data struct {
  hand string
  bet int
}

func parse_file(filepath string) {
  c := make(chan string)
  go tokenize(filepath, c)
  state := state_hand
  var hand string
  hands := []bet_data{}

  for s := range c {
    switch state {
    case state_hand:
      hand = s
      state = state_bet
    case state_bet:
      bet, err := strconv.Atoi(s)
      if err != nil {
        panic(err)
      }
      hands = append(hands, bet_data{hand, bet})
    }
  }

}

