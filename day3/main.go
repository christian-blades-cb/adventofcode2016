package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////
// --- Day 3: Squares With Three Sides ---								  //
// 																		  //
// Now that you can think clearly, you move deeper into the labyrinth of  //
// hallways and office furniture that makes up this part of Easter Bunny  //
// HQ. This must be a graphic design department; the walls are covered in //
// specifications for triangles.										  //
// 																		  //
// Or are they?															  //
// 																		  //
// The design document gives the side lengths of each triangle it		  //
// describes, but... 5 10 25? Some of these aren't triangles. You can't	  //
// help but mark the impossible ones.									  //
// 																		  //
// In a valid triangle, the sum of any two sides must be larger than the  //
// remaining side. For example, the "triangle" given above is impossible, //
// because 5 + 10 is not larger than 25.								  //
// 																		  //
// In your puzzle input, how many of the listed triangles are possible?	  //
////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////
// --- Part Two ---														 //
// 																		 //
// Now that you've helpfully marked up their design documents, it occurs //
// to you that triangles are specified in groups of three				 //
// vertically. Each set of three numbers in a column specifies a		 //
// triangle. Rows are unrelated.										 //
// 																		 //
// For example, given the following specification, numbers with the same //
// hundreds digit would be part of the same triangle:					 //
// 																		 //
// 101 301 501															 //
// 102 302 502															 //
// 103 303 503															 //
// 201 401 601															 //
// 202 402 602															 //
// 203 403 603															 //
// 																		 //
// In your puzzle input, and instead reading by columns, how many of the //
// listed triangles are possible?										 //
///////////////////////////////////////////////////////////////////////////

func main() {
	fd, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fmt.Println(numValidRows(fd))
	fd.Seek(0, 0)
	fmt.Println(numValidCols(fd))
}

func numValidCols(in io.Reader) int {
	vals := readColumnarVals(in)
	var validTriangles int
	var triangle [3]int
	for i, x := range vals {
		triangle[i%3] = x
		if i%3 == 2 {
			a, b, c := triangle[0], triangle[1], triangle[2]
			if a+b > c && a+c > b && b+c > a {
				validTriangles++
			}
		}
	}
	return validTriangles
}

func readColumnarVals(in io.Reader) []int {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)
	columnar := [3][]int{[]int{}, []int{}, []int{}}
	i := 0
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		columnar[i] = append(columnar[i], x)
		if i == 2 {
			i = 0
		} else {
			i++
		}
	}
	return append(append(columnar[0], columnar[1]...), columnar[2]...)
}

func numValidRows(in io.Reader) int {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanWords)
	var err error
	var validTriangles int
	var triangle [3]int
	i := 0
	for scanner.Scan() {
		triangle[i], err = strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		if i == 2 {
			if triangle[0]+triangle[1] > triangle[2] &&
				triangle[0]+triangle[2] > triangle[1] &&
				triangle[1]+triangle[2] > triangle[0] {
				validTriangles++
			}
			i = 0
		} else {
			i++
		}
	}
	return validTriangles
}
