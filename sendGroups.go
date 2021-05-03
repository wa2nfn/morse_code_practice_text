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
	"regexp"
	"strings"
)

var (
	char2psReplacer = strings.NewReplacer("a", "<AS>", "b", "<AR>", "c", "<BT>", "d", "<KA>", "e", "<HH>", "f", "<SK>", "g", "<SN>", "0", "\u00D8")

	ps2charReplacer = strings.NewReplacer("<AS>", "a", "<AR>", "b", "<BT>", "c", "<KA>", "d", "<HH>", "e", "<SK>", "f", "<SN>", "g", "\u00D8", "0")

	MCPTps2charReplacer = strings.NewReplacer("<AS>", "a", "<AR>", "b", "<BT>", "c", "<KA>", "d", "<HH>", "e", "<SK>", "f", "<SN>", "g", "\u00D8", "0")

	gotCarat      bool
	validCharPS   string
	invalidCharPS string
)

// determine which function to do
func doSendOpts(fp *os.File) {
	if flagsend != "" && flagsendcheck != "" {
		fmt.Printf("\n Error: Option <-send> and <-sendCheck> are mutually exclusive.\n")
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

	if strings.Contains(flagsendcheck, sep) {
		// ps format <xx>
		validCharPS = "<>"
		invalidCharPS = "^"
	} else {
		// ps format ^xx
		sep = "^"
		validCharPS = sep
		invalidCharPS = "<>"
		char2psReplacer = strings.NewReplacer("a", "^AS", "b", "^AR", "c", "^BT", "d", "^KA", "e", "^HH", "f", "^SK", "g", "^SN", "0", "\u00D8")
		ps2charReplacer = strings.NewReplacer("^AS", "a", "^AR", "b", "^BT", "c", "^KA", "d", "^HH", "e", "^SK", "f", "^SN", "g", "\u00D8", "0")
		gotCarat = true
	}

	path = strings.Split(flagsendcheck, sep)

	if len(path) == 0 || len(path) != 2 {
		fmt.Printf("\n Error: Option <-sendCheck> requires 2 file names, in format like: file1,file2 (if Prosigns have <XX> format\nor file1^file2, if ProSigns have ^XX format.\n")
		fmt.Printf("\n        The file from the CW capture, MUST have \"%s\" before it. E.g. -sendCheck=C:capture.txt,practice.txt\n")
		os.Exit(1)
	}

	var errVal int

	errVal = readLines(path)

	switch errVal {
	case 1:
		fmt.Printf("\n Error: Only one file can be a CW captured text from your sending software. Prefix its name in the sendCheck with: %s.\n\n", gchalk.BrightRed("C: or c:"))
	case 2:
		fmt.Printf("\n Error: Only one file can be a practice text file. E.g. <practice.txt>.\n        Do NOT prefix its name, with \"C: or c:\".\n\n")
	default:
		os.Exit(0)
	}
	os.Exit(0)
}

/*
// based on codeGroups - this is used for -send option
// generates groups for sending practice based on  number indicating a group of characters with similar attributes
*/

// make random send groups
func doSendGroups(fp *os.File) {
	var tmpOut []rune
	// always use <> format

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

	return newChar, randCharSlice
}

//
// buildSendSlice - create a byte slice to use for codeGroups
//
func buildSendSlice() []rune {
	var gotDigit bool
	// based on the ints in flagsend we populate the charSlice
	sendCharSlice := make([]rune, 0, flagcgmax*flagnum) // may be extra

	charSlice := make([]rune, 0, 50) // may be extra

	if strings.Contains(flagsend, "0") {
		// no op allows use of cglist
		gotDigit = true
	}
	if strings.Contains(flagsend, "1") {
		charSlice = append(charSlice, []rune("EIHMOST50")...)
		gotDigit = true
	}
	if strings.Contains(flagsend, "2") {
		charSlice = append(charSlice, []rune("ABDGJNUVWZ")...)
		gotDigit = true
	}
	if strings.Contains(flagsend, "3") {
		charSlice = append(charSlice, []rune("12346789")...)
		gotDigit = true
	}
	if strings.Contains(flagsend, "4") {
		charSlice = append(charSlice, []rune("FKLRPQXY")...)
		gotDigit = true
	}
	if strings.Contains(flagsend, "5") {
		charSlice = append(charSlice, []rune("C.,/?")...)
		gotDigit = true
	}
	// special: use lc letters which we will later map to prosigns
	if strings.Contains(flagsend, "6") {
		charSlice = append(charSlice, []rune("abcdefg")...)
		gotDigit = true
	}

	if flagcglist != "" {
		charSlice = append(charSlice, []rune(flagcglist)...)
	}

	if gotDigit == false {
		fmt.Printf("\n Error: option <-send> must include at least one digit from 1-5\n")
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
	var sendGroupsCompareNoSpace string
	var userGroupsCompareNoSpace string
	var maxIndex int
	var totalCorrect int
	var totalChars int
	var warningMsg string
	var colorExtra = gchalk.BrightGreen
	var colorError = gchalk.BrightRed
	var colorMiss = gchalk.WithBgBlue().BrightCyan
	var miss bool
	var extra bool
	var extraForever bool
	var captureFile = ""
	var practiceFile = ""
	badPS := regexp.MustCompile(`<[^>]*>`)  // find <..--..> errors and anything since supported PD already done
	hhErr := regexp.MustCompile(`<.{8,8}>`) // find <HH> as code

	gchalk.SetLevel(gchalk.LevelAnsi256)

	// do both files
	for fIndex := 0; fIndex <= 1; fIndex++ {
		// determine which file we have
		whoIsIt, b, fName := determineFile(path[fIndex])
		// b is already in UC
		// convert any NL to space
		b = bytes.ReplaceAll(b, []byte("\n"), []byte(" "))

		if whoIsIt == 'm' {
			// process MCPT file
			gotMCPT = true
			practiceFile = fName
			// good PS replace
			bStr := MCPTps2charReplacer.Replace(string(b))

			if strings.ContainsAny(bStr, "<>^") {
				fmt.Printf("\n Warning: Your practice file <%s> contains character(s) \"<>^\" in addition to those in supported ProSigns.", practiceFile)
				fmt.Printf("\n          This will add to your error count.\n\n")

				// look for < anything>
				if hhErr.MatchString(bStr) {
					bStr = hhErr.ReplaceAllString(bStr, "*")
				}
			}

			sendGroupsCompare = strings.Fields(bStr)
			sendGroupsCompareNoSpace = strings.ReplaceAll(bStr, " ", "")

		} else if whoIsIt == 'u' {
			// process User file
			gotUser = true
			captureFile = fName

			// good PS replaced
			var tmp = ps2charReplacer.Replace(string(b))

			// look for <.......> errors first and invalid prosigns
			if badPS.MatchString(tmp) {
				tmp = badPS.ReplaceAllString(tmp, "*")
			}

			if gotCarat {
				// got carat so PS needs looklike ^BT NOT <BT>
				if strings.ContainsAny(tmp, "<>") {
					// used ^ sep should use ","
					fmt.Printf("\n Warning: Your CW capture file <%s> had unsupported ProSign format\n          characters \"<>\" (i.e. <>).", captureFile)
					fmt.Printf("\n\n          If those are correct for your ProSigns, use a comma \",\"\n          between the file names (i.e.-send=%cature.txt,practice.txt).\n", gchalk.Yellow("C:"))

					os.Exit(88)
				}
			} else {
				// have , so need <>  ie <BT> NOT ^BT
				if strings.ContainsAny(tmp, "^") {
					// used , sep should use "^"
					fmt.Printf("\n Warning: Your CW capture file <%s> had unsupported ProSign format\n          character \"^\" (i.e. ^ ).", captureFile)
					fmt.Printf("\n\n          If that is correct for your ProSigns, use a carat \"^\"\n          between the file names (i.e.-send=%scapture.txt^practice.txt).\n", gchalk.Yellow("C:"))
					os.Exit(88)
				}
			}

			userGroupsCompare = strings.Fields(tmp)
			userGroupsCompareNoSpace = strings.ReplaceAll(tmp, " ", "")

		} else {
			panic("Got bad file response: report program error.")
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
		warningMsg = fmt.Sprintf("\n Note: The practice text file <%s> had <%d> groups.\n       Your CW capture file <%s> had <%d>.\n       Only the first <%d> groups will be checked.\n\n Your accuracy score is limited to checked groups!\n", practiceFile, sendLen, captureFile, userLen, userLen)
		maxIndex = userLen
	} else if userLen > sendLen {
		warningMsg = fmt.Sprintf("\n Note: Your CW capture file <%s> had too many groups <%d>.\n       Your practice text file <%s> only had <%d>.\n       Only the first <%d> groups will be checked.\n\n Your accuracy score is limited to checked groups!\n", captureFile, userLen, practiceFile, sendLen, sendLen)
		maxIndex = sendLen
	} else {
		maxIndex = sendLen
	}

	fmt.Printf(`
Levenshtein             Practice Text                      CW Capture
  Distance                 Groups                            Groups
===========    ==============================    ==============================
`)

	// compare char by char the MCPT group to user Group
	var lineCnt int

	for index, sgChar := range sendGroupsCompare {
		var out string
		var missed string
		var missedRaw string
		var tmpChar string
		var Index int

		// get users groups
		ugChar := userGroupsCompare[index]

		// diff strings
		totalChars += len(sgChar)

		// trim length if excessive
		levErrors := levenshtein(sgChar, ugChar)
		// print the table

		fmt.Printf("   %2d          ", levErrors) // first to keep alignment
		if levErrors > 0 {
			// so the string do NOT agree
			for ; Index < len(sgChar) && Index < len(ugChar); Index++ {

				if sgChar[Index] != ugChar[Index] {
					// mismatch - color bad data
					// col3
					tmpChar = char2psReplacer.Replace(string(ugChar[Index]))

					if tmpChar == "*" {
						out += gchalk.BrightMagenta(tmpChar)
					} else {
						out += colorError(tmpChar)
					}
				} else {
					// both matched good!
					// col3
					out += char2psReplacer.Replace(string(sgChar[Index]))
					totalCorrect++
				}
			}

			// one array is done
			if len(ugChar[Index:]) >= 1 && Index+1 > len(sgChar) { // walked off the sgChar
				// extra sent by user
				// but there might also be invaild as *
				tmpChar = char2psReplacer.Replace(ugChar[Index:])
				tmpChar = out + colorExtra(tmpChar)
				extra = true // once set, forever set
			} else if len(sgChar[Index:]) >= 1 {
				// more in column 2
				missedRaw = char2psReplacer.Replace(sgChar[Index:])
				//missed = colorMiss(char2psReplacer.Replace(sgChar[Index:]))
				missed = colorMiss(missedRaw)
				miss = true // once set, forever set
			}

		} else {
			totalCorrect += len(sgChar)
			out = char2psReplacer.Replace(ugChar)
		}

		sgCharToPrint := char2psReplacer.Replace(sgChar[:Index])
		if missed != "" {
			var padded string

			olen := len(sgCharToPrint) + len(missedRaw) - strings.Count(sgCharToPrint, "\u00D8")

			if olen < 30 {
				padded = sgCharToPrint + missed + strings.Repeat(" ", 30-olen)
			} else {
				padded = sgCharToPrint + missed
			}

			fmt.Printf("%-30s    %-30s\n", padded, out) // missed stuff from col 1
			missed = ""
		} else if extra {
			fmt.Printf("%-30s    %-30s\n", char2psReplacer.Replace(sgChar), tmpChar) // extra stuff in second col
			extra = false
			extraForever = true
		} else {
			// matched OR errors
			fmt.Printf("%-30s    %-30s\n", char2psReplacer.Replace(sgChar), out)
		}

		lineCnt++
		if lineCnt >= 10 {
			fmt.Printf("%s\n", strings.Repeat("-", 79))
			lineCnt = 0
		}

		out = ""
		tmpChar = ""
		index++

		if index == maxIndex {
			break
		}

	}

	if totalChars == 0 {
		fmt.Printf("\n\n Error: The practice file <%s> is empty.\n", practiceFile)
		os.Exit(1)
	}

	if totalCorrect == totalChars {
		grn := gchalk.BrightGreen("100 %")
		fmt.Printf("\n\n Accuracy: %s\n", grn)
	} else {
		fmt.Printf("\n\n Accuracy: %3.2f %%\n", float32(totalCorrect)*100.0/float32(totalChars))
	}

	if warningMsg != "" {
		fmt.Printf("%s", warningMsg)
	}

	if miss {
		fmt.Printf("\n Warning: You %s sending some characters, (column 2) in %s (ProSign counts as 1).\n", colorMiss("missed"), colorMiss("blue"))
		fmt.Printf("\n          If your CW capture groups following the %s characters are all %s,\n          you MAY have split a group with an extra space.", colorMiss("missed"), colorError("errors"))
		fmt.Printf("\n          The space should be just before practice group turned %s and your CW capture groups turned %s.\n", colorMiss("blue"), colorError("red"))
		fmt.Printf("\n          Edit your CW capture file <%s>, fix the space error, and rerun.\n", captureFile)
	}

	if extraForever {
		fmt.Printf("\n Warning: You sent some %s characters, (column 3) in %s.\n", colorExtra("extra"), colorExtra("green"))
		fmt.Printf("\n          Or you missed a space and combined two groups.")
		fmt.Printf("\n          The space should be just before your CW capture groups all turned %s.\n", colorError("red"))
		fmt.Printf("\n          Edit your capture file <%s>, fix the space error, and rerun.\n", captureFile)
	}

	fmt.Printf("\n\n Note: INVALID morse characters (or ProSigns) are shown as asterisks \"%s%s%s\".\n", gchalk.WithBgBlue().BrightCyan("*"), gchalk.BrightMagenta("*"), gchalk.BrightGreen("*"))

	// look at timing compare
	if totalCorrect != totalChars && miss || extraForever {
		var reply string
		fmt.Printf("\n Do you want to see the captured text from a timing perspective? (y or n): ")

		fmt.Scanf("%s", &reply)
		if reply != "" && reply[0] == byte('y') {
			timingCheck(sendGroupsCompareNoSpace, userGroupsCompareNoSpace)
		}
	}

	return 0
}

// see if the file was the MCPT generated file, or from the users software
func determineFile(path string) (byte, []byte, string) {
	gotMCPT := false

	// see if its the MCPT file
	if strings.HasPrefix(path, "c:") || strings.HasPrefix(path, "C:") {
		// the users cw capture groups
		path = strings.TrimPrefix(path, "c:")
		path = strings.TrimPrefix(path, "C:")
	} else {
		// the MCPT file or practice text  file
		gotMCPT = true
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("\n Error: reading file <%s>. %v\n", path, err)
		os.Exit(1)
	}

	b = bytes.ToUpper(b) // in case source was not -send

	if gotMCPT {
		// the MCPT or practice text file
		return byte('m'), b, path
	} else {
		// the users cw capture groups
		return byte('u'), b, path
	}
}

// compute the levenshtein error distance of the two groups
// str1 is what MCPT generated, str2 is what the user sent
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

/*
** chunk and compare just for errors to see if issue was spacing
 */

func timingCheck(practice string, capture string) {
	var strLen int = 60 // how many chars will we display
	var out string
	var captureSection string
	var practiceSection string

	// need to compare to the length of shortest
	var minLen int = len(practice)
	if len(capture) < minLen {
		minLen = len(capture)
	}

	for len(capture) > 0 {
		if len(practice) == 0 {
			break
		} else {
			if len(practice) >= strLen {
				practiceSection = practice[0:strLen]
				captureSection = capture[0:strLen]
				practice = practice[strLen:]
				capture = capture[strLen:]
			} else {
				practiceSection = practice[0:]
				captureSection = capture[0:]
				practice = ""
				capture = ""
			}
		}

		for i, pChar := range practiceSection {
			if i >= len(captureSection) {
				break
			}

			if byte(pChar) != captureSection[i] {
				tmpChar := char2psReplacer.Replace(string(captureSection[i]))

				if tmpChar == "*" {
					out += gchalk.BrightMagenta(tmpChar)
				} else {
					out += gchalk.BrightRed(tmpChar)
				}
			} else {
				// both matched good!
				out += char2psReplacer.Replace(string(pChar))
			}
		}

		fmt.Printf("\npractice: %s", char2psReplacer.Replace(practiceSection))
		fmt.Printf("\n capture: %s\n", out)

		out = ""
	}

	return
}
