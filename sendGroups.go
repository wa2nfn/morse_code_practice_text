package main

import (
	"os"
	"fmt"
	"strings"
)

/*
// based on codeGroups - this is used for -send option
// generates groups for sending practice based on  number indicating a group of characters with similar attributes
*/

// make random send groups
func doSendGroups(fp *os.File) {
	var tmpOut []rune
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

	for i := 0; i < gl; i++ {
		tmp, charSlice = getRandomSendChar(charSlice)
		cg = append(cg, tmp)
	}

	cg = append(cg, ' ')

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
