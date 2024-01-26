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
	points := get_points("1 2 3 4 34 23456 0")
	expect(7, len(points))
	expect(1, points[0])
	expect(34, points[4])
	expect(23456, points[5])
	expect(0, points[6])
}

func Test_iterate(t *testing.T) {
	expect := make_expect(t)
	points1 := get_points("0 3 6 9 12 15")
	points2 := iterate_points(points1)

	expect(len(points1)-1, len(points2))
	expect(3, points2[0])
	expect(3, points2[1])
	expect(3, points2[2])
	expect(3, points2[3])
	expect(3, points2[4])

	points3 := iterate_points(points2)

	expect(len(points2)-1, len(points3))
	expect(0, points3[0])
	expect(0, points3[1])
	expect(0, points3[2])
	expect(0, points3[3])
}

func Test_end_state(t *testing.T) {
	expect := make_expect(t)
	points1 := get_points("0 3 6 9 12 15")
	points2 := iterate_points(points1)
	points3 := iterate_points(points2)

	expect(false, end_state(points1))
	expect(false, end_state(points2))
	expect(true, end_state(points3))
}

func Test_get_point_arrays(t *testing.T) {
	expect := make_expect(t)
	arrs := get_point_arrays("0 3 6 9 12 15")

	expect(3, len(arrs))
	expect(6, len(arrs[0]))
	expect(5, len(arrs[1]))
	expect(4, len(arrs[2]))

	expect([6]int{0, 3, 6, 9, 12, 15}, [6]int(arrs[0]))
	expect([5]int{3, 3, 3, 3, 3}, [5]int(arrs[1]))
	expect([4]int{0, 0, 0, 0}, [4]int(arrs[2]))
}

func Test_get_next_point(t *testing.T) {
	expect := make_expect(t)
	next1 := get_next_point(get_point_arrays("0 3 6 9 12 15"))
	next2 := get_next_point(get_point_arrays("1 3 6 10 15 21"))
	next3 := get_next_point(get_point_arrays("10 13 16 21 30 45"))

	expect(18, next1)
	expect(28, next2)
	expect(68, next3)
}

func Test_part1(t *testing.T) {
	expect := make_expect(t)
	v := part1("test.txt")

	expect(114, v)
}

func Test_part2(t *testing.T) {
	expect := make_expect(t)
	v := part2("test.txt")

	expect(2, v)
}
