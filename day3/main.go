package main

import (
	"bufio"
	"fmt"
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
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
	fmt.Println(validTriangles)
}
