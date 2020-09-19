package main

import (
	"os"
)

// make random code groups
// uses the presaved chars in charSlice based on uniform distribution
func makeGroups(fp *os.File) {
	strBuf := ""
	var tmpOut []rune

	charSlice := buildCharSlice()

	// make the code groups
	for i := 0; i < flagnum; i++ {

		var xOut []rune
		rand := 3

		// tmpOut is our code group
		tmpOut, charSlice = makeSingleGroup(charSlice)
		// do we need prefix or suffix

		if flagrandom {
			if flagsufmin >= 1 || flagpremin >= 1 {
				// 0 - neither ix, 1 prefix,2 do suffix, 3 both
				rand = rng.Intn(4)
			}
		}

		// end raw word, and get back word to print
		// do we need prefix?
		if flagpremin >= 1 && (rand == 3 || rand == 1) {
			xOut = []rune(ixStr("p"))
			xOut = append(xOut, tmpOut...)
			tmpOut = xOut
		}

		// do we need a suffix or just a space
		if flagsufmin >= 1 && (rand == 3 || rand == 2) {
			tmpOut = tmpOut[0 : len(tmpOut)-1]
			// make string to rune slice
			xOut = []rune(ixStr("s"))
			tmpOut = append(tmpOut, xOut...)
			tmpOut = append(tmpOut, ' ')
		}

		// text repeat!
		if flagrepeat > 0 {
			// we need to repeat
			temp := tmpOut

			for cnt := 1; cnt < flagrepeat; cnt++ {
				// wordOut is the word plus trailing space already
				temp = append(temp, tmpOut...)
			}
			strBuf += string(temp)

		} else {
			// non repeat case
			strBuf += string(tmpOut)
		}

		if flagDMmin >= 1 && (flagDR == false || (flagDR == true && flipFlop())) {
			if flagDMmin == flagDMmax {
				for count := 0; count < flagDMmax; count++ {
					strBuf += delimiterSlice[rng.Intn(len(delimiterSlice))]
				}
			} else {
				for count := 0; count < (flagDMmin + rng.Intn(flagDMmax-flagDMmin+1)); count++ {
					strBuf += delimiterSlice[rng.Intn(len(delimiterSlice))]
				}
			}

			strBuf += " "
		}
	}

	printStrBuf(strBuf, fp)
}

/*
** make a code group of random length
** character pulled from byte slice that even distribution
** of characters.
 */
func makeSingleGroup(charSlice []rune) ([]rune, []rune) {
	var cg []rune
	var tmp rune
	gl := flagcgmin

	// choose random group len from min to max
	if flagcgmax != flagcgmin {
		gl = rng.Intn(flagcgmax-flagcgmin+1) + flagcgmin
	}

	if gl > len(charSlice) {
		gl = len(charSlice)
		gl--
	}

	for i := 0; i < gl; i++ {
		tmp, charSlice = getRandomChar(charSlice)

		if rune2psMap != nil {
			s := rune2psMap[tmp]
			if s != "" {
				cg = append(cg, []rune(s)...)
			} else {
				cg = append(cg, tmp)
			}
		} else {
			cg = append(cg, tmp)
		}
	}

	cg = append(cg, ' ')
	return cg, charSlice
}

// used for codeGroup
func getRandomChar(randCharSlice []rune) (rune, []rune) {
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
// buildCharSlice - create a byte slice to use for codeGroups
//
func buildCharSlice() []rune {

	// make slice of chars for MAY NEED use later
	// if word mode, we only need for MAX possible codeGroups
	numChars := 0

	if flagMixedMode == 0 {
		numChars = flagcgmax * flagnum // may be extra
	} else {
		m := flagnum / flagMixedMode
		if m < 1 {
			m = 1
		}
		numChars = flagcgmax * m // may be extra
	}

	charSlice := make([]rune, 0, numChars)
	cgSlice := flagCglistRune // start with cglist chars

	// see if we need to add the prosign runes
	if ps2runeMap != nil {
		for _, val := range ps2runeMap {
			cgSlice = append(cgSlice, val)
		}
	}

	charSlice = cgSlice

	// charSlice now has the user given list of chars
	// then just copy cgSlice into charSlice as needed
	if len(cgSlice) < numChars {
		// flush out the charSlice to max we may need
		factor := numChars / len(cgSlice)
		factor-- // we have the original already

		// only does FULL slice
		for ; factor > 0; factor-- {
			charSlice = append(charSlice, cgSlice...)
		}

		// may still be a partial shortage
		howShort := numChars - len(charSlice)

		for _, key := range cgSlice {

			if howShort == 0 {
				break
			}
			charSlice = append(charSlice, key)
			howShort--
		}
	}

	return charSlice
}
