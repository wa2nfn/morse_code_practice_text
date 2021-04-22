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
	char2psReplacer = strings.NewReplacer("a", "<AS>", "b", "<AR>", "c", "<BT>", "d", "<KA>", "e", "<HH>", "f", "<SK>", "g", "<BK>", "h", "<AA>", "i", "<CT>", "j", "<KN>", "k", "<VA>", "l", "<SN>", "0", "\u00D8")

	ps2charReplacer = strings.NewReplacer("<AS>", "a", "<AR>", "b", "<BT>", "c", "<KA>", "d", "<HH>", "e", "<SK>", "f", "<BK>", "g", "<AA>", "h", "<CT>", "i", "<KN>", "j", "<VA>", "k", "<SN>", "l", "\u00D8", "0","<err>","*")

	MCPTps2charReplacer = strings.NewReplacer("<AS>", "a", "<AR>", "b", "<BT>", "c", "<KA>", "d", "<HH>", "e", "<SK>", "f", "<BK>", "g", "<AA>", "h", "<CT>", "i", "<KN>", "j", "<VA>", "k", "<SN>", "l", "\u00D8", "0")

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
		char2psReplacer = strings.NewReplacer("a", "^AS", "b", "^AR", "c", "^BT", "d", "^KA", "e", "^HH", "f", "^SK", "g", "^BK", "h", "^AA", "i", "^CT", "j", "^KN", "k", "^VA", "l", "^SN", "0", "\u00D8","<err>","*")
		ps2charReplacer = strings.NewReplacer("^AS", "a", "^AR", "b", "^BT", "c", "^KA", "d", "^HH", "e", "^SK", "f", "^BK", "g", "^AA", "h", "^CT", "i", "^KN", "j", "^VA", "k", "^SN", "l", "\u00D8", "0")
		gotCarat = true
	}

	path = strings.Split(flagsendcheck, sep)

	if len(path) == 0 || len(path) != 2 {
		fmt.Printf("\n Error: Option <-sendCheck> requires 2 file names, in format like: file1,file2 (if Prosigns have <XX> format\nor file1^file2, if ProSigns have ^XX format.\n")
		os.Exit(1)
	}

	var errVal int

	errVal = readLines(path)

	switch errVal {
	case 1:
		fmt.Printf("\n Error: One file MUST be the captured text from your morse sending.\n")
	case 2:
		fmt.Printf("\n Error: One file MUST be a MCPT generated file of practice material.\n")
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
	var maxIndex int
	var totalCorrect int
	var totalChars int
	var warningMsg string
	var colorExtra = gchalk.BrightGreen
	var colorError = gchalk.BrightRed
	var colorMiss = gchalk.BrightYellow
	var miss bool
	var extra bool
	var extraForever bool
	m := regexp.MustCompile(`<[\-.]+>`)  // find <..--..> errors
	hh := regexp.MustCompile(`<.{8,8}>`) // find <HH> as code

	gchalk.SetLevel(gchalk.LevelAnsi16m)

	// do both files
	for fIndex := 0; fIndex <= 1; fIndex++ {
		// determine which file we have
		whoIsIt, b := determineFile(path[fIndex])
		// b is already in UC

		if whoIsIt == 'm' {
			// process MCPT file
			gotMCPT = true
			sendGroupsCompare = strings.Fields(MCPTps2charReplacer.Replace(string(b)))
			continue
		} else if whoIsIt == 'u' {
			// process User file
			gotUser = true

			// look for <HH>  errors first
			if m.Match(b) {
				b = hh.ReplaceAll(b, []byte("e"))
			}

			// look for <.......> errors first
			if m.Match(b) {
				b = m.ReplaceAll(b, []byte("*"))
			}

			if gotCarat {
				// got carat so PS needs looklike ^BT NOT <BT>
				if bytes.ContainsAny(b, "<>") {
					// used ^ sep should use ","
					fmt.Printf("\n Warning: Your file <%s> had unsupported ProSign format\n          characters \"<>\" (i.e. <BT>).", path[fIndex])
					fmt.Printf("\n\n          If those are correct for your ProSigns, use a comma \",\"\n          between the file names (i.e.-send=file1,file2).\n")

					os.Exit(88)
				}
			} else {
				// have , so need <>  ie <BT> NOT ^BT
				if bytes.ContainsAny(b, "^") {
					// used , sep should use "^"
					fmt.Printf("\n Warning: Your file <%s> had unsupported ProSign format\n          character \"^\" (i.e. ^BT ).", path[fIndex])
					fmt.Printf("\n\n          If that is correct for your ProSigns, use a carat \"^\"\n          between the file names (i.e.-send=file1^file2).\n")
					os.Exit(88)
				}
			}
			// convert any NL to space
			b = bytes.ReplaceAll(b, []byte("\n"), []byte(" "))

			var tmp = ps2charReplacer.Replace(string(b))
			if gotCarat == false && (strings.Contains(tmp, "<") || strings.Contains(tmp, ">")) {
				fmt.Printf(" Warning: Your file contains unsupported ProSigns or \"< or >\",\n          they will add to errors.\n")
			} else if gotCarat == true && strings.Contains(tmp, "^") {
				fmt.Printf(" Warning: Your file contains unsupported ProSigns or \"^\",          they will add to errors.\n")
			}

			userGroupsCompare = strings.Fields(ps2charReplacer.Replace(string(b)))
		} else {
			panic ( "Got back bad file response")
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
			warningMsg = fmt.Sprintf("\n Note: The MCPT file had <%d> groups, your file had <%d>.\n Only the first <%d> will be checked.\n\n Your accuracy score is limited to checked groups!\n", sendLen, userLen, userLen)
			maxIndex = userLen
		} else if userLen > sendLen {
			warningMsg = fmt.Sprintf("\n Note: Your file had too many groups <%d>, the MCPT file only had <%d>.\n Only the first <%d> will be checked.\n\n Your accuracy score is limited to checked groups!\n", userLen, sendLen, sendLen)
			maxIndex = sendLen
		} else {
			maxIndex = sendLen
		}

		fmt.Printf(`
Levenshtein             MCPT Created                       Your Sent
  Distance                 Group                             Group
===========    ==============================    ==============================
`)

		// compare char by char the MCPT group to user Group
		for index, sgChar := range sendGroupsCompare {
			var out string
			var missed string
			var tmpChar string
			var Index int

			// get users groups
			ugChar := userGroupsCompare[index]

			// diff strings
			totalChars += len(sgChar)

			// trim length if excessive
			levErrors := levenshtein(sgChar, ugChar)
			// print the table

			fmt.Printf("   %2d    ", levErrors) // first to keep alignment
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

				if Index+1 == len(sgChar) && Index+1 == len(ugChar) {
					break // finished
				}

				// one array is done
				if len(ugChar[Index:]) >= 1 && Index+1 > len(sgChar) { // walked off the sgChar
					// extra sent by user
					tmpChar = char2psReplacer.Replace(ugChar[Index:])
					tmpChar = out + colorExtra(tmpChar)
					extra = true // once set, forever set
				} else if len(sgChar[Index:]) >= 1 {
					// more in system column, user missed some
					missed = colorMiss(char2psReplacer.Replace(sgChar[Index:]))
					miss = true // once set, forever set
				}

			} else {
				totalCorrect += len(sgChar)
				out = char2psReplacer.Replace(ugChar)
			}

			sgChar = char2psReplacer.Replace(sgChar)
			if missed != "" {
				olen := len(sgChar)
				padded := sgChar[:Index] + missed + strings.Repeat(" ", 30 - olen)
				fmt.Printf("      %-30s    %-30s", padded, ugChar) // missed stuff from col 1
				missed = ""
			} else if extra {
				fmt.Printf("      %-30s    %-30s", sgChar, tmpChar) // extra stuff in second col
				extra = false
				extraForever = true
			} else {
				// matched OR errors
				fmt.Printf("      %-30s    %-30s", sgChar, out)
			}

			fmt.Println()

			out = ""
			tmpChar = ""
			index++

			if index == maxIndex {
				break
			}

		}

		if totalChars == 0 {
			fmt.Printf("\n\n Error: The MCPT file is empty.\n")
			os.Exit(1)
		}

		if totalCorrect == totalChars {
			grn := gchalk.BrightGreen("100%")
			fmt.Printf("\n\n Accuracy: %s\n\n", grn)
		} else {
			fmt.Printf("\n\n Accuracy: %3.2f%%\n\n", float32(totalCorrect)*100.0/float32(totalChars))
		}

		if warningMsg != "" {
			fmt.Printf("%s", warningMsg)
		}

		if miss {
			fmt.Printf("\n Warning: You %s sending some characters (ProSign counts as 1).\n",colorMiss("missed"))
			fmt.Printf("\n          If your sent groups following the %s characters are all (%s),\n          you may have split a group with an extra space.", colorMiss("missed"),colorError("errors"))
			fmt.Printf("\n          The space should be just before your sent groups turned %s.\n", colorError("red"))
			fmt.Printf("\n          Edit your sent file, fix the space error, and rerun.\n")
		}

		if extraForever {
			fmt.Printf("\n Warning: You sent some %s characters.\n",colorExtra("extra"))
			fmt.Printf("\n          Or you missed a space and combined two groups.")
			fmt.Printf("\n          The space should be just before your sent groups all turned %s.\n", colorError("red"))
			fmt.Printf("\n          Edit your sent file, fix the space error, and rerun.")
		}

	}
	return 0 
}

// see if the file was the MCPT generated file, or from the users software
func determineFile(path string) (byte, []byte) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("\n Error: reading file <%s>. %v\n", path, err)
		os.Exit(1)
	}

	b = bytes.ToUpper(b) // in case source was not -send

	// see if its the MCPT file
	if bytes.Contains(b, []byte("\u0008")) {
		// the MCPT file
		b = bytes.ReplaceAll(b, []byte("\u0008"), []byte(""))
		return byte('m'), b
	} else {
		// the users send groups
		return byte('u'), b
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
