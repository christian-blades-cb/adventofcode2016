package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"unicode"
)

/////////////////////////////////////////////////////////////////////////////////
// --- Day 4: Security Through Obscurity ---								   //
// 																			   //
// Finally, you come across an information kiosk with a list of rooms. Of	   //
// course, the list is encrypted and full of decoy data, but the			   //
// instructions to decode the list are barely hidden nearby. Better			   //
// remove the decoy data first.												   //
// 																			   //
// Each room consists of an encrypted name (lowercase letters separated		   //
// by dashes) followed by a dash, a sector ID, and a checksum in square		   //
// brackets.																   //
// 																			   //
// A room is real (not a decoy) if the checksum is the five most common		   //
// letters in the encrypted name, in order, with ties broken by				   //
// alphabetization. For example:											   //
// 																			   //
// aaaaa-bbb-z-y-x-123[abxyz] is a real room because the most common		   //
// letters are a (5), b (3), and then a tie between x, y, and z, which		   //
// are listed alphabetically.												   //
// 																			   //
// a-b-c-d-e-f-g-h-987[abcde] is a real room because although the letters	   //
// are all tied (1 of each), the first five are listed alphabetically.		   //
// 																			   //
// not-a-real-room-404[oarel] is a real room.								   //
// 																			   //
// totally-real-room-200[decoy] is not.										   //
// 																			   //
// Of the real rooms from the list above, the sum of their sector IDs is 1514. //
// 																			   //
// What is the sum of the sector IDs of the real rooms?						   //
/////////////////////////////////////////////////////////////////////////////////

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	sectorSum := 0
	i := 0
	for scanner.Scan() {
		txt := scanner.Text()
		if sector, ok := sectorOf(txt); ok {
			sectorSum += sector
		}
		i++
	}

	fmt.Println(sectorSum)
}

type letterFreq struct {
	letter rune
	freq   int
}

type freqSorter []letterFreq

func newFreqSorter(freq map[rune]int) freqSorter {
	fs := make(freqSorter, len(freq))

	i := 0
	for k, v := range freq {
		fs[i] = letterFreq{letter: k, freq: v}
		i++
	}
	return fs
}

func (fs freqSorter) Len() int {
	return len(fs)
}

func (fs freqSorter) Less(i, j int) bool {
	if fs[i].freq > fs[j].freq {
		return true
	}
	if fs[i].freq < fs[j].freq {
		return false
	}
	return fs[i].letter < fs[j].letter
}

func (fs freqSorter) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func sectorOf(input string) (sector int, ok bool) {
	freq := map[rune]int{}
	sector, err := strconv.Atoi(input[len(input)-10 : len(input)-7])
	if err != nil {
		panic(err)
	}

	for _, r := range input[:len(input)-10] {
		if unicode.IsLetter(r) {
			freq[r]++
		}
	}

	checksum := input[len(input)-6 : len(input)-1]
	ok = isRealRoom(freq, checksum)
	return
}

func isRealRoom(freq map[rune]int, checksum string) bool {
	sorted := newFreqSorter(freq)
	sort.Sort(&sorted)

	for i, r := range checksum {
		if sorted[i].letter != r {
			return false
		}
	}

	return true
}
