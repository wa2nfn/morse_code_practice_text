package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

/*
  read input file looking for AnyString that matches current tutor/lesson.
  First we will strip characters out of the possible characters to listen to.
*/

func readStringsFile(fp *os.File) {

	discarded := false

	if flagtext == "" {
		fmt.Printf("\nError: an text input file must be given to -textFile.\n")
		os.Exit(0)
	}

	// modify so -lesson can be used with -text
	if flaglesson != "0:0" {
		// user wants lessons to be used in -in=file
		flaginlist = flagcglist // swap meaning
	}

	if flaginlist == "" {

		if flaginlist == "" {
			fmt.Printf("\nError: inList can't be empty or nothing gets matched.\n")
			os.Exit(0)
		}
	}

	content, err := os.ReadFile(flagtext)
	if err != nil {
		fmt.Printf("\n%s For file name <%s>.\n", err, flagtext)
		os.Exit(1)
	}

	// to match what user wants
	//var s string
	flaginlist += "^<>" // for prosigns

	re := fmt.Sprintf(`^[%s]{%d,%d}$`, flaginlist, flagmin, flagmax)
	isInSet := regexp.MustCompile(re).MatchString

	ps := regexp.MustCompile(`^<[A-Za-z]{2}>$|^\^[A-Za-z]{2}$`)

	replacer := strings.NewReplacer(
		"x-ray", "xray",
		"u-turn", "uturn",
		"n't'", " not",
		"i'll", "i will",
		"i'd", "i would",
		"i'm", "i am",
		"y've", "y have",
		"u've", "u have",
		"i've", "i have",
		"'s", " is",
		"'d", " would",
		"'re", " are",
		"!", "", "#", "", "$", "", "%", "", "&", "", "*", "", "(", "", ")", "", "-", " ", "_", " ", "{", "", "}", "", "`", "", "'", "", ":", "", ";", "", "\"", "", "|", "")

	line := string(content)
	line = strings.ToLower(line)
	line = replacer.Replace(line)
	line = strings.ToUpper(line)

	for _, wd := range strings.Fields(line) {

		// input line pruned

		if isInSet(wd) {

			// if prosign check it or ignore it
			if len(wd) == 3 || len(wd) == 4 {
				if ps.MatchString(wd) {
					if !ckProsign(wd) {
						// its invalid so skip it
						continue
					}
				}
			}

			// reverse the string
			if flagreverse {
				wd = reverse(wd)
			}

			/*
			** always ordered
			 */
			wordArray = append(wordArray, wd)
		} else {
			discarded = true
		}

	}

	if len(wordArray) == 0 {
		fmt.Println("\nThere is nothing to output.\nVerify options (or defaults) are not overly restrictive (inLen, inList, lesson).\nVerify your input file is sufficiently populated with matchable text.\n")

		if discarded {
			fmt.Printf("Your input file DID have some text, but nothing matched your criteria.\n")
		}
		os.Exit(0)
	}

	if flagprosign != "" && len(proSign) >= 1 {
		replaceIndex := make(map[int]struct{})

		ll := len(wordArray) - len(proSign)
		if ll <= 0 {
			fmt.Printf("\n Error: matched words <%d> must be > number of ProSigns in file <%s>.\n", len(wordArray), flagprosign)
			os.Exit(56)
		}

		// trim and then pad with proSign
		wordArray = wordArray[0:ll]

		for len(proSign)*100/len(wordArray) <= 15 {
			proSign = append(proSign, proSign...)
		}

		wordArray = append(wordArray, proSign...)

		for i := 0; i < len(proSign); {
			rand := rng.Intn(len(wordArray))
			if _, ok := replaceIndex[rand]; ok != true {
				replaceIndex[rand] = struct{}{}
				i++
			}
		}

		// now do the substitutions
		j := 0
		for index := range replaceIndex {
			temp := append([]string{proSign[j]}, wordArray[index:]...)
			wordArray = append(wordArray[:index], temp...)
			j++
		}
	}

	ct := len(wordArray)
	if ct < flagnum {
		ct = flagnum - ct
		// we need to append more words in order
		for i := 0; i < ct; i++ {
			wordArray = append(wordArray, wordArray[i])
		}
	}

	if flagordered == false {
		rand.Seed(time.Now().UTC().UnixNano())
		rand.Shuffle(len(wordArray), func(i, j int) { wordArray[i], wordArray[j] = wordArray[j], wordArray[i] })
	}

	// trim
	if flagnum > 0 {
		// if zero we limit to what we found
		wordArray = wordArray[0:flagnum]
	}
}
