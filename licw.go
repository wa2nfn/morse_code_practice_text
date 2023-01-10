package main

import(
	"fmt"
	"strings"
	"os"
	"strconv"
)

func licw() string {
	b1seed := "REATINPGSLCDHOFUWB"
	b2seed := "KMY59,QXV73?+sb16.ZJ/28@40"
	str := ""
	tmpSet := ""
	var max int
	var min int

	switch flaglesson {
	case "B1":
		if flagtutor == "B2C" || flagtutor == "B2S" {
			fmt.Printf("\nError: lesson value incomplete with the tutor=<%s>, require lesson ID(s)\n", flagtutor)
			os.Exit(99)
		}
		return b1seed
	case "B2":
		if flagtutor == "B1C" || flagtutor == "B1S" {
			fmt.Printf("\nError: lesson value <B2> is invalid with the tutor=<%s>\n", flagtutor)
			os.Exit(99)
		}
		return b2seed
	case "B1:B2", "B2:B1":
		return b1seed + b2seed
	}

	b1seedDbl := "REATINPGSLCDHOFUWBREATINPGSLCDHOFUWB"
	b2seedDbl := "KMY59,QXV73?+sb16.ZJ/28@40KMY59,QXV73?+sb16.ZJ/28@40"

	// analyse min and max
	if strings.HasPrefix(flaglesson, "B1:") {

		if flagtutor == "B2C" || flagtutor == "B2S" {
			flaglesson = strings.TrimPrefix(flaglesson,"B1:")
			tmpSet = b1seed
			str = b2seedDbl
		} else {
			fmt.Printf("\nError: lesson value for this tutor can ONLY have\n       these formats: B1, <ID>, <ID>:<ID>\n")
			os.Exit(99)
		}

	} else if strings.HasPrefix(flaglesson, "B2:") {
			fmt.Printf("\nError: lesson value never contains B2\n")
			fmt.Printf("\nlesson value for this tutor can ONLY have these formats: B2, B1:<ID>, B1:<ID>:<ID>, <ID>, <ID>:<ID>\n")
			os.Exit(99)
	}

	// char mode, session done later
	if flagtutor == "B1C" {
		str = b1seedDbl
	} else if flagtutor == "B2C" {
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
	if flagtutor == "B1C" || flagtutor == "B2C" {
		if max > len(str) {
			fmt.Printf("\nError: the value max (<min>:<max>) in the lesson pair is too large, max is: %d.\n",len(str))
			os.Exit(99)
		}
		min--
		max--

		for i,c := range str { 

			if i >= min && i <= max { 
				tmpSet += string(c)
			}
		}
	} else {
		// its B1S or B2S
		num := 0
		arrB1S := []string{"REA","TIN","PGS","LCD","HOF","UWB","REA","TIN","PGS","LCD","HOF","UWB"}
		arrB2S := []string{"KMY","59,","QXV","73?","+sb","16.","ZJ/","28@","40","KMY","59,","QXV","73?","+sb","16.","ZJ/","28$","40"}

		if flagtutor == "B1S" {
			num = len(arrB1S)
		} else {
			num = len(arrB2S)
		}

		if max > num {
			fmt.Printf("\nError: the value max (<min>:<max>) in the lesson pair is too large.\n")
			os.Exit(99)
		}
		min--
		max--

		 // session approach
		 for i := 0 ; i < num; i++ { 
			if i >= min && i <= max { 
				if flagtutor == "B1S" {
					tmpSet += arrB1S[i]
				} else {
					tmpSet += arrB2S[i]
				}
			}
		}
	}

	return tmpSet
}
