package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
)

//////////////////////////////////////////////////////////////////////////////////////////
// --- Day 2: Bathroom Security ---													    //
// 																					    //
// You arrive at Easter Bunny Headquarters under cover of							    //
// darkness. However, you left in such a rush that you forgot to use the			    //
// bathroom! Fancy office buildings like this one usually have keypad				    //
// locks on their bathrooms, so you search the front desk for the code.				    //
// 																					    //
// "In order to improve security," the document you find says, "bathroom			    //
// codes will no longer be written down. Instead, please memorize and				    //
// follow the procedure below to access the bathrooms."								    //
// 																					    //
// The document goes on to explain that each button to be pressed can be			    //
// found by starting on the previous button and moving to adjacent					    //
// buttons on the keypad: U moves up, D moves down, L moves left, and R				    //
// moves right. Each line of instructions corresponds to one button,				    //
// starting at the previous button (or, for the first line, the "5"					    //
// button); press whatever button you're on at the end of each line. If a			    //
// move doesn't lead to a button, ignore it.										    //
// 																					    //
// You can't hold it much longer, so you decide to figure out the code as			    //
// you walk to the bathroom. You picture a keypad like this:						    //
// 																					    //
// 1 2 3																			    //
// 4 5 6																			    //
// 7 8 9																			    //
// Suppose your instructions are:													    //
// 																					    //
// ULL																				    //
// RRDDD																			    //
// LURDL																			    //
// UUUUD																			    //
// 																					    //
// You start at "5" and move up (to "2"), left (to "1"), and left (you				    //
// can't, and stay on "1"), so the first button is 1.								    //
// 																					    //
// Starting from the previous button ("1"), you move right twice (to				    //
// "3") and then down three times (stopping at "9" after two moves and				    //
// ignoring the third), ending up with 9.											    //
// 																					    //
// Continuing from "9", you move left, up, right, down, and left, ending with 8.	    //
// Finally, you move up four times (stopping at "2"), then down once, ending with 5.    //
// So, in this example, the bathroom code is 1985.									    //
// 																					    //
// Your puzzle input is the instructions from the document you found at the front desk. //
// What is the bathroom code?														    //
//////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////
// --- Part Two ---															 //
// 																			 //
// You finally arrive at the bathroom (it's a several minute walk from		 //
// the lobby so visitors can behold the many fancy conference rooms and		 //
// water coolers on this floor) and go to punch in the code. Much to your	 //
// bladder's dismay, the keypad is not at all like you imagined				 //
// it. Instead, you are confronted with the result of hundreds of			 //
// man-hours of bathroom-keypad-design meetings:							 //
// 																			 //
//     1																	 //
//   2 3 4																	 //
// 5 6 7 8 9																 //
//   A B C																	 //
//     D																	 //
// You still start at "5" and stop when you're at an edge, but given the	 //
// same instructions as above, the outcome is very different:				 //
// 																			 //
// You start at "5" and don't move at all (up and left are both edges),		 //
// ending at 5.																 //
// 																			 //
// Continuing from "5", you move right twice and down three times			 //
// (through "6", "7", "B", "D", "D"), ending at D.							 //
// 																			 //
// Then, from "D", you move five more times (through "D", "B", "C", "C",	 //
// "B"), ending at B.														 //
// 																			 //
// Finally, after five more moves, you end at 3.							 //
// 																			 //
// So, given the actual keypad layout, the code would be 5DB3.				 //
// 																			 //
// Using the same instructions in your puzzle input, what is the correct	 //
// bathroom code?															 //
///////////////////////////////////////////////////////////////////////////////

func main() {
	kp := keypad{
		[]rune{'1', '2', '3'},
		[]rune{'4', '5', '6'},
		[]rune{'7', '8', '9'},
	}

	fp, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	printCodes(fp, kp, 1, 1)

	kp2 := keypad{
		[]rune{'🌭', '🌭', '1', '🌭', '🌭'},
		[]rune{'🌭', '2', '3', '4', '🌭'},
		[]rune{'5', '6', '7', '8', '9'},
		[]rune{'🌭', 'A', 'B', 'C', '🌭'},
		[]rune{'🌭', '🌭', 'D', '🌭', '🌭'},
	}

	fp.Seek(0, 0)
	printCodes(fp, kp2, 0, 2)
}

func printCodes(in io.Reader, kp keypad, x, y int) {
	finger := newFinger(kp, x, y)

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		instructions := bytes.Runes(scanner.Bytes())
		finger.do(instructions)
		fmt.Printf("%c", finger.button())
	}
	fmt.Print("\n")
}

type keypad [][]rune

type finger struct {
	x, y       int
	maxX, maxY int
	keypad     keypad
}

func newFinger(keypad keypad, startx, starty int) finger {
	return finger{
		x:      startx,
		y:      starty,
		maxX:   len(keypad[0]) - 1,
		maxY:   len(keypad) - 1,
		keypad: keypad,
	}
}

func (f *finger) do(instructions []rune) {
	var x, y int
	for _, inst := range instructions {
		switch inst {
		case 'U', 'u':
			x = f.x
			y = maxInt(0, f.y-1)
		case 'D', 'd':
			x = f.x
			y = minInt(f.maxY, f.y+1)
		case 'R', 'r':
			y = f.y
			x = minInt(f.maxX, f.x+1)
		case 'L', 'l':
			y = f.y
			x = maxInt(0, f.x-1)
		default:
			panic("invalid instruction")
		}
		if f.keypad[y][x] != '🌭' {
			f.x = x
			f.y = y
		}
	}
}

func (f *finger) button() rune {
	return f.keypad[f.y][f.x]
}

func minInt(x, y int) int {
	return int(math.Min(float64(x), float64(y)))
}

func maxInt(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}
