package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

//
// permute - pseudo permutations of "lesson" characters
//
func permute(mode string, fp *os.File) {
	var cMap = make(map[string]struct{})
	var tMap = make(map[string]struct{})
	var numToPrt = flagnum
	var strBuf = ""
	cnt := 0
	MAX := 10000

	/*
		if len(flagcglist) == 0 {
			flagcglist = kochChars
		}
	*/

	if len(flagcglist) < 2 {
		fmt.Println("\nError: input string too short, increase <lesson> for your <tutor>.\n")
		os.Exit(9)
	}

	// see if we need to add the prosign runes
	if ps2runeMap != nil {
		for _, val := range ps2runeMap {
			flagcglist += string(val)
		}
	}

	// randomize the cglist
	var chars = []rune(flagcglist)
	var chars2 = []rune(flagcglist)
	for i, val := range flagcglist {
		chars[i] = val
		chars2[i] = val
	}

	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(flagcglist), func(i, j int) { chars[i], chars[j] = chars[j], chars[i] })
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(flagcglist), func(i, j int) { chars2[i], chars2[j] = chars2[j], chars2[i] })

	// make pairs
	if cMap != nil && mode != "t" {

		for _, H := range chars {
			for _, L := range chars2 {
				pair := string(H) + string(L)
				cMap[pair] = struct{}{}
				cnt++
			}

			if len(cMap) >= flagnum {
				numToPrt = flagnum
				break
			} else if cnt >= MAX {
				numToPrt = MAX
				break
			}
		}

		if mode == "p" {
			cnt = 0
			for out := range cMap {
				if cnt >= numToPrt {
					break
				}

				strBuf += string(out) + " "
			}

			if flagtutor == "LICWB2" || flagtutor == "B2" {
				strBuf = strings.ReplaceAll(strBuf, "$", " BK ")
				strBuf = strings.ReplaceAll(strBuf, "  ", " ")
			}

			printStrBuf(convertRunes(strBuf), fp)
			return
		}
	}

	// triples
	if tMap != nil {
		for _, H := range chars {
			for _, M := range chars2 {
				for _, L := range flagcglist {
					triple := string(H) + string(M) + string(L)
					tMap[triple] = struct{}{}
					cnt++
				}

			}
		}

		if len(tMap) >= flagnum {
			numToPrt = flagnum
		} else if cnt >= MAX {
			numToPrt = MAX
		}

		cnt = 0
		if mode == "t" {
			for out := range tMap {
				if cnt >= numToPrt {
					break
				}

				strBuf += string(out) + " "
				cnt++
			}

			printStrBuf(convertRunes(strBuf), fp)
			return
		} else {
			// mode = b
			for p := range cMap {
				tMap[p] = struct{}{}
			}

			// now print b
			for out := range tMap {
				if cnt >= numToPrt {
					break
				}

				strBuf += string(out) + " "
				cnt++
			}

			printStrBuf(convertRunes(strBuf), fp)
			return
		}
	}
}

func convertRunes(strBuf string) string {
	var out string

	if rune2psMap != nil {

		for _, s := range strBuf {
			if rune2psMap[s] != "" {
				out += string(rune2psMap[s])
			} else {
				out += string(s)
			}
		}
		return out
	}

	return strBuf
}
