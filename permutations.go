package main

import (
	"fmt"
	"os"
)

//
// permute - pseudo permutations of "lesson" characters
//
func permute(mode string, fp *os.File) {
	var cMap = make(map[string]struct{})
	var tMap = make(map[string]struct{})
	var numToPrt = flagnum
	var str = flagcglist
	var strBuf = ""
	length := len(str)
	cnt := 0

	if len(flagcglist) == 0 {
		flagcglist = kochChars
	}

	if length < 2 {
		fmt.Println("\nError: input string to short, increase <lesson> for your <tutor>.\n")
		os.Exit(9)
	}

	fullstr := str
	for str != "" {
		for _, c := range str {
			for _, full := range fullstr {
				makeTuple(c, full, mode, cMap, tMap)
			}
		}

		if len(str) >= 1 {
			str = str[1:]
		}
	}

	// now do output
	if cMap != nil && mode != "t" {
		// if both, let user pick any length
		if mode == "p" && len(cMap) >= flagnum {
			numToPrt = len(cMap)
		}

		for cnt < numToPrt {
			// print tuples, the triples unless "b" print all
			for key := range cMap {
				strBuf += key + " "
				cnt++
				if cnt >= numToPrt {
					break
				}
			}
		}
	}

	if tMap != nil && mode == "t" {
		if len(tMap) >= flagnum {
			numToPrt = len(tMap)
		}

		// print triples
		for cnt < numToPrt {
			// print triples
			for key := range tMap {
				strBuf += key + " "
				cnt++
				if cnt >= numToPrt {
					break
				}
			}
		}
	}

	printStrBuf(strBuf, fp)
}

func makeTuple(low rune, high rune, m string, cMap map[string]struct{}, tMap map[string]struct{}) {
	r := []rune{}
	r = append(r, low)
	r = append(r, high)
	if m != "t" {
		cMap[string(r)] = struct{}{}
	}

	r[0], r[1] = r[1], r[0]
	if m != "t" {
		cMap[string(r)] = struct{}{}
	}

	if m != "p" && low != high {
		r = append(r, high)
		if m == "t" {
			tMap[string(r)] = struct{}{}
		} else if m == "b" {
			cMap[string(r)] = struct{}{}
		}
	}
}
