package main

import (
	"os"
	"fmt"
	"strings"
)

var (
	rune2prosign = make(map[rune][]rune)
)

/*
// based on codeGroups - this is used for -send option
// generates groups for sending practice based on  number indicating a group of characters with similar attributes
*/

// make random send groups
func doSendGroups(fp *os.File) {
	var tmpOut []rune
	//wdl rune2prosign = make([]rune, 0, 6)  // only for group 6
	outBuf := make([]rune, 0, (flagcgmax * flagnum)+flagnum ) 

	sendCharSlice := buildSendSlice() // gets enough chars for entire output

	// make the send groups
	for i := 0; i < flagnum; i++ {

		// tmpOut is our send group
		tmpOut, sendCharSlice = makeSingleSendGroup(sendCharSlice)
		tmpOut = append(tmpOut, ' ')
		outBuf = append(outBuf, tmpOut...)
	}


	printStrBuf(string(outBuf), fp)
}

/*
** make a send group of random length
** character pulled from byte slice that even distribution
** of characters.
 */
func makeSingleSendGroup(charSlice []rune) ([]rune, []rune) {
	var cg []rune
	var tmp rune
	gl := flagcgmin

	// choose random group len from min to max
	if flagcgmax != flagcgmin && flaghead == false {
		gl = rng.Intn(flagcgmax-flagcgmin+1) + flagcgmin
	}

	if gl > len(charSlice) {
		gl = len(charSlice)
		gl--
	}

	cg = append(cg, ' ')
	for i := 0; i < gl; i++ {
		tmp, charSlice = getRandomSendChar(charSlice)
		// see if its from group6 prosigns
		switch tmp {
		case 'a','b','c','d','e','f','g':
			cg = append(cg, rune2prosign[tmp]...)
		default:
			cg = append(cg, tmp)
		}
	}

	return cg, charSlice
}

// used for sendGroup
func getRandomSendChar(randCharSlice []rune) (rune, []rune) {
	sLen := len(randCharSlice)

	index := rng.Intn(sLen)
	newChar := randCharSlice[index] // to be returned

	// eat the value used
	sLen--
	randCharSlice[index] = randCharSlice[sLen]
	randCharSlice = randCharSlice[:sLen]

	if newChar == '0' {
		 newChar = '\u00D8' // make zeros more readable ? wdl
	}

	return newChar, randCharSlice
}

//
// buildSendSlice - create a byte slice to use for codeGroups
//
func buildSendSlice() []rune {
	// based on the ints in flagsend we populate the charSlice
	sendCharSlice := make([]rune, 0, flagcgmax * flagnum) // may be extra

	charSlice := make([]rune, 0, 50) // may be extra
	if strings.Contains(flagsend, "1") {
		charSlice = append(charSlice, []rune("EIHMOST50")...)
	}
	if strings.Contains(flagsend, "2") {
		charSlice = append(charSlice, []rune("ABDGJNUVWZ")...)
	}
	if strings.Contains(flagsend, "3") {
		charSlice = append(charSlice, []rune("12346789")...)
	}
	if strings.Contains(flagsend, "4") {
		charSlice = append(charSlice, []rune("FKLRPQXY")...)
	}
	if strings.Contains(flagsend, "5") {
		charSlice = append(charSlice, []rune("C.,/?=")...)
	}
	// special: use lc letters which we will later map to prosigns
	if strings.Contains(flagsend, "6") {
		charSlice = append(charSlice, []rune("abcdefg")...)
		// populate the map
		rune2prosign['a'] = []rune("<AR>")
		rune2prosign['b'] = []rune("<AS>")
		rune2prosign['c'] = []rune("<BT>")
		rune2prosign['d'] = []rune("<KA>")
		rune2prosign['e'] = []rune("<HH>")
		rune2prosign['f'] = []rune("<SK>")
		rune2prosign['g'] = []rune("<BK>")
	}

	if len(charSlice) < 5 {
		fmt.Printf("\nError: option -send must include at least one digit from 1-5\n")
		os.Exit(1)
	}

	sendCharSlice = append(sendCharSlice, charSlice...) // get it seeded

	// charSlice now has the user given list of chars
	// then just copy charSlice into charSlice as needed

	need := flagnum * flagcgmax
	for ( len(sendCharSlice) < need ) {
		sendCharSlice = append(sendCharSlice, charSlice...)
	}

	return sendCharSlice
}
