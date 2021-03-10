package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strings"
)

/*
** added for Ordered write of input "NR"
 */

// read input file and create words. vs do code groups
func readFileMode(localSkipFlag bool, localSkipCount int, fp *os.File) {

	discarded := false

	if flaginput == "" {
		fmt.Printf("\nError: an input file must be given to -in, unless -codeGroups is used.\n")
		os.Exit(0)
	}

	if flaginlist == "" {
		fmt.Printf("\nError: inlist can't be empty or nothing gets matched.\n")
		os.Exit(0)
	}

	file, err := os.Open(flaginput)
	if err != nil {
		fmt.Printf("\n%s File name <%s>.\n", err, flaginput)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// to match what user wants
	var trimChars string
	var s string

	// added to fix regex to have unicode for international
	if strings.Contains(flaginlist, "%") {
		flaginlist = strings.ReplaceAll(flaginlist, "%C0", "\u00C0")
		flaginlist = strings.ReplaceAll(flaginlist, "%E0", "\u00E0")
		flaginlist = strings.ReplaceAll(flaginlist, "%C4", "\u00C4")
		flaginlist = strings.ReplaceAll(flaginlist, "%E4", "\u00E4")
		flaginlist = strings.ReplaceAll(flaginlist, "%C9", "\u00C9")
		flaginlist = strings.ReplaceAll(flaginlist, "%E9", "\u00E9")
		flaginlist = strings.ReplaceAll(flaginlist, "%C8", "\u00C8")
		flaginlist = strings.ReplaceAll(flaginlist, "%E8", "\u00E8")
		flaginlist = strings.ReplaceAll(flaginlist, "%C7", "\u00C7")
		flaginlist = strings.ReplaceAll(flaginlist, "%E7", "\u00E7")
		flaginlist = strings.ReplaceAll(flaginlist, "%D1", "\u00D1")
		flaginlist = strings.ReplaceAll(flaginlist, "%F1", "\u00F1")
		flaginlist = strings.ReplaceAll(flaginlist, "%D6", "\u00D6")
		flaginlist = strings.ReplaceAll(flaginlist, "%F6", "\u00F6")
		flaginlist = strings.ReplaceAll(flaginlist, "%DC", "\u00DC")
		flaginlist = strings.ReplaceAll(flaginlist, "%FC", "\u00FC")
	}

	if strings.Contains(flaginlist, "?") {
		trimChars = ".\",!"
		flaginlist = strings.ReplaceAll(flaginlist, "?", "")
		re := fmt.Sprintf(`[%s]{%d,%d}\?{0,1}`, flaginlist, flagmin, flagmax)
		s = fmt.Sprintf(`^%s$|^%s[%s]*$^<[A-Za-z]{2}>$|^\^[A-Za-z]{2}$`, re, flaginlist, trimChars)
	} else {
		trimChars = ".\",!?"
		re := fmt.Sprintf(`[%s]{%d,%d}`, flaginlist, flagmin, flagmax)
		s = fmt.Sprintf(`^%s$|^%s[%s]*$^<[A-Za-z]{2}>$|^\^[A-Za-z]{2}$`, re, flaginlist, trimChars)
	}

	word := regexp.MustCompile(s)
	ps := regexp.MustCompile(`^<[A-Za-z]{2}>$|^\^[A-Za-z]{2}$`)

	if !flagNR {
		wordMap = make(map[string]struct{})
	}

	// read ProSigns
	if flagprosign != "" {
		psfile, err := os.Open(flagprosign)
		if err != nil {
			fmt.Printf("\n%s File name <%s>.\n", err, flagprosign)
			os.Exit(1)
		}

		// fill the map with prosigns
		doProSigns(psfile)
		psfile.Close()
	}

	for scanner.Scan() {
		// first way to split the string on spaces

		// lets preprocess for '_-
		var line string
		replacer := strings.NewReplacer("!", "", "#", "", "$", "", "%", "", "&", "", "*", "", "(", "", ")", "", "-", " ", "_", " ", "{", "", "}", "", "`", "", ":", "", ";", "", "'", "", "\"", "")
		line = replacer.Replace(scanner.Text())

		textWords := strings.FieldsFunc(line, func(r rune) bool {
			if r == ' ' {
				return true
			}
			return false
		})

		for index := 0; index < len(textWords); index++ {
			// every token is now a string of space separated characters
			tmpWord := strings.TrimRight(textWords[index], trimChars)
			tmpWord = strings.TrimLeft(tmpWord, "\"")
			tmpWord = strings.ToUpper(tmpWord)

			if word.MatchString(tmpWord) {

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
				** if -NR words are ordered so we store and retrieve from an array
				** else we use a map
				 */
				if flagNR {
					wordArray = append(wordArray, tmpWord)
				} else {
					// add to map if not there
					if _, ok := wordMap[tmpWord]; ok != true {
						wordMap[tmpWord] = struct{}{}
					}
				}
			} else {
				discarded = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	msg := "\nSorry there is nothing to output.\nMake sure the options (or defaults) are not to restrictive (inlen, inlist).\nVerify your input file is sufficiently populated with matchable text.\n"

	if flagNR {
		// ORDERED so use array
		// ok so either the MAP or ARRAY is populated from input

		if len(wordArray) == 0 {
			fmt.Println(msg)

			if discarded {
				fmt.Printf("Your input file DID have some text.\n")
				if localSkipFlag {
					fmt.Printf("Your -skip X option maybe too aggressive.\n")
				}
			}
			os.Exit(0)
		}

		ct := len(wordArray)
		if ct < flagnum {
			ct = flagnum - ct
			// we need to append more words in order
			for i := 0; i < ct; i++ {
				wordArray = append(wordArray, wordArray[i])
			}
		}

		// proSigns for NR = false done differently
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

		wordArray = wordArray[0:flagnum]
	} else {
		// RANDOM so use MAP to fill array
		if len(wordMap) == 0 {
			fmt.Println(msg)
			if discarded {
				fmt.Printf("Your input file DID have some text.\n")
				if localSkipFlag {
					fmt.Printf("Your -skip X option maybe too aggresive.\n")
				}
			}
			os.Exit(0)
		}

		// trim the saved wordMap to save time and memory later
		if len(wordMap) > flagnum {
			cntr := 0
			m := make(map[string]struct{})
			for v := range wordMap {
				if cntr == flagnum {
					wordMap = m
					break
				} else {
					m[v] = struct{}{}
					cntr++
				}
			}
			// swap maps
			wordMap = m

		}
		runtime.GC()

		// wordsMap has already gotten appropriate words
		// we just need to divide into length buckets
		if flaghead {

			mlw := map[int][]string{}

			keys := make([]int, maxWordLen, maxWordLen)
			for ind := flagmin; ind <= flagmax; ind++ {
				keys[ind] = ind
			}

			// Loop over map and append keys to empty slice.
			for key, _ := range mlw {
				keys = append(keys, key)
			}

			// put words in buckets by length
			for wd := range wordMap {
				l := len(wd)
				mlw[l] = append(mlw[l], wd)
			}

			wordMap = nil //recover space

			// now use the slice to find words to use
			for i := 0; i <= flagnum; {
				ind := 0

				for ind = flagmin; ind <= flagmax; ind++ {

					if mlw[ind] == nil {
						continue
					}

					rand := rng.Intn(len(mlw[ind]))
					wordArray = append(wordArray, (mlw[ind])[rand])

					i++

					if i == flagnum {
						break
					}
				}

				if flagDMmin >= 1 && (flagDR == false || (flagDR == true && flipFlop())) {
					buf := ""

					if flagDMmin == flagDMmax {
						for count := 0; count < flagDMmax; count++ {
							buf += delimiterSlice[rng.Intn(len(delimiterSlice))]
						}
					} else {
						for count := 0; count < (flagDMmin + rng.Intn(flagDMmax-flagDMmin+1)); count++ {
							buf += delimiterSlice[rng.Intn(len(delimiterSlice))]
						}
					}

					wordArray = append(wordArray, buf)
					buf = ""
				}
			}
			return
		}

		fillArray(fp)
	}
}

// ready to print the users practice word
func doOutput(words []string, fp *os.File) {
	strBuf := ""
	strOut := ""
	// for eb options
	firstSlowFast := true
	lastSpeed := flagLCWOlow
	lastSpeedEff := flagLCWOeff
	counter := 1
	sectionSize := 0
	LCWOspeeds := []int{}
	LCWOspeedsRepeat := []int{}
	ebslowcnt := 0
	ebfastcnt := 0
	ebinslow := false
	flagRANDOM := false
	flagSLOWFAST := false
	flagRAMP := false
	flagEFFRAMP := false
	flagREPEAT := false
	speedCount := 0
	anyLCWO := false
	var charSlice []rune

	if flagMixedMode > 0 {
		charSlice = buildCharSlice()
	}

	// header
	if flagheader != "" {
		strOut += fmt.Sprintf("%s\n", flagheader)
	}

	// in case we have mixedMode and speed change need to override speed in header
	if flagMixedMode > 1 && flagLCWOstep > 0 {
		strOut = fmt.Sprintf(" |e0 |w%d ", flagLCWOlow+flagLCWOstep)
	}

	// for runtime efficiency?
	if flagLCWOramp {
		flagRAMP = true
	}
	if flagLCWOefframp {
		flagEFFRAMP = true
	}
	if flagLCWOslow > 0 {
		flagSLOWFAST = true
		anyLCWO = true
	}
	if flagLCWOrandom {
		if flagLCWOnum > 1 && flagLCWOlow > 0 && flagLCWOstep >= 1 && !flagRAMP && !flagSLOWFAST && !flagEFFRAMP && !flagREPEAT {
			flagRANDOM = true
			anyLCWO = true
		} else {
			fmt.Printf("\nError: Invalid combination of LCWO options.\n")
			os.Exit(5)
		}
	}
	if flagLCWOrepeat >= 1 {
		flagREPEAT = true
		anyLCWO = true
	}

	if flagLCWOlow > 0 && flagLCWOstep > 0 && flagLCWOnum > 0 && (flagRAMP == false && flagRANDOM == false && flagREPEAT == false) {
		fmt.Printf("\nError: You're missing an LCWO option to indicate a feature: i.e. LCWO_random, LCWO_repeat, ...\n")
		os.Exit(5)
	}

	if anyLCWO {
		// ONLY if lcwo is used
		// since LCWO made eff=5 default in Convert text to CW in July 2020
		strOut += " |e0 "
	}

	///////////////////////////////////
	////// LCWO handling - intial setup
	///////////////////////////////////
	// seed array with LCWO_low, and fill as appropriate
	if flagLCWOnum >= 1 && flagLCWOstep > 0 {
		for i := 0; i < flagLCWOnum; i++ {
			spd := flagLCWOlow + (i * flagLCWOstep)
			LCWOspeeds = append(LCWOspeeds, spd)
		}
	}

	if flagLCWOrepeat >= 1 && flagLCWOstep > 0 {
		for i := 0; i < flagLCWOrepeat; i++ {
			LCWOspeedsRepeat = append(LCWOspeedsRepeat, flagLCWOlow+(i*flagLCWOstep))
		}
	}

	// LCWO flagRAMP how many words per ramp section
	if flagRAMP && !flagEFFRAMP {
		sectionSize = flagnum / flagLCWOnum
		if sectionSize < 1 {
			fmt.Printf("\nError: LCWO_num is too large for the -num value.\nThere would not be any words in each speed change section.\n")
			os.Exit(1)
		}

		lastSpeed = LCWOspeeds[0]

		if !flagREPEAT {
			if flagLCWOeff > 0 {
				strOut += fmt.Sprintf("|w%d |e%d ", LCWOspeeds[0], flagLCWOeff)
			} else {
				strOut += fmt.Sprintf("|w%d |e0 ", LCWOspeeds[0])
			}
			speedCount++
		}
	}

	// LCWO effective_ramp how many words per ramp section
	if flagEFFRAMP && !flagRAMP {
		sectionSize = flagnum / flagLCWOnum

		if sectionSize < 1 {
			fmt.Printf("\nError: LCWOnum is too large for the -num value.\nThere would not be any words in each speed change section.\n")
			os.Exit(1)
		}

		lastSpeedEff = flagLCWOeff
		strOut += fmt.Sprintf("|w%d |e%d ", flagLCWOlow, flagLCWOeff)
		counter = 0
	}

	////////////////////////////////////////
	/// setup done, now process input words
	/////////////////////////////////////////
	// select words from array, then lower high water mark so all words get used

	for index, wordOut := range words {

		//////////////
		// LCWO flagRANDOM
		//////////////
		if flagRANDOM {
			speed := LCWOspeeds[0]

			if len(LCWOspeeds) >= 1 {

				for {
					speed = LCWOspeeds[rng.Intn(len(LCWOspeeds))]

					if speed != lastSpeed {
						lastSpeed = speed
						break
					}
				}
			}

			if flagLCWOeff > 0 {
				strOut += fmt.Sprintf("|w%d |e%d ", speed, speed-effDelta)
			} else {
				strOut += fmt.Sprintf("|w%d ", speed)
			}
		}

		/////////////////
		/// LCWO FAST_SLOW
		/////////////////
		if flagSLOWFAST {
			s := flagLCWOlow

			if ebinslow {
				if ebslowcnt >= flagLCWOslow {
					s = flagLCWOlow + flagLCWOstep // now fast
					// slow words are done
					ebfastcnt = 0

					// keep eff same
					strOut += fmt.Sprintf("%s |w%d ", flagLCWOsf, s)
					ebinslow = false
				}
			} else {
				if ebfastcnt >= flagLCWOfast || firstSlowFast {
					firstSlowFast = false
					s := flagLCWOlow // now slow

					// fast words are done
					ebslowcnt = 0

					// set up slow section
					if index == 0 {
						if flagLCWOeff > 0 {
							strOut += fmt.Sprintf("|e%d |w%d ", flagLCWOeff, s)
						} else {
							strOut += fmt.Sprintf("|w%d ", s)
						}
					} else {
						strOut += fmt.Sprintf("%s |w%d ", flagLCWOfs, s)
					}
					ebinslow = true
				}
			}
		}

		// end raw word, and get back word to print
		wordOut, charSlice = prepWord(wordOut, lastSpeed, index, charSlice)

		///////////////////////////////////
		// LCWO CHECK FOR SPEED MARKERS
		///////////////////////////////////

		if flagRAMP {

			if counter >= sectionSize && speedCount < len(LCWOspeeds) {
				sf := ""
				if flagLCWOsf != "" {
					sf = " " + flagLCWOsf
				}

				if index+flagLCWOnum <= flagnum {
					if flagLCWOeff > 0 {
						strOut += fmt.Sprintf("%s%s |e%d |w%d ", wordOut, sf, LCWOspeeds[speedCount]-effDelta, LCWOspeeds[speedCount])
					} else {
						strOut += fmt.Sprintf("%s%s |w%d ", wordOut, sf, LCWOspeeds[speedCount])
					}
					wordOut = ""
					speedCount++
				}

				if speedCount < len(LCWOspeeds) {
					lastSpeed = LCWOspeeds[speedCount]
				}
				counter = 1
			} else {
				counter++
			}
		}

		////////////
		// flagEFFRAMP
		///////////
		if flagEFFRAMP {
			if counter == sectionSize {
				// ck if eff is going to over take word speed
				if lastSpeedEff+flagLCWOstep <= flagLCWOlow {
					// cap the eff speed
					lastSpeedEff += flagLCWOstep
				}
				strOut += fmt.Sprintf("%s |e%d ", wordOut, lastSpeedEff)
				counter = 1
				strOut += wordOut
				wordOut = ""
			} else {
				counter++
			}

		}

		//////////////
		/// flagSLOWFAST
		//////////////
		if flagSLOWFAST {
			if ebinslow {
				ebslowcnt++
			} else {
				ebfastcnt++
			}
		}

		// get back to individual words
		wordOut = strings.ReplaceAll(wordOut, "~", " ")

		// this is the processed word to be used
		strBuf += strOut + wordOut
		strOut = ""
	}

	printStrBuf(strBuf, fp)
}

/*
** take in a raw word from input file and tack on: prefix, suffix
** repeat if necessay,do mixedMode
 */
func prepWord(wordOut string, lastSpeed int, index int, charSlice []rune) (string, []rune) {
	strOut := ""
	rand := 3
	mustLen := len(flagmust)

	if flagrandom {
		if flagsufmin >= 1 || flagpremin >= 1 {
			// 0 - neither ix, 1 prefix,2 do suffix, 3 both
			rand = rng.Intn(4)
		}
	}

	// see if -must set, if so do one substitution now
	if mustLen > 0 {
		ll := 0
		bytearr := make([]byte, len(wordOut))
		for k, v := range wordOut {
			bytearr[k] = byte(v)
		}
		ll = rng.Intn(len(wordOut))
		bytearr[ll] = byte(flagmust[rng.Intn(mustLen)])
		wordOut = string(bytearr)
	}

	// end raw word, and get back word to print
	// do we need prefix?
	if flagpremin >= 1 && (rand == 3 || rand == 1) {
		wordOut = ixStr("p") + wordOut
	}

	// do we need a suffix or just a space
	if flagsufmin >= 1 && (rand == 3 || rand == 2) {
		wordOut += ixStr("s")
	}

	// text repeat!
	if flagrepeat > 0 {
		// we need to repeat
		wordOut += " "
		temp := wordOut

		for cnt := 1; cnt < flagrepeat; cnt++ {
			// wordOut is the word plus trailing space already
			wordOut += temp
		}
	}

	// LCWO_repeat
	if flagLCWOrepeat > 1 {

		for i := 0; i < flagLCWOrepeat; i++ {
			// if we ALSO have flagRAMP we must offset speed
			spd := lastSpeed + (i * flagLCWOstep)

			if flagLCWOeff > 0 {
				strOut += fmt.Sprintf("|w%d |e%d %s", spd, spd-effDelta, wordOut)
			} else {
				strOut += fmt.Sprintf("|w%d %s", spd, wordOut)
			}

			// speed offset
			if flagLCWOramp {
				spd += lastSpeed
			}
		}
	}

	/*
		// why? WDL old logic, can we remove it
			if flagLCWOramp && flagLCWOrepeat == 0 {
				strOut = wordOut
				wordOut = ""
			}
	*/

	// mixedMode put out code Group
	if flagMixedMode > 1 && (flagMMR == false || (flagMMR == true && flipFlop())) {
		g := []rune{}

		if index%flagMixedMode == 0 {
			// set slow speed
			if flagLCWOstep > 0 {
				strOut += fmt.Sprintf("|w%d ", flagLCWOlow)
			}
			if flagLCWOfs != "" {
				strOut += flagLCWOfs + " "
			}

			g, charSlice = makeSingleGroup(charSlice)

			strOut += string(g)
			if flagLCWOsf != "" {
				strOut += flagLCWOsf + " "
			}
			// set fast speed
			if flagLCWOstep > 0 {
				strOut += fmt.Sprintf("|w%d ", flagLCWOlow+flagLCWOstep)
			}
		}
	}

	// this means NOTHING was done to the word
	if strOut == "" {
		strOut = wordOut
	}

	// use delimiter if NOT headcopy
	if flagDMmin >= 1 && flaghead == false && (flagDR == false || (flagDR == true && flipFlop())) {
		if flagDMmin == flagDMmax {
			for count := 0; count < flagDMmax; count++ {
				strOut += delimiterSlice[rng.Intn(len(delimiterSlice))]
			}
		} else {
			for count := 0; count < (flagDMmin + rng.Intn(flagDMmax-flagDMmin+1)); count++ {
				strOut += delimiterSlice[rng.Intn(len(delimiterSlice))]
			}
		}

		strOut += " "
	}

	return strOut, charSlice
}

// returns a single random prefix/suffix from list to add to output word
func ixStr(ps string) string {
	retStr := ""

	if ps == "s" {
		// user wants a suffix
		if flagsufmin == flagsufmax {
			for count := 0; count < flagsufmax; count++ {
				retStr += string(flagSuflistRune[rng.Intn(len(flagSuflistRune))])
			}
		} else {
			for count := 0; count < (flagsufmin + rng.Intn(flagsufmax-flagsufmin+1)); count++ {
				retStr += string(flagSuflistRune[rng.Intn(len(flagSuflistRune))])
			}
		}
	} else {
		// user wants a prefix
		if flagpremin == flagpremax {
			for count := 0; count < flagpremax; count++ {
				retStr += string(flagPrelistRune[rng.Intn(len(flagPrelistRune))])
			}
		} else {
			for count := 0; count < (flagpremin + rng.Intn(flagpremax-flagpremin+1)); count++ {
				retStr += string(flagPrelistRune[rng.Intn(len(flagPrelistRune))])
			}
		}
	}

	return retStr
}

// fill the array from the word map but might need to stuff more values
func fillArray(fp *os.File) {
	var wordArray = make([]string, 0, flagnum)

	for key := range wordMap {
		// make first population of slice
		wordArray = append(wordArray, key)
	}

	// see if initial array satisfies the number of words the user wanted
	// if less, we will reuse words from map to grow the array (or slice)
	if !flagunique && len(wordArray) < flagnum {

		howShort := flagnum - len(wordArray)
		factor := flagnum / len(wordArray)
		factor-- // we have the original already

		// only does FULL maps
		var key string
		for ; factor > 0; factor-- {
			for key = range wordMap {
				wordArray = append(wordArray, key)
			}
		}

		// may still be a partial shortage
		howShort = flagnum - len(wordArray)
		i := 1
		for key := range wordMap {
			wordArray = append(wordArray, key)
			if i == howShort {
				break
			}
			i++
		}
	}

	// trash the map to conserve memory
	wordMap = nil

	if len(proSign) >= 1 {
		// enough words, now stuff in prosigns
		ll := len(wordArray) - len(proSign)
		if ll <= 0 {
			fmt.Printf("\nError: matched words <%d> must be > number of proSigns in file <%s>.\n", len(wordArray), flagprosign)
			os.Exit(54)
		}

		// trim and then pad with proSign
		wordArray = wordArray[0:ll]

		for len(proSign)*100/len(wordArray) <= 15 {
			proSign = append(proSign, proSign...)
		}

		// trim and then pad with proSign
		wordArray = append(wordArray, proSign...)
	}

	// shuffle array
	NwordArray := make([]string, 0, 0)
	in := 0
	for _, i := range rand.Perm(len(wordArray)) {
		NwordArray = append(NwordArray, wordArray[i])
		in++
	}

	wordArray = nil
	doOutput(NwordArray, fp)
}
