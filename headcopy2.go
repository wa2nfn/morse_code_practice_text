package main

import (
	"strings"
)

// headcopy is a stripped down version of prepWord elimationg the use of many options that
// are not required and disassembling the word to its parts as needed
func headcopy2(wordOut string, index int, charSlice []rune) (string, []rune) {
	inWord := strings.Split(wordOut,"")
	outStr := ""
	lastStr := ""

	// break up "word" and send out strings until all chars have been used
	for _,char := range inWord {
		lastStr = lastStr + char
		outStr += lastStr + " "
	}


	// text repeat!
	if flagrepeat > 0 {
		// we need to repeat

		for cnt := 1; cnt < flagrepeat; cnt++ {
			// wordOut is the word plus trailing space already
			outStr += wordOut
			outStr += " "
		}
	}


	// use delimiter 
	if flagDMmin >= 1 && (flagDR == false || (flagDR == true && flipFlop())) {

		if flagDMmin == flagDMmax {
			for count := 0; count < flagDMmax; count++ {
				outStr += delimiterSlice[rng.Intn(len(delimiterSlice))]
			}
		} else {
			for count := 0; count < (flagDMmin + rng.Intn(flagDMmax-flagDMmin+1)); count++ {
				outStr += delimiterSlice[rng.Intn(len(delimiterSlice))]
			}
		}
	}

	outStr += "\n" // special case here so Precision CW Tutor can read a "word" as a line for Ctrl-L use

	return outStr, charSlice
}

