/*
** processes the options: send, sendCheck
 */

package main

import (
	"bytes"
	"fmt"
	"github.com/jwalton/gchalk"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

var (
	char2psReplacer *strings.Replacer
	ps2charReplacer *strings.Replacer
	gotCarat bool
	validCharPS string
	invalidCharPS string
)

// determine which function to do
func doSendOpts(fp *os.File) {
	if flagsend != "" && flagsendcheck != "" {
		fmt.Printf("\n     Error: Option <send> and <sendCheck> are mutually exclusive.\n")
	} else if flagsend != "" {
		doSendGroups(fp)
	} else {
		doSendCheck(fp)
	}
	os.Exit(1)
}

func doSendCheck(fp *os.File) {
	var path []string
	var sep = ","

	if strings.Contains(flagsendcheck,sep) {
		// ps format <xx>
		char2psReplacer = strings.NewReplacer("a", "<AS>", "b", "<AR>", "c", "<BT>", "d", "<KA>", "e", "<HH>", "f", "<SK>", "g", "<BK>","h","<AA>","i","<CT>","j","<KN>","k","<VA>","l","<SN>")
		ps2charReplacer = strings.NewReplacer("<AS>", "a", "<AR>", "b", "<BT>", "c", "<KA>", "d", "<HH>", "e", "<SK>", "f", "<BK>", "g","<AA>","h","<CT>","i","<KN>","j","<VA>","k","<SN>","l" )
		validCharPS = "<>"
		invalidCharPS = "^"
	} else {
		// ps format ^xx
		sep = "^"
		validCharPS = sep
		invalidCharPS = "<>"
		char2psReplacer = strings.NewReplacer("a", "^AS", "b", "^AR", "c", "^BT", "d", "^KA", "e", "^HH", "f", "^SK", "g", "^BK","h","^AA","i","^CT","j","^KN","k","^VA","l","^SN")
		ps2charReplacer = strings.NewReplacer("^AS", "a", "^AR", "b", "^BT", "c", "^KA", "d", "^HH", "e", "^SK", "f", "^BK", "g","^AA","h","^CT","i","^KN","j","^VA","k","^SN","l")
	}

	path = strings.Split(flagsendcheck, sep)

	if len(path) == 0 || len(path) != 2 {
		fmt.Printf("\n     Error: Option <sendCheck> requires 2 file names, in format like: file1,file2 (if prosigns have <XX> format\nor file1^file2, if prosigms have ^XX format.\n")
		os.Exit(1)
	}

	var errVal int

	errVal = readLines(path)

	switch errVal {
	case 1:
		fmt.Printf("\n     Error: One of files MUST be the captured text from your morse sending.\n")
	case 2:
		fmt.Printf("\n     Error: One of files MUST be an MCPT generated file of practice material.\n")
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
	// always use <> fomat

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
		newChar = '\u00D8' // make zeros more readable 
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
		fmt.Printf("\n     Error: option -send must include at least one digit from 1-5\n")
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
	var invalidPScharsFlag bool
	//var colorError = gchalk.WithBgBlack().BrightRed
	var colorExtra = gchalk.BrightGreen
	//var colorError = gchalk.WithBgBlack().BrightRed
	var colorError = gchalk.BrightRed
	var miss bool
	var extra bool

	gchalk.SetLevel(gchalk.LevelAnsi16m)

	// do both files
	for fIndex := 0; fIndex <= 1; fIndex++ {
		// determine which file we have
		whoIsIt, b := determineFile(path[fIndex])
		b = bytes.ReplaceAll(b, []byte("0"), []byte("\u00D8"))  // unfancy all zeros

		if whoIsIt == 'm' {
			// process MCPT file
			gotMCPT = true
			sendGroupsCompare = strings.Fields(ps2charReplacer.Replace(string(b)))
		} else {
			// process User file
			gotUser = true
			// convert any NL to space
			b = bytes.ReplaceAll(b, []byte("\n"), []byte(" "))

			tmpStr := ps2charReplacer.Replace(string(b))
			// must see if the user had invalid prosign delimiters <> if so make them *
			if strings.ContainsAny(tmpStr, invalidCharPS) {
				invalidPScharsFlag = true
				b = []byte(tmpStr)

				if gotCarat {
					b = bytes.ReplaceAll(b, []byte("^"), []byte("*"))
				} else {
					b = bytes.ReplaceAll(b, []byte("<"), []byte("*"))
					b = bytes.ReplaceAll(b, []byte(">"), []byte("*"))
				}
			}
			userGroupsCompare = strings.Fields(ps2charReplacer.Replace(string(b)))
		}
	}

	if gotUser == false {
		return 1
	}
	if gotMCPT == false {
		return 2
	}

	// needed so we don't over run array
	sendLen := len(sendGroupsCompare)
	userLen := len(userGroupsCompare)

	if sendLen > userLen {
		warningMsg = fmt.Sprintf("\n     Note: The MCPT file had <%d> groups, your file only had <%d>.\n     Only the first <%d> groups will be checked.\n\n     Your accuracy is overstated!\n", sendLen, userLen, userLen)
		maxIndex = userLen
	} else if userLen > sendLen {
		warningMsg = fmt.Sprintf("\n     Note: Your file had too many groups <%d>, the MCPT file only had <%d>.\n     Only the first <%d> groups will be checked.\n\n     Your accuracy is overstated!\n", userLen, sendLen, sendLen)
		maxIndex = sendLen
	} else {
		maxIndex = sendLen
	}

	fmt.Printf(`
            Levenshtein       MCPT Created            Your Sent    
             Distance            Group                  Group     
            ===========   ====================   ==================== 
`)

	// compare char by char the MCPT group to user Group
	for index, sgChar := range sendGroupsCompare {
		var out string
		var missed string
		var tmpChar string
		var missCnt int
		var Index int

		// get users groups
		ugChar := userGroupsCompare[index]
		// diff strings
		totalChars += len(sgChar)

		levErrors := levenshtein(sgChar, ugChar)
		// print the table

		fmt.Printf("            %6d        ", levErrors)

		if levErrors > 0 {
			for ; Index < len(sgChar) && Index < len(ugChar); Index++ {

				if Index < len(sgChar) && Index < len(ugChar) {
					if sgChar[Index] != ugChar[Index] {
						// mismatch - color bad data
						tmpChar = char2psReplacer.Replace(string(ugChar[Index]))
						out += colorError(tmpChar)
					} else {
						// both matched good!
						out += char2psReplacer.Replace(string(sgChar[Index]))
						totalCorrect++
					}
				}
			}

			if Index + 1 == len(sgChar) && Index + 1 == len(ugChar) {
				break // finished
			}

			// one array is done
			if len(ugChar[Index:]) >= 1 && Index + 1  > len(sgChar) { // walked off the sgChar
				// extra sent by user
				tmpChar = char2psReplacer.Replace(ugChar[Index:])
				tmpChar = out + colorExtra(tmpChar)
				extra = true  // once set, forever set
			} else if len(sgChar[Index:]) >= 1 {
				// more in system column, user missed some
				missCnt = len(sgChar) - len(ugChar)
				missed = char2psReplacer.Replace(sgChar[:])
				miss = true  // once set, forever set
			}

		} else {
			totalCorrect += len(sgChar)
			out = ugChar
		}

		// FINN
		sgChar = char2psReplacer.Replace(sgChar)
		if missed != "" {
			s := missed + " -" + strconv.Itoa(missCnt)
			fmt.Printf("%-20s   %-20s\n", s, out) // missed stuff from col 1
			missed = ""
			missCnt = 0
		} else if extra {
			fmt.Printf("%-20s   %-20s\n", sgChar, tmpChar)  // extra stuff in second col
			extra = false
		} else {
			fmt.Printf("%-20s   %-20s\n", sgChar, out)  // all matched
		}

		out = ""
		tmpChar = ""
		sgChar = ""
		index++

		if index == maxIndex {
			break
		}

	}

	if totalChars == 0 {
		fmt.Printf("\n     Error: The MCPT file is empty.\n")
		os.Exit(1)
	}

	if totalCorrect == totalChars {
		fmt.Printf("\n     Accuracy: 100%%\n")
	} else {
		fmt.Printf("\n     Accuracy: %3.2f%%\n", float32(totalCorrect)*100.0/float32(totalChars))
	}

	if warningMsg != "" {
		fmt.Printf("%s", warningMsg)
	}

	if invalidPScharsFlag {
		fmt.Printf("\n     Warning: Your file had unsupported ProSigns or ProSign character(s) \"%s\",\nthey will be shown as \"*\" and will add to your errors.\n",invalidCharPS)
	}

	if miss {
		fmt.Printf("\n     Warning: You MISSED sending some characters, noted by -X (a prosign counts as 1).\n")
		fmt.Printf("\n              If your send groups following the missed indicator are (errors), an EXTRA sent space may have split a group.")
		fmt.Printf("\n              Edit your send file, fix the space error, and rerun.\n")
		fmt.Printf("\n              The space should be just before the user groups turned RED.\n")
	}

	return 0

}

// see if the file was the MCPT generated file, or from the users software
func determineFile(path string) (byte, []byte) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("\n     Error: reading file <%s>. %v\n", path, err)
		os.Exit(1)
	}

	b = bytes.ToUpper(b) // in case source was not -send

	// see if its the MCPT file
	if  bytes.Contains(b,[]byte("\u0008")) {
		// the MCPT file
		b = b[:len(b) - 2]
		return byte('m'), b
	} else {
		// the users send groups
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

