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
	case "<AR>", "<AS>", "<BT>", "<KA>", "<HH>", "<SK>", "<VA>", "<SN>", "<VE>", "<DU>", "<SOS>":
		return true
	case "^AR", "^AS", "^BT", "^KA", "^HH", "^SK", "^VA", "^SN", "^VE", "^DU", "^SOS": // for g4fon maybe others
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

	word := regexp.MustCompile("^<[A-Za-z][A-Za-z]>|^\\^[A-Za-z][A-Za-z]|^<SOS>")

	var myRune rune = 'p'

	for scanner.Scan() {
		ps = strings.TrimSpace(scanner.Text())
		if len(ps) == 0 || ps[0] == byte('#') {
			// comment line
			continue
		}
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

			if flagCG || flagMixedMode > 1 || flagpermute != "" {
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
