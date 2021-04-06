/*
** processes the options: send, sendCheck
 */
package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"
)

var (
	char2psReplacer = strings.NewReplacer("a", "<AS>", "b", "<AR>", "c", "<BT>", "d", "<KA>", "e", "<HH>", "f", "<SK>", "g", "<BK>")
)

// determine which funtion to do
func doSendOpts(fp *os.File) {
	if flagsend != "" && flagsendcheck != "" {
		fmt.Printf("\nError: Option <send> and <sendCheck) are mutually exclusive.\n")
		os.Exit(1)
	} else if flagsend != "" {
		doSendGroups(fp)
	} else {
		doSendCheck(fp)
	}
}

func doSendCheck(fp *os.File) {
	path := strings.Split(flagsendcheck, ",")

	if len(path) == 0 || len(path) != 2 {
		fmt.Printf("\nError: Option <sendCheck> reuires 2 file names in format like: file1,file2.\n")
		os.Exit(1)
	}

	var errVal int

	errVal = readLines(path)

	switch errVal {
	case 1:
		fmt.Printf("\nError: One file must be from MCPT's send option output, the other from a morse  sending output capture (not MCPT).\n")
	default:
		os.Exit(0)
	}
}

/*
// based on codeGroups - this is used for -send option
// generates groups for sending practice based on  number indicating a group of characters with similar attributes
*/

// make random send groups
func doSendGroups(fp *os.File) {
	var tmpOut []rune

	outBuf := make([]rune, 0, (flagcgmax*flagnum)+flagnum)

	sendCharSlice := buildSendSlice() // gets enough chars for entire output

	// make the send groups
	for i := 0; i < flagnum; i++ {

		// tmpOut is our send group
		tmpOut, sendCharSlice = makeSingleSendGroup(sendCharSlice)
		tmpOut = append(tmpOut, ' ')
		outBuf = append(outBuf, tmpOut...)
	}

	// substitue prosigns
	printStrBuf(char2psReplacer.Replace(string(outBuf)), fp)
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
	sendCharSlice := make([]rune, 0, flagcgmax*flagnum) // may be extra

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
		charSlice = append(charSlice, []rune("C.,/?")...)
	}
	// special: use lc letters which we will later map to prosigns
	if strings.Contains(flagsend, "6") {
		charSlice = append(charSlice, []rune("abcdefg")...)
	}

	if len(charSlice) < 5 {
		fmt.Printf("\nError: option -send must include at least one digit from 1-5\n")
		os.Exit(1)
	}

	sendCharSlice = append(sendCharSlice, charSlice...) // get it seeded

	// charSlice now has the user given list of chars
	// then just copy charSlice into charSlice as needed

	need := flagnum * flagcgmax
	for len(sendCharSlice) < need {
		sendCharSlice = append(sendCharSlice, charSlice...)
	}

	return sendCharSlice
}

/*
** below is sendCheck functionality
 */

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path []string) int {
	var gotUser bool
	var gotMCPT bool
	var sendGroupsCompare []string
	var userGroupsCompare []string
	var maxIndex int
	var totalCorrect int
	var totalChars int
	var warningMsg string
	var invalidPSchars bool
	info := color.New(color.FgRed).SprintFunc() // make error visually readable
	ps2charReplacer := strings.NewReplacer("<AS>", "a", "<AR>", "b", "<BT>", "c", "<KA>", "d", "<HH>", "e", "<SK>", "f", "<BK>", "g")
	char2psReplacer := strings.NewReplacer("a", "<AS>", "b", "<AR>", "c", "<BT>", "d", "<KA>", "e", "<HH>", "f", "<SK>", "g", "<BK>")

	// do both files
	for fIndex := 0; fIndex <= 1; fIndex++ {
		// determine which file we have
		whoIsIt, b := determineFile(path[fIndex])

		if whoIsIt == 'm' {
			// process MCPT file
			gotMCPT = true
			sendGroupsCompare = strings.Fields(ps2charReplacer.Replace(string(b)))
		} else {
			// process User file
			gotUser = true
			// make zeros reformatted before compare
			b = bytes.ReplaceAll(b, []byte("0"), []byte("\u00D8"))
			b = bytes.ReplaceAll(b, []byte("\n"), []byte(" "))

			tmpStr := ps2charReplacer.Replace(string(b))
			// must see if the user had invalid prosign delimiters ^<> if so make them *
			if strings.ContainsAny(tmpStr, "^<>") {
				invalidPSchars = true
				b = []byte(tmpStr)
				b = bytes.ReplaceAll(b, []byte("^"), []byte("*"))
				b = bytes.ReplaceAll(b, []byte("<"), []byte("*"))
				b = bytes.ReplaceAll(b, []byte(">"), []byte("*"))
			}
			userGroupsCompare = strings.Fields(ps2charReplacer.Replace(string(b)))
		}
	}

	if gotUser == false || gotMCPT == false {
		return 1
	}

	// needed so we don't over run array
	sendLen := len(sendGroupsCompare)
	userLen := len(userGroupsCompare)

	if sendLen > userLen {
		warningMsg = fmt.Sprintf("\nWarning: The MCPT file had <%d> groups, your file only had <%d> groups.\nOnly the first <%d> groups will be checked.\nYour accuracy is overstated!\n", sendLen, userLen, userLen)
		maxIndex = userLen
	} else if userLen > sendLen {
		warningMsg = fmt.Sprintf("\nWarning: Your file had too many groups <%d>, the MCPT file only had <%d> groups.\nOnly the first <%d> groups will be checked.\nYour accuracy is overstated!\n", userLen, sendLen, sendLen)
		maxIndex = sendLen
	} else {
		maxIndex = sendLen
	}

	fmt.Printf(`
  Levenshtein    MCPT Created        User Sent    
   Distance          Group            Group     
 ============   ===============   =============== 
`)

	// compare char by char the MCPT group to user Group
	for index, sgChar := range sendGroupsCompare {
		if index == maxIndex {
			break
		}

		// get users groups
		ugChar := userGroupsCompare[index]
		// diff strings
		totalChars += len(sgChar)

		levErrors := levenshtein(sgChar, ugChar)
		// print the table
		out := ""
		sp := "   "

		if levErrors > 0 {
			for sgCharIndex, ugCharIndex := 0, 0; ; {
				var tmpChar string

				if sgCharIndex < len(sgChar) && ugCharIndex < len(ugChar) {
					if sgChar[sgCharIndex] != ugChar[ugCharIndex] {
						tmpChar = char2psReplacer.Replace(string(ugChar[ugCharIndex]))
						out += info(tmpChar)
					} else {
						out += char2psReplacer.Replace(string(ugChar[ugCharIndex]))
						totalCorrect++
					}
					sgCharIndex++
					ugCharIndex++
				} else if sgCharIndex == len(sgChar) { // walked off the sgChar
					tmpChar = char2psReplacer.Replace(ugChar[ugCharIndex:])
					out += info(tmpChar)
					break
				} else if ugCharIndex == len(ugChar) {
					tmpChar = char2psReplacer.Replace(ugChar[ugCharIndex:])
					out += info(tmpChar)
					break
				}
			}
		} else {
			totalCorrect += len(sgChar)
			out = ugChar
		}

		fmt.Printf("  %6d        ", levErrors)
		fmt.Printf("%-15s%s%-15s\n", char2psReplacer.Replace(sgChar), sp, char2psReplacer.Replace(out))

		index++
	}

	if totalChars == 0 {
		fmt.Printf("\nError: The MCPT file is empty.\n")
		os.Exit(1)
	}

	if totalCorrect == totalChars {
		fmt.Printf("\nAccuracy: 100%%\n")
	} else {
		fmt.Printf("\nAccuracy: %3.2f%%\n", float32(totalCorrect)*100.0/float32(totalChars))
	}

	if warningMsg != "" {
		fmt.Printf("%s", warningMsg)
	}

	if invalidPSchars {
		fmt.Printf("\n%s", "Warning: Your file had unsupported ProSigns or ProSign characters \"^<>\",\nthey will be shown as \"*\" and will add to your errors.\n")
	}

	return 0

}

// see if the file was the MCPT generated file, or from the users software
func determineFile(path string) (byte, []byte) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("\nError: reading file <%s>. %v\n", path,err)
		os.Exit(1)
	}

	// see if its the MCPT file
	if len(b) >=2 && b[len(b)-2] == byte(0x08) {
		// the MCPT file
		return byte('m'), b[:len(b)-2]
	} else {
		// the users send groups
		b = bytes.ToUpper(b)
		return byte('u'), b
	}
}

// compute the levenshtein error distance of the two groups
// str1 is what MCPT generated, str2 is what the user sent
// func levenshtein(str1, str2 []rune) int {
func levenshtein(str1, str2 string) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}

	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

// utility for levenshtein distance
func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
