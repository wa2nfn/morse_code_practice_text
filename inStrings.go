package main

import (
	"os"
	"strings"
	"bufio"
	"regexp"
	"math/rand"
	"fmt"
)


/*
  read input file looking for AnyString that matches current tutor/lesson.
  First we will strip characters out of the possible characters to listen to.
  (i.e. we will NOT handle the internationals since they are NOT in any of the lessons)
 */

func readStringsFile(localSkipFlag bool, localSkipCount int, fp *os.File) {

	discarded := false

	if flagtext == "" {
		fmt.Printf("\nError: an input file must be given to -text.\n")
		os.Exit(0)
	}

	if flaginlist == "" {
		fmt.Printf("\nError: inlist can't be empty or nothing gets matched.\n")
		os.Exit(0)
	}

	file, err := os.Open(flagtext)
	if err != nil {
		fmt.Printf("\n%s File name <%s>.\n", err, flagtext)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// to match what user wants
	var s string
	flaginlist += "^<>" // for prosigns

	re := fmt.Sprintf(`^[%s]{%d,%d}$`, flaginlist, flagmin, flagmax)
	isInSet := regexp.MustCompile(re).MatchString

	ps := regexp.MustCompile(`^<[A-Za-z]{2}>$|^\^[A-Za-z]{2}$`)
	replacer := strings.NewReplacer("!","","#","","$","","%","","&","","*","","(","",")","","-"," ","_"," ","{","","}","","`","","'","",":","",";","","\"","","|","")

	for scanner.Scan() {
		// first prune chars we don't want
		s = replacer.Replace(scanner.Text())
		// now s is the input line pruned
		// first way to split the string on spaces

		s = strings.ToUpper(s)

		textWords := strings.FieldsFunc(s, func(r rune) bool {
			if r == ' ' {
				return true
			}
			return false
		})


		for index := 0; index < len(textWords) && index <= flagnum; index++ {
			// tokens now a string of space separated characters

			tmpWord := textWords[index]

			if isInSet(tmpWord) {

				// skip only viable matching words
				if localSkipFlag {
					if localSkipCount > 0 {
						localSkipCount--
						continue
					} else {
						localSkipFlag = false
					}
				}


				// if prosign check it or ignore it
				if len(tmpWord) == 3 || len(tmpWord) == 4 {
					if ps.MatchString(tmpWord) {
						if !ckProsign(tmpWord) {
							// its invalid so skip it
							continue
						}
					}
				}

				// reverse the string
				if flagreverse {
					tmpWord = reverse(tmpWord)
				}

				/*
				** always ordered 
				 */
					wordArray = append(wordArray, tmpWord)
			} else {
				discarded = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	if len(wordArray) == 0 {
		fmt.Println("\nThere is nothing to output.\nVerify options (or defaults) are not overly restrictive (inlen, inlist, lesson).\nVerify your input file is sufficiently populated with matchable text.\n")

		if discarded {
			fmt.Printf("Your input file DID have some text, but nothing matched your criteria.\n")
			if localSkipFlag {
				fmt.Printf("Your -skip X option maybe too aggressive.\n")
			}
		}
		os.Exit(0)
	}

	if flagprosign != "" && len(proSign) >= 1 {
		replaceIndex := make(map[int]struct{})

		ll := len(wordArray) - len(proSign)
		if ll <= 0 {
			fmt.Printf("\nError: matched words <%d> must be > number of proSigns in file <%s>.\n", len(wordArray), flagprosign)
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

	if flagNR == false {
		rand.Shuffle(len(wordArray),func(i,j int){wordArray[i], wordArray[j] = wordArray[j], wordArray[i]})
	}

	if len(wordArray) > flagnum {
		wordArray = wordArray[0:flagnum]
	}
}
