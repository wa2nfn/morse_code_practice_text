package main

import (
	"fmt"
	"os"
)

// seed array with LCWO_low, and fill as appropriate
func doCallSigns(chars string, fp *os.File) {
	// first decompose the kochChars to see what we have

	var haveSlash bool
	var haveP bool
	var haveM bool
	var haveQ bool
	var haveR bool
	var ltr []rune
	var dgt []rune
	var trailer []string

	for _, i := range chars {
		// could be either case based on flagcaps so do both
		if i >= rune('A') && i <= rune('Z') {
			ltr = append(ltr, i)

			if i == rune('M') {
				haveM = true
			} else if i == rune('P') {
				haveP = true
			} else if i == rune('R') {
				haveR = true
			} else if i == rune('Q') {
				haveQ = true
			}

		} else if i >= rune('0') && i <= rune('9') {
			dgt = append(dgt, i)
		} else if i == '/' {
			haveSlash = true
		}
	}

	// see if we can continue
	if len(ltr) < 1 {
		fmt.Printf("\nError: to use callSigns, the lesson # must be high enough to have letters.")
		fmt.Printf("\n       tutor <%s> lesson <%s> characters <%s>\n", flagtutor, flaglesson, chars)
		os.Exit(2)
	}

	if len(ltr) == 1 && ltr[0] == 'Q' {
		fmt.Printf("\nError: to use callSigns, the lesson # must be high enough to have a letter other than Q.")
		fmt.Printf("\n       tutor <%s> lesson <%s> characters <%s>\n", flagtutor, flaglesson, chars)
		os.Exit(2)
	}

	if len(dgt) < 1 {
		fmt.Printf("\nError: to use callSigns, the lesson # must be high enough to have numbers.")
		fmt.Printf("\n       tutor <%s> lesson <%s> characters <%s>\n", flagtutor, flaglesson, chars)
		os.Exit(2)
	}

	if haveSlash {
		if haveQ && haveR && haveP {
			trailer = append(trailer, "QRP")
		}
		if haveP {
			trailer = append(trailer, "P")
		}
		if haveM {
			trailer = append(trailer, "M")
			trailer = append(trailer, "MM")
		}

		for _, i := range dgt {
			trailer = append(trailer, string(i))
		}
	}

	doCallSignOutput(string(ltr), string(dgt), trailer, fp)
}

// loop to make the callSigns as needed
func doCallSignOutput(ltr string, dgt string, trailer []string, fp *os.File) {
	strBuf := ""
	tmp := ""

	// make the num callSigns desired
	for i := 0; i < flagnum; i++ {
		tmp = genCS(ltr, dgt, trailer)

		// repeat ?
		if flagrepeat > 1 {
			for cnt := 0; cnt < flagrepeat; cnt++ {
				strBuf += tmp + " "
			}
		} else {
			strBuf += tmp + " "
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

func genCS(ltr string, dgt string, trailer []string) string {

	var outStr []byte
	var curChar rune
	prelen := 2
	rand := 0
	suflen := 3

	// each cs has a Prefix Number Suffix PNS
	// P can be 1-3 chars NOT a Q in first position
	// N is single digit
	// S can be 1-4 char

	// we will have limits on format!
	// No digits in P or S
	// S will just be len 1-3

	// do prefix
	rand = rng.Intn(prelen) + 1

	for cnt := 0; cnt < rand; cnt++ {
		// get index of a ltr
		index := rng.Intn(len(ltr))

		if cnt == 0 {
			for {
				curChar = rune(ltr[index])
				if curChar != rune('q') && curChar != rune('Q') {
					outStr = append(outStr, byte(curChar))
					break
				}
				index = rng.Intn(len(ltr))
			}
		} else {
			curChar = rune(ltr[index])
			outStr = append(outStr, byte(curChar))
		}
	}

	// do dgt
	randDgt := rng.Intn(len(dgt))
	outStr = append(outStr, byte(dgt[randDgt]))

	// do suffix

	rand = rng.Intn(suflen) + 1
	for cnt := 0; cnt < rand; cnt++ {
		// get index of a ltr
		index := rng.Intn(len(ltr))
		outStr = append(outStr, byte(ltr[index]))
	}

	// randomly add a trailer
	if len(trailer) > 0 {
		tStr := ""
		if rng.Intn(9) == 0 {
			outStr = append(outStr, '/')

			// make sure trailer is not the same as N in PNS
			for {
				i := rng.Intn(len(trailer))
				tStr = trailer[i]
				if byte(dgt[randDgt]) != byte(tStr[0]) {
					break
				}
			}

			for _, i := range tStr {
				outStr = append(outStr, byte(i))
			}
		}
	}

	return string(outStr)
}
