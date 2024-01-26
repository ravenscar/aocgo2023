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

func Test_load(t *testing.T) {
	expect := make_expect(t)
	data := loadAll("test1.txt")
	expect(2, len(data.directions))
	expect("R", data.directions[0])
	expect("L", data.directions[1])
	expect("BBB", data.lookup["AAA"].L)
	expect("CCC", data.lookup["AAA"].R)
	expect("ZZZ", data.lookup["ZZZ"].L)
	expect("ZZZ", data.lookup["ZZZ"].R)
}

func Test_part1(t *testing.T) {
	expect := make_expect(t)
	expect(2, part1("test1.txt"))
	expect(6, part1("test2.txt"))
	expect(6, part2("test3.txt"))
	expect(6, part2lcm("test3.txt"))
}
