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
	"flag"
	"math"
)

var (
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
		doSendCheck()
	}
	os.Exit(1)
}

func doSendCheck() {
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
		/*
		// cannot use lc e or w, conflict with LCWO
		*/
		char2psReplacer = strings.NewReplacer(
			"(", "^AS", 
			")", "^AR", 
			"{", "^BT", 
			"}", "^KA", 
			"#", "^HH", 
			"[", "^SK", 
			"]", "^SN", 
			"0", "\u00D8",
			"$","^VE",
			"%","^DU",
			"*","^SOS")
		ps2charReplacer = strings.NewReplacer(
			"^AS", "(", 
			"^AR", ")", 
			"^BT", "{", 
			"^KA", "}", 
			"^HH", "#", 
			"^SK", "[", 
			"^SN", "]", 
			"\u00D8", "0",
			"+",")",
			"=","{",
			"^VE","$",
			"^DU","%",
			"^SOS","*")
		gotCarat = true
	}

	path = strings.Split(flagsendcheck, sep)

	if len(path) == 0 || len(path) != 2 {
		fmt.Printf("\n Error: Option <-sendCheck> requires 2 file names, in format like: file1,file2 (if ProSigns have <XX> format\n        or file1^file2, if ProSigns have ^XX format.\n")
		os.Exit(1)
	}

	readLines(path)
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
	if flagcgmax != flagcgmin && flagheadcopy == 0 {
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

	var myList string 

	if flag.Arg(0) != "" {
		myList = flag.Arg(0)
	}
	if len(flag.Args()) > 1 {
		fmt.Printf("\n Error: you can only have 1 string to the right of the option <-send=value>.\n")
		os.Exit(1)
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
		charSlice = append(charSlice, []rune("abcdfgh")...)
		gotDigit = true
	}
	if strings.Contains(flagsend, "7") {
		// cut numbers
		charSlice = append(charSlice, []rune("TAUV456BDN")...)
		gotDigit = true
	}
	if strings.Contains(flagsend, "8") {
		// cut numbers
		charSlice = append(charSlice, []rune("(){}#[]$%*")...)
		gotDigit = true
	}

	if myList != "" {
		tStr := ""
		// no op allows use of cglist
		tStr = ps2charReplacer.Replace(strings.ToUpper(myList))

		charSlice = append(charSlice, []rune(tStr)...)
		gotDigit = true
	}

	if gotDigit == false {
		fmt.Printf("\n Error: option <-send> must include at least one digit from 0-7.\n        Digit \"0\" is required, if it's the only digit and you are specifying\n        your own group list.\n        E.g. -send=0 \"EOMT5\"\n")
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
func readLines(path []string) {
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
	badPS := regexp.MustCompile(`<[^>]*>|\*`)  // find <..--..> errors and anything since supported PD already done
	hhErr := regexp.MustCompile(`<.{8,8}>`) // find <HH> as code

	gchalk.SetLevel(gchalk.LevelAnsi256)

	// do both files
	for fIndex := 0; fIndex <= 1; fIndex++ {

		b, err := ioutil.ReadFile(path[fIndex])
		if err != nil {
			fmt.Printf("\n %s reading file <%s>. %v\n", gchalk.Red("Error:"), path[fIndex], err)
			os.Exit(1)
		}

		// b is already in UC
		b = bytes.ToUpper(b) 
		// convert any NL to space
		b = bytes.ReplaceAll(b, []byte("\n"), []byte(" "))

		// practice file
		if fIndex == 0 {
			practiceFile = path[0]
			// good PS replace
			bStr := ps2charReplacer.Replace(string(b))

			if strings.ContainsAny(bStr, "<>^") {
				fmt.Printf("\n %s Your practice file <%s> contains character(s) \"<>^\" in addition to\n          those in supported ProSigns. This will add to your error count.\n\n", gchalk.Yellow("Warning:"), practiceFile)

				// look for < anything>
				if hhErr.MatchString(bStr) {
					bStr = hhErr.ReplaceAllString(bStr, "e")
				}
			}

			sendGroupsCompare = strings.Fields(bStr)
			sendGroupsCompareNoSpace = strings.ReplaceAll(bStr, " ", "")
			continue
		}

		// capture file
		if fIndex == 1 {
			// process User file
			captureFile = path[1]

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
					fmt.Printf("\n %s Your CW capture file <%s> had unsupported ProSign format\n          characters \"<>\" (i.e. <>).", gchalk.Yellow("Warning:"), captureFile)
					fmt.Printf("\n\n          If those are correct for your ProSigns, use a comma \",\"\n          between the file names (i.e.-send=%cature.txt,practice.txt).\n", gchalk.Yellow("C:"))

					os.Exit(88)
				}
			} else {
				// have , so need <>  ie <BT> NOT ^BT
				if strings.ContainsAny(tmp, "^") {
					// used , sep should use "^"
					fmt.Printf("\n %s Your CW capture file <%s> had unsupported ProSign format\n          character \"^\" (i.e. ^ ).", gchalk.Yellow("Warning:"), captureFile)
					fmt.Printf("\n\n          If that is correct for your ProSigns, use a carat \"^\"\n          between the file names (i.e.-send=practice.txt^capture.txt).\n")
					os.Exit(88)
				}
			}

			userGroupsCompare = strings.Fields(tmp)
			userGroupsCompareNoSpace = strings.ReplaceAll(tmp, " ", "")

		}
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
		fmt.Printf("\n\n %s The practice file <%s> is empty.\n", gchalk.Red("Error:"), practiceFile)
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
		fmt.Printf("\n %s You %s sending some characters, (column 2) in %s.\n          ProSigns count as 1 char.\n", gchalk.Yellow("Warning:"), colorMiss("missed"), colorMiss("blue"))
		fmt.Printf("\n          If your CW capture groups following the %s characters are all \n          %s, you MAY have split a group with an extra space.", colorMiss("missed"), colorError("errors"))
		fmt.Printf("\n          The space should be just before practice group turned %s and your CW\n          capture groups turned %s.\n", colorMiss("blue"), colorError("red"))
		fmt.Printf("\n          Edit your CW capture file <%s>, fix the space error, and rerun.\n", captureFile)
	}

	if extraForever {
		fmt.Printf("\n %s You sent some %s characters, (column 3) in %s.\n", gchalk.Yellow("Warning:"),colorExtra("extra"), colorExtra("green"))
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

	return
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
	var strLen int = 72 // how many chars will we display
	var out string
	var pLen = len(practice)
	var cLen = len(capture)
	var min = math.Min(float64(pLen),float64(cLen))

	strLen = int( math.Min(float64(strLen), float64(min)) )
	practice = practice[:strLen]
	capture = capture[:strLen]
	// now we can compare because we have the smallest length

	for i, pChar := range practice {

		if byte(pChar) != capture[i] {
			tmpChar := char2psReplacer.Replace(string(capture[i]))

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


	fmt.Printf("\n\nBoth practice and capture have been limited to 72 characters.\n")
	fmt.Printf("\npractice (below):\n\n%s", char2psReplacer.Replace(practice))
	fmt.Printf("\n%s\n\ncapture (above):", out)

	out = ""

	return
}
