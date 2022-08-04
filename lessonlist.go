package main

import (
	"fmt"
	"os"
	"time"
	"math/rand"
)

//
//
func lessonlist(fp *os.File) {
	tmpArr := []rune(flagcglist)

	if len(flagcglist) < 1 {
		fmt.Println("\nError: input string too short, increase <lesson> for your <tutor>.\n")
		os.Exit(9)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	if flagrandom {
		rand.Shuffle(len(tmpArr), func(i, j int) { tmpArr[i], tmpArr[j] = tmpArr[j], tmpArr[i] })
	}

	// we use the cglist to send even increasing strings in order they
	// are given

	outBuf := ""

	out:
	for j := 0; ; {
		i := len(string(tmpArr))
		notDone := false

		for k := 0; k <= i; {
			outBuf += string(tmpArr[0:k])
			if outBuf != "" {
				outBuf += " "
			}

			if j >= flagnum {
			     break out
			}
			k++
			j++
			outBuf += "\n"
			notDone = true
		}

		if flagrandom && notDone {
			rand.Shuffle(len(tmpArr), func(i, j int) { tmpArr[i], tmpArr[j] = tmpArr[j], tmpArr[i] })
			notDone = false
		}
	}

	printStrBuf(convertRunes(outBuf),fp)
}
