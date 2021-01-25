package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

//
// ckProsign
// these are only valid prosigns
//
func ckProsign(ps string) bool {
	ps = strings.ToUpper(ps)
	switch ps {
	case "<AA>", "<AR>", "<AS>", "<BT>", "<CT>", "<KA>", "<HH>", "<KN>", "<SK>", "<SN>":
		return true
	case "^AA", "^AR", "^AS", "^BT", "^CT", "^KA", "^HH", "^KN", "^SK", "^SN": // for g4fon
		return true
	default:
		return false
	}
}

//
// process the file of prosigns, check their validity
//
func doProSigns(file *os.File) {
	ps := ""
	fs := ""

	scanner := bufio.NewScanner(file)

	word := regexp.MustCompile("^<[A-Za-z][A-Za-z]>|^\\^[A-Za-z][A-Za-z]")

	var myRune rune = '\u00A1'

	for scanner.Scan() {
		ps = strings.TrimSpace(scanner.Text())
		if ps == "" {
			continue
		}

		field := strings.Fields(ps)
		if field == nil {
			continue
		}

		if len(field[0]) < 3 || len(field[0]) > 5 {
			continue
		}
		fs = field[0]

		if word.MatchString(fs) {
			fs = strings.ToUpper(fs)

			if !ckProsign(fs) {
				break
			}

			if flagCG || flagMixedMode > 1 {
				// save in a pair of arrays
				ps2runeMap[fs] = myRune
				rune2psMap[myRune] = fs
				myRune++
			} else {
				// so in words
				proSign = append(proSign, fs)
			}

		} // ignore non matching ProSigns
	}
}
