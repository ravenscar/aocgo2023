package main

import (
	"testing"
)

func make_expect(t *testing.T) func(i, j interface{}) {
	return func(expected, i interface{}) {
		if i != expected {
			switch i.(type) {
			case string:
				t.Fatalf("Expected %s but got %s", expected, i)
			case int:
				t.Fatalf("Expected %d but got %d", expected, i)
			case bool:
				t.Fatalf("Expected %b but got %b", expected, i)
			default:
				t.Fatalf("Expected %v but got %v", expected, i)
			}
		}
	}
}

func Test_parse_file(t *testing.T) {
	expect := make_expect(t)
	bets := parse_file("test.txt")
	expect(5, len(bets))

	expect("32T3K", bets[0].hand)
	expect(765, bets[0].bet)
	expect("T55J5", bets[1].hand)
	expect(684, bets[1].bet)
	expect("KK677", bets[2].hand)
	expect(28, bets[2].bet)
	expect("KTJJT", bets[3].hand)
	expect(220, bets[3].bet)
	expect("QQQJA", bets[4].hand)
	expect(483, bets[4].bet)
}

func Test_process_hand(t *testing.T) {
	expect := make_expect(t)
	lookup := build_hand_lookup("9TJJT")

	expect(0, lookup["A"])
	expect(1, lookup["9"])
	expect(2, lookup["T"])
	expect(2, lookup["J"])
}

func Test_get_hand_value(t *testing.T) {
	expect := make_expect(t)

	expect(five_oaK, get_hand_value(bet_data{lookup: build_hand_lookup("44444")}, '*'))
	expect(four_oaK, get_hand_value(bet_data{lookup: build_hand_lookup("44544")}, '*'))
	expect(full_house, get_hand_value(bet_data{lookup: build_hand_lookup("44545")}, '*'))
	expect(three_oaK, get_hand_value(bet_data{lookup: build_hand_lookup("44543")}, '*'))
	expect(two_pair, get_hand_value(bet_data{lookup: build_hand_lookup("9TJJT")}, '*'))
	expect(pair, get_hand_value(bet_data{lookup: build_hand_lookup("9TJJK")}, '*'))
	expect(high_card, get_hand_value(bet_data{lookup: build_hand_lookup("9TJ3K")}, '*'))

	expect(pair, get_hand_value(bet_data{lookup: build_hand_lookup("9TJ3K")}, 'J'))
	expect(three_oaK, get_hand_value(bet_data{lookup: build_hand_lookup("9JJ3K")}, 'J'))
	expect(four_oaK, get_hand_value(bet_data{lookup: build_hand_lookup("9JJJK")}, 'J'))
	expect(five_oaK, get_hand_value(bet_data{lookup: build_hand_lookup("9JJJJ")}, 'J'))
}

func Test_compare_hand_position(t *testing.T) {
	expect := make_expect(t)
	expect(-1, compare_hand_position("2", "3", card_rank))
	expect(1, compare_hand_position("3", "2", card_rank))
	expect(9, compare_hand_position("A", "5", card_rank))
	expect(-7, compare_hand_position("6", "K", card_rank))
	expect(-7, compare_hand_position("A6", "AK", card_rank))
	expect(-7, compare_hand_position("AAAAAAAA6", "AAAAAAAAK", card_rank))
	expect(3, compare_hand_position("6", "3", card_rank))
	expect(3, compare_hand_position("64444", "35555", card_rank))
	expect(8, compare_hand_position("J4444", "35555", card_rank))
	expect(-2, compare_hand_position("J4444", "35555", card_rank2))
}

func Test_compare_hands(t *testing.T) {
	expect := make_expect(t)
	makeHand := func(hand string) bet_data {
		lookup := build_hand_lookup(hand)
		bet := 0
		return bet_data{hand, bet, lookup}
	}
	compare_hand := make_compare_hand(card_rank, '*')
	expect(1, compare_hand(makeHand("44444"), makeHand("44445")))
	expect(-1, compare_hand(makeHand("44445"), makeHand("44444")))
	expect(-1, compare_hand(makeHand("44444"), makeHand("55555")))
	expect(1, compare_hand(makeHand("55555"), makeHand("44444")))
	expect(3, compare_hand(makeHand("64444"), makeHand("35555")))
}

func Test_part1(t *testing.T) {
	expect := make_expect(t)
	v := part1("test.txt")
	expect(6440, v)
}

func Test_part2(t *testing.T) {
	expect := make_expect(t)
	v := part2("test.txt")
	expect(5905, v)
}
