package main

import (
	"fmt"
	"strings"
	"os"
)

var (
	newChar byte
	oldChar byte
	strBuf  string
	delimiter0  string
	delimiter1  string
	LastRand int = 999
	wantReviewAll string 
	summaryList = flagcglist
)

// this is the code for -review flag
// review each char in groups of random format adding then to previous char
// for more intence practice
func review(fp *os.File) {

	// do we have delimiters
	if len(delimiterSlice) == 1 {
		delimiter0 = char2psReplacer.Replace(delimiterSlice[0])
	}

	if len(delimiterSlice) == 2 {
		delimiter0 = char2psReplacer.Replace(delimiterSlice[0])
		delimiter1 = char2psReplacer.Replace(delimiterSlice[1])
	}

	if flagcglist == "" {
		fmt.Printf("\nError: the review option requires a non-empty cglist\n")
		os.Exit(1)
	}

	if flaglessonstart != 0 {
		fmt.Printf("\nInclude characters BEFORE lesson: %d in summary?\n",flaglessonstart+1)
		fmt.Printf("They are: %s\n",char2psReplacer.Replace(kochChars[0:flaglessonstart]))
		
		fmt.Printf("\nEnter \"y\", for yes include them: ")

		fmt.Scanf("%s", &wantReviewAll)
		fmt.Println()
	}

	for number, char := range flagcglist {

		if newChar != '0' {
			oldChar = newChar
		}

		newChar = byte(char)
		do5(number) // 5x5 for each char
		doSingle(number)
		doPair(number)
		doEndBlock(number)
	}

	printStrBuf(strBuf, fp)
	return
}

//prevent same rand vals in a row
func genRand(val int) int {

	for ;; {
		rn := rng.Intn(val)
		if rn != LastRand {
			LastRand = rn
			return rn
		}
	}
}

//does 5x5 ie AAAAA AAAAA
func do5(number int) {
	doBlockAnnounce()
	str := strings.Repeat(string(newChar),5)
	strBuf += (str + " " + str + " ")
	doDel(number)
	return
}

// make random groups of 1 char ie AA AAAA AAA
func doSingle(number int) {
	// 8 groups 2-5 chars

	for i := 1; i <= 8; i++ {
		strBuf += fmt.Sprintf("%s ",strings.Repeat(string(newChar),2+genRand(3)))
	}

	doDel(number)
	// special for a blank 
	if number == 0 {
		strBuf += fmt.Sprintln()
	}

	LastRand = 999 // reset
	return
}

// do groups of the new and old char ie KK KM MM MMK
func doPair(number int) {
	var pairList [2]byte

	if number < 1 {
		return
	} else {
		pairList[0] = oldChar
		pairList[1] = newChar
	}

	// 10 groups len 2-6
	for i := 1; i <= 10; i++ {

		for j := 0; j <= 2 + genRand(5); j++ {
			strBuf += string(pairList[genRand(2)])
		}
		strBuf += " "
	}

	doDel(number)
	LastRand = 999 // reset

	return
}

// to block announce delimiter
func doBlockAnnounce() {
	if delimiter0 != "" {
		strBuf += fmt.Sprintf("%s\n",delimiter0)
	}

	return
}

// for end of line delimiter
func doDel(number int) {

	if delimiter1 != "" {
		strBuf += fmt.Sprintf("%s\n",delimiter1)
	} else {
		strBuf += "\n"
	}

	return

}

func doEndBlock(number int) {
	var summaryCount int = 30

	if number == 0 {
		return
	}

	if wantReviewAll == "y" {
		summaryList = char2psReplacer.Replace(kochChars[0:flaglessonstart] + flagcglist[0:number+1])
	} else {
		summaryList = char2psReplacer.Replace(flagcglist[0:number+1])
	}

	// 20 groups
	for i := 1; i <= 20; i++ {

		// length of group
		for j := 0; j <= 2 + genRand(4); j++ {
			strBuf += string(summaryList[rng.Intn(len(summaryList))])
		}
		strBuf += " "
	}

	doDel(number)

	if number % 5 == 0 {
		if delimiter0 != "" {
			strBuf += "\n" + delimiter0 + delimiter0 + "\n"
		} else {
			strBuf += "\n"
		}

		cnt := 0
		for ;; {
			// length of group
			for j := 0; j <= 2 + rng.Intn(7); j++ {
				strBuf += string(summaryList[genRand(len(summaryList))])
			}
			strBuf += " "

			if cnt >= summaryCount {
				break
			} else {
				cnt++
			}
		}

		cnt += 5

		if delimiter1 != "" {
			strBuf += fmt.Sprintf("%s\n",delimiter1)
		} else {
			strBuf += "\n"
		}
	}

	LastRand = 999 // reset
	strBuf += "\n"
}
