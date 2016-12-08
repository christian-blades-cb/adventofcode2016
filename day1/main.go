package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode/utf8"
)

/////////////////////////////////////////////////////////////////////////////////////////////////
// --- Day 1: No Time for a Taxicab ---														   //
// 																							   //
// Santa's sleigh uses a very high-precision clock to guide its								   //
// movements, and the clock's oscillator is regulated by									   //
// stars. Unfortunately, the stars have been stolen... by the Easter						   //
// Bunny. To save Christmas, Santa needs you to retrieve all fifty stars					   //
// by December 25th.																		   //
// 																							   //
// Collect stars by solving puzzles. Two puzzles will be made available						   //
// on each day in the advent calendar; the second puzzle is unlocked when					   //
// you complete the first. Each puzzle grants one star. Good luck!							   //
// 																							   //
// You're airdropped near Easter Bunny Headquarters in a city								   //
// somewhere. "Near", unfortunately, is as close as you can get - the						   //
// instructions on the Easter Bunny Recruiting Document the Elves							   //
// intercepted start here, and nobody had time to work them out further.					   //
// 																							   //
// The Document indicates that you should start at the given coordinates					   //
// (where you just landed) and face North. Then, follow the provided						   //
// sequence: either turn left (L) or right (R) 90 degrees, then walk						   //
// forward the given number of blocks, ending at a new intersection.						   //
// 																							   //
// There's no time to follow such ridiculous instructions on foot,							   //
// though, so you take a moment and work out the destination. Given that					   //
// you can only walk on the street grid of the city, how far is the							   //
// shortest path to the destination?														   //
// 																							   //
// For example:																				   //
// 																							   //
// Following R2, L3 leaves you 2 blocks East and 3 blocks North, or 5 blocks away.			   //
// R2, R2, R2 leaves you 2 blocks due South of your starting position, which is 2 blocks away. //
// R5, L5, R5, R3 leaves you 12 blocks away.												   //
// How many blocks away is Easter Bunny HQ?													   //
/////////////////////////////////////////////////////////////////////////////////////////////////

type turn int

const (
	rightTurn turn = iota
	leftTurn
)

func main() {
	instructionChan := make(chan instruction)
	vectorChan := make(chan vector)
	go lexInstructions(os.Stdin, instructionChan)
	go parseInstructions(instructionChan, vectorChan)

	stops := []vector{vector{x: 0, y: 0}}

	finalVector := vector{}
	for v := range vectorChan {
		finalVector.sum(&v)
		currentCoord := vector{x: finalVector.x, y: finalVector.y}
		stops = append(stops, currentCoord)
	}

	fmt.Println(finalVector.cityBlockDistance())
	firstRevisit := findFirstCrossing(stops)
	fmt.Println(firstRevisit.cityBlockDistance())
}

func direction(origin, destination vector) (direction vector) {
	if destination.x < origin.x {
		direction.x = -1
	} else if destination.x > origin.x {
		direction.x = 1
	}

	if destination.y < origin.y {
		direction.y = -1
	} else if destination.y > origin.y {
		direction.y = 1
	}

	return
}

func findFirstCrossing(stops []vector) vector {
	visits := map[vector]int{}

	prev := stops[0]
	for _, stop := range stops {
		if stop == prev {
			continue
		}

		direction := direction(prev, stop)
		for current := prev; current != stop; current.sum(&direction) {
			visits[current]++
			if visits[current] > 1 {
				return current
			}
		}

		prev = stop
	}
	panic("no revisits")
}

func (v *vector) cityBlockDistance() (blocks int64) {
	var distance int64
	if v.x < 0 {
		distance -= v.x
	} else {
		distance += v.x
	}

	if v.y < 0 {
		distance -= v.y
	} else {
		distance += v.y
	}
	return distance
}

type vector struct {
	x, y int64
}

func (v *vector) sum(v2 *vector) {
	v.x += v2.x
	v.y += v2.y
}

func newVector(direction uint8, steps int64) vector {
	masked := 0x3 & direction
	switch masked {
	case 0x0: // north
		return vector{x: 0, y: steps}
	case 0x1: // east
		return vector{x: steps, y: 0}
	case 0x2: // south
		return vector{x: 0, y: (0 - steps)}
	case 0x3: // west
		return vector{x: (0 - steps), y: 0}
	default:
		panic("how did we even get here?")
	}
}

func parseInstructions(in <-chan instruction, out chan<- vector) {
	var direction uint8
	for token := range in {
		switch token.turn {
		case rightTurn:
			direction++
		case leftTurn:
			direction--
		}
		out <- newVector(direction, token.walk)
	}

	close(out)
}

type instruction struct {
	turn turn
	walk int64
}

func lexInstructions(in io.Reader, out chan<- instruction) {
	scanner := bufio.NewScanner(in)
	scanner.Split(scanCsv)
	for scanner.Scan() {
		word := scanner.Bytes()
		if len(word) < 2 {
			continue
		}

		inst := instruction{}

		turn := word[0]
		switch turn {
		case 'R', 'r':
			inst.turn = rightTurn
		case 'L', 'l':
			inst.turn = leftTurn
		default:
			continue
		}

		steps, err := strconv.ParseInt(string(word[1:]), 10, 64)
		if err != nil {
			continue
		}

		inst.walk = steps
		out <- inst
	}

	close(out)
}

func scanCsv(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpaceOrComma(r) {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpaceOrComma(r) {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func isSpaceOrComma(r rune) bool {
	if r == ',' {
		return true
	}
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}
