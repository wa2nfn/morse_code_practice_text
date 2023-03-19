package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func licw() string {
	b1seed := "REATINPGSLCDHOFUWB"
	b2seed := "KMY59,QXV73?+sb16.ZJ/28@40"
	str := ""
	tmpSet := ""
	var max int
	var min int

	switch flaglesson {
	case "BC1":
		if flagtutor == "BC2C" || flagtutor == "BC2S" {
			fmt.Printf("\nError: lesson value incomplete with the tutor=<%s>, require lesson ID(s)\n", flagtutor)
			os.Exit(99)
		}
		return b1seed
	case "BC2":
		if flagtutor == "BC1C" || flagtutor == "BC1S" {
			fmt.Printf("\nError: lesson value <BC2> is invalid with the tutor=<%s>\n", flagtutor)
			os.Exit(99)
		}
		return b2seed
	case "BC1:BC2", "BC2:BC1":
		return b1seed + b2seed
	}

	b1seedDbl := "REATINPGSLCDHOFUWBREATINPGSLCDHOFUWB"
	b2seedDbl := "KMY59,QXV73?+sb16.ZJ/28@40KMY59,QXV73?+sb16.ZJ/28@40"

	// analyse min and max
	if strings.HasPrefix(flaglesson, "BC1:") {

		if flagtutor == "BC2C" || flagtutor == "BC2S" {
			flaglesson = strings.TrimPrefix(flaglesson, "BC1:")
			tmpSet = b1seed
			str = b2seedDbl
		} else {
			fmt.Printf("\nError: lesson value for this tutor can ONLY have\n       these formats: BC1, <ID>, <ID>:<ID>\n")
			os.Exit(99)
		}

	} else if strings.HasPrefix(flaglesson, "BC2:") {
		fmt.Printf("\nError: lesson value never contains BC2\n")
		fmt.Printf("lesson value for this tutor can ONLY have these formats: BC1:<ID>, BC1:<ID>:<ID>, <ID>, <ID>:<ID>\n")
		os.Exit(99)
	}

	// char mode, session done later
	if flagtutor == "BC1C" {
		str = b1seedDbl
	} else if flagtutor == "BC2C" {
		str = b2seedDbl
	}

	if strings.Contains(flaglesson, ":") {
		op := strings.Split(flaglesson, ":")

		if len(op) > 2 || len(op) == 0 {
			fmt.Printf("\nError: invalid format for option <lesson>\n")
			os.Exit(99)
		}

		Min, error := strconv.Atoi(op[0])

		if error != nil {
			fmt.Printf("\nError: invalid format for option <lesson>\n")
			os.Exit(99)
		}
		min = Min

		if len(op) == 2 && op[1] != "" {
			Max, error := strconv.Atoi(op[1])

			if error != nil {
				fmt.Printf("\nError: invalid format for option <lesson>\n")
				os.Exit(99)
			}
			max = Max
		}
	} else {
		// no :
		Min, error := strconv.Atoi(flaglesson)
		if error != nil {
			fmt.Printf("\nError: invalid format for option <lesson>\n")
			os.Exit(99)
		}
		min = Min
		max = min
	}

	if min < 0 || max < 0 {
		fmt.Printf("\nError: invalid value for option <lesson>, min and max <%s> must be >= 0.\n", flaglesson)
		os.Exit(99)
	}

	if max < min {
		fmt.Printf("\nError: invalid value for option <lesson>, max must be >= min <%s>.\n", flaglesson)
		os.Exit(99)
	}

	// different process for session vs char access tutor
	if flagtutor == "BC1C" || flagtutor == "BC2C" {
		if max > len(str) {
			fmt.Printf("\nError: the value max (<min>:<max>) in the lesson pair is too large, max is: %d.\n", len(str))
			os.Exit(99)
		}
		min--
		max--

		for i, c := range str {

			if i >= min && i <= max {
				tmpSet += string(c)
			}
		}
	} else {
		// its BC1S or BC2S
		num := 0
		arrBC1S := []string{"REA", "TIN", "PGS", "LCD", "HOF", "UWB", "REA", "TIN", "PGS", "LCD", "HOF", "UWB"}
		arrBC2S := []string{"KMY", "59,", "QXV", "73?", "+sb", "16.", "ZJ/", "28@", "40", "KMY", "59,", "QXV", "73?", "+sb", "16.", "ZJ/", "28$", "40"}

		if flagtutor == "BC1S" {
			num = len(arrBC1S)
		} else {
			num = len(arrBC2S)
		}

		if max > num {
			fmt.Printf("\nError: the value max (<min>:<max>) in the lesson pair is too large.\n")
			os.Exit(99)
		}
		min--
		max--

		// session approach
		for i := 0; i < num; i++ {
			if i >= min && i <= max {
				if flagtutor == "BC1S" {
					tmpSet += arrBC1S[i]
				} else {
					tmpSet += arrBC2S[i]
				}
			}
		}
	}

	return tmpSet
}
