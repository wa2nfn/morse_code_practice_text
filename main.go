//
// Copyright 2019, 2020, 2021 Bill Lanahan - WA2NFN
//

package main

import (
	"flag"
	"fmt"
	_ "log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	program       = "mcpt"
	version       = "2.0" // 4/6/2021
	maxWordLen    = 40
	maxUserWords  = 5000
	maxLineLen    = 500
	maxSuffix     = 20
	maxPrefix     = 20
	maxDelimChars = 20
	maxRepeat     = 20
	maxSkips      = 5000
	maxMixedMode  = 20
	inListStr     = "A-Za-z"
)

var (
	kochChars      = ""
	seed           = time.Now().UTC().UnixNano()
	rng            = rand.New(rand.NewSource(seed))
	wordMap        = make(map[string]struct{})
	wordArray      = make([]string, 0, 0)
	delimiterSlice []string
	effDelta       int
	proSign        []string
	runeMap        = make(map[rune]struct{})
	ps2runeMap     = make(map[string]rune)
	rune2psMap     = make(map[rune]string)
)

var (
	flagLF            bool
	flagTAB           bool
	flagdisplayFormat string
	flagmax           int
	flagcgmax         int
	flaglen           int
	flagmin           int
	flagcgmin         int
	flagrepeat        int
	flagnum           int
	flagskip          int
	flagsufmin        int
	flagsufmax        int
	flagpremin        int
	flagpremax        int
	flagDM            string
	flagDMmin         int
	flagDMmax         int
	flagDR            bool
	flaglesson        string
	flagMixedMode     int
	flagLCWOsf        string
	flagLCWOfs        string
	flagLCWOstep      int
	flagLCWOnum       int
	flagLCWOslow      int
	flagLCWOlow       int
	flagLCWOfast      int
	flagLCWOrepeat    int
	flagLCWOeff       int
	flagLCWOramp      bool
	flagLCWOrandom    bool
	flagLCWOefframp   bool
	flagheader        string
	flagprelist       string
	flagPrelistRune   []rune
	flagsuflist       string
	flagSuflistRune   []rune
	flaginlist        string
	flagcglist        string
	flagCglistRune    []rune
	flaginput         string
	flagtext          string
	flagoutput        string
	flagopt           string
	flagprosign       string
	flagdelimit       string
	flagtutor         string
	flagrandom        bool
	flagunique        bool
	flagNR            bool
	flagMMR           bool
	flagCG            bool
	flagreverse       bool
	flaghelp          string
	flagpermute       string
	flaginlen         string
	flagprelen        string
	flagsuflen        string
	flagcglen         string
	flaglessonend     int
	flaglessonstart   int
	flagcallsigns     bool
	flagmust          string
	flaghead          bool
	flagsend          string
	flagsendcheck     string
)

var message string = `

  Error: 

  Either you forgot a required option (listed below), or your values were so 
  restrictive that there was nothing to show.

  If you are a new user, you might want to run: mcpt -help=tour or 
  run: mcpt -help, to review options.
  Examples are in the MCPT User Guide, as well.

  One of these is required:
  =========================
  -text       (for strings/words from a file)
  -in         (for only words from a file)
  -codeGroups (code groups)
  -permute    (permutations of characters based on tutor/lesson)
  -help       (option lisiting)
  -send       (create sending practice by specific groups of characters)
  -sendCheck  (to verify your accuracy of sending)
  -callsigns  (simple generated call signs based on tutor/lesson)

  	` // end message

func init() {
  flag.StringVar(&flaginlen, "inlen", "1:5", "# characters in a word. inlen=min:max")
	flag.StringVar(&flagcglen, "cglen", "5:5", "# characters in a code group. cglen=min:max.")
	flag.IntVar(&flagrepeat, "repeat", 1, "Number of times to repeat word sequentially.")
	flag.IntVar(&flagnum, "num", 100, fmt.Sprintf("Number of words (or code groups) to output. Min 1, max %d.\n", maxUserWords))
	flag.IntVar(&flaglen, "len", 80, fmt.Sprintf("Length characters in an output line (max %d).", maxLineLen))
	flag.IntVar(&flagskip, "skip", 0, fmt.Sprintf("Number of the first unique words in the input to skip. Max %d", maxSkips))
	flag.StringVar(&flagsuflen, "suflen", "0:0", "The number of suffix characters to append to words. suflen=min:max")
	flag.StringVar(&flagprelen, "prelen", "0:0", "The number of prefix characters to affix to words. prelen=min:max")
	flag.BoolVar(&flagrandom, "random", false, "If prefix/suffix is set, determines if used on a\nword-by-word basis. (default false)")
	flag.StringVar(&flagsuflist, "suflist", "0-9/,.?=", "Characters for a word suffix.")
	flag.StringVar(&flagprelist, "prelist", "0-9/,.?=", "Characters for a word prefix.")
	flag.StringVar(&flaginlist, "inlist", inListStr, "Characters to define an input word.")
	flag.StringVar(&flaginput, "in", "", "Input text file name, for words (including extension).")
	flag.StringVar(&flagtext, "text", "", "Input text file name, for any strings in file (including extension).")
	flag.StringVar(&flagoutput, "out", "", "Output file name.")
	flag.StringVar(&flagopt, "opt", "", "Specify an option file name to read or create.")
	flag.StringVar(&flagprosign, "prosign", "", "ProSign file name. One ProSigns per line. i.e. <BT>")
	flag.StringVar(&flagdelimit, "delimiter", "", "Output an inter-word delimiter string. A \"|\" separates delimiters e.g. <SK>|abc|123.\nA blank field e.g. aa| |bb, is valid to get a space. ")
	flag.BoolVar(&flagunique, "unique", false, "Each output word is sent only once (num option quantity may be reduced).\n (default false)")
	flag.StringVar(&flagtutor, "tutor", "LCWO", "Only with -lessons. Sets order and # of characters by tutor type.\nLCWO, JustLearnMorseCode, G4FON, MorseElmer, MorseCodeNinja, HamMorse, LockdownMorse, MFJ418.\nUse -help=tutors for more info.")
	flag.StringVar(&flagDM, "DM", "0:0", "Delimiter multiple, (if delimiter is used.) Between 1 and DM delimiter\nstrings are concatenated. DM=min:max")
	flag.StringVar(&flaglesson, "lesson", "0:0", "Given the lesson number by <tutor>, populates options inlist and cglist with appropriate characters.")
	flag.BoolVar(&flagDR, "DR", false, "Delimiter random, (if DM > 0) DR=true makes a delimiter randomly print on. (default false)")
	flag.IntVar(&flagMixedMode, "mixedMode", 0, fmt.Sprintf("mixedMode X, If X > 1 and  X < %d , a code group will print every X words.", maxMixedMode))
	flag.BoolVar(&flagreverse, "reverse", false, "Reverses the spelling of words from inlist file. (default false)")
	flag.BoolVar(&flagCG, "codeGroups", false, "Random code groups from cglist characters (default false).")
	flag.BoolVar(&flagNR, "NR", false, "Non-Randomized output words read from input (default false).")
	flag.BoolVar(&flagMMR, "MMR", false, "Mixed-Mode-Random, randomize code group occurance in mixed mode. (default false)")
	flag.StringVar(&flagcglist, "cglist", "A-Z0-9/.,?=", "Set of characters to make code groups.")
	flag.StringVar(&flagheader, "header", "", "string copied verbatim to head of output")
	flag.BoolVar(&flaghead, "headCopy", false, "Used with codeGroups to increment length by 1 for each group, until cglen is max")
	flag.IntVar(&flagLCWOlow, "LCWO_low", 15, "low character speed setting (wpm).")
	flag.IntVar(&flagLCWOstep, "LCWO_step", 0, "speed change increment (wpm).")
	flag.IntVar(&flagLCWOslow, "LCWO_slow", 0, "number of words to send at slower speed.")
	flag.IntVar(&flagLCWOfast, "LCWO_fast", 0, "number of words to send at faster speed.")
	flag.IntVar(&flagLCWOnum, "LCWO_num", 0, "number of speed change steps.")
	flag.BoolVar(&flagLCWOramp, "LCWO_ramp", false, "ramp speed up in steps (default false).")
	flag.BoolVar(&flagLCWOrandom, "LCWO_random", false, "random speed per word (default false).")
	flag.IntVar(&flagLCWOrepeat, "LCWO_repeat", 0, "times to repeat each word with increasing speed.")
	flag.IntVar(&flagLCWOeff, "LCWO_effective", 0, "effective (aka Farnsworth) speed must be < LCWO_low (wpm).")
	flag.BoolVar(&flagLCWOefframp, "LCWO_effective_ramp", false, "ramp effective speed (char speed constant) must be < LCWO_low. (default false)")
	flag.StringVar(&flagLCWOsf, "LCWO_sf", "", "to alert transition from LCWO_low to LCWO_low+LCWO_step for plain text in mixedMode\nor LCWO_slow text to LCWO_fast text,")
	flag.StringVar(&flagLCWOfs, "LCWO_fs", "", "to alert transition from LCWO_low+LCWO_step speed for plain text to LCWO_low for codeGroup mixedMode\nor LCWO_fast text to LCWO_slow text.")
	flag.StringVar(&flaghelp, "help", "", "[TOUR|FILES|LCWO|OPTIONS|TUTORS] more help of given topics.")
	flag.StringVar(&flagpermute, "permute", "", "Selected permutations of current \"lesson\" characters [p,t,b([pairs,triples,both)].")
	flag.BoolVar(&flagcallsigns, "callSigns", false, "Call signs with current lesson's characters.")
	flag.StringVar(&flagdisplayFormat, "displayFormat", "", "LF, TAB, TAB_LF, LF_TAB. Cosmetic output options to give more whitespace for easier screen reading.")
	flag.StringVar(&flagmust, "must", "", "A string of characters. Each output codeGroup/string/word, MUST get one character from this string.")
	flag.StringVar(&flagsend, "send", "", "A string of group numbers (1-5) to make sending practice groups.")
	flag.StringVar(&flagsendcheck, "sendCheck", "", "Two files: <mcptSend.txt,yourSent.txt>, one is output of -send, the other from you CW practice.")

	// fill the rune map which is used to validate option string like: cglist, prelist, delimiter
	runeMap['a'] = struct{}{}
	runeMap['b'] = struct{}{}
	runeMap['c'] = struct{}{}
	runeMap['d'] = struct{}{}
	runeMap['e'] = struct{}{}
	runeMap['f'] = struct{}{}
	runeMap['g'] = struct{}{}
	runeMap['h'] = struct{}{}
	runeMap['i'] = struct{}{}
	runeMap['j'] = struct{}{}
	runeMap['k'] = struct{}{}
	runeMap['l'] = struct{}{}
	runeMap['m'] = struct{}{}
	runeMap['n'] = struct{}{}
	runeMap['o'] = struct{}{}
	runeMap['p'] = struct{}{}
	runeMap['q'] = struct{}{}
	runeMap['r'] = struct{}{}
	runeMap['s'] = struct{}{}
	runeMap['t'] = struct{}{}
	runeMap['u'] = struct{}{}
	runeMap['v'] = struct{}{}
	runeMap['w'] = struct{}{}
	runeMap['x'] = struct{}{}
	runeMap['y'] = struct{}{}
	runeMap['z'] = struct{}{}
	runeMap['A'] = struct{}{}
	runeMap['B'] = struct{}{}
	runeMap['C'] = struct{}{}
	runeMap['D'] = struct{}{}
	runeMap['E'] = struct{}{}
	runeMap['F'] = struct{}{}
	runeMap['G'] = struct{}{}
	runeMap['H'] = struct{}{}
	runeMap['I'] = struct{}{}
	runeMap['J'] = struct{}{}
	runeMap['K'] = struct{}{}
	runeMap['L'] = struct{}{}
	runeMap['M'] = struct{}{}
	runeMap['N'] = struct{}{}
	runeMap['O'] = struct{}{}
	runeMap['P'] = struct{}{}
	runeMap['Q'] = struct{}{}
	runeMap['R'] = struct{}{}
	runeMap['S'] = struct{}{}
	runeMap['T'] = struct{}{}
	runeMap['U'] = struct{}{}
	runeMap['V'] = struct{}{}
	runeMap['W'] = struct{}{}
	runeMap['X'] = struct{}{}
	runeMap['Y'] = struct{}{}
	runeMap['Z'] = struct{}{}
	runeMap['0'] = struct{}{}
	runeMap['1'] = struct{}{}
	runeMap['2'] = struct{}{}
	runeMap['3'] = struct{}{}
	runeMap['4'] = struct{}{}
	runeMap['5'] = struct{}{}
	runeMap['6'] = struct{}{}
	runeMap['7'] = struct{}{}
	runeMap['8'] = struct{}{}
	runeMap['9'] = struct{}{}
	runeMap[','] = struct{}{}
	runeMap['.'] = struct{}{}
	runeMap['/'] = struct{}{}
	runeMap['?'] = struct{}{}
	runeMap['='] = struct{}{}
	runeMap['+'] = struct{}{}
	runeMap['@'] = struct{}{}
	runeMap['!'] = struct{}{}  // added at bottom of LCWO
	runeMap['"'] = struct{}{}  // added at bottom of LCWO
	runeMap['\''] = struct{}{} // added at bottom of LCWO
	runeMap['('] = struct{}{}  // added at bottom of LCWO
	runeMap[')'] = struct{}{}  // added at bottom of LCWO
	runeMap['-'] = struct{}{}  // added at bottom of LCWO
	runeMap[':'] = struct{}{}  // added at bottom of LCWO
	runeMap[';'] = struct{}{}  // added at bottom of LCWO
	runeMap['*'] = struct{}{}  // DUMMY value for delimiter and LCWO users
}

func init() {
	fmt.Fprintf(os.Stderr, "\n                         MCPT - Morse Code Practice Text\n                                  version %s\n                                    by WA2NFN\n\n", version)
}

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nIf the output above was NOT from entry of -help, scroll to the first line which gives a hint.\nusually: a misspelling, missing \"-\", missing space before \"-\", illegal space after \"-\", unmatched \".")

	}

	kochChars = "KMURESNAPTLWI.JZ=FOY,VG5/Q92H38B?47C1D60X" // default for LCWO
	var fp *os.File
	localSkipFlag := false
	localSkipCount := 0

	flag.Parse() // first parse to see if we had -opt

	if flagopt != "" {
		_, err := os.Stat(flagopt)

		if err != nil {
			fmt.Printf("\nError: Can't find options file=<%s>.\n\nDo you want to create one with the current options? Enter \"y\" or \"n\": ", flagopt)

			ans := ""
			fmt.Scanf("%s", &ans)
			if ans != "y" {
				fmt.Printf("\nBye!\n")
				os.Exit(0)
			} else {
				// create the file
				fp, err := os.Create(flagopt)
				if err != nil {
					fmt.Printf("\n%s File name <%s>.\n", err, flagopt)
				}

				defer func() {
					closeErr := fp.Close()
					if closeErr != nil {
						if err == nil {
							err = closeErr
						} else {
							fmt.Println("Error occured while closing the file :", closeErr)
							os.Exit(2)
						}
					}
				}()

				skip := regexp.MustCompile(`^help|^opt`)
				for _, arg := range os.Args[1:] {
					arg = strings.TrimLeft(arg, "-")
					if skip.MatchString(arg) {
						continue
					}

					arg += "\n"

					_, err := fp.WriteString(arg)
					if err != nil {
						fmt.Println(err)
						fp.Close()
						os.Exit(0)
					}
				}
				fmt.Printf("\nOption file=<%s> has been created.\n", flagopt)
			}

			fp.Close()
			os.Exit(1)
		}

		optFile, err := os.Open(flagopt)
		if err != nil {
			fmt.Printf("\n%s File name <%s>.\n", err, flagopt)
			os.Exit(1)
		}

		// do file parse
		doOptFile(optFile)
		optFile.Close()

		flag.Parse() // second parse since options read
	}

	inListChanged := false
	flaginlist = strings.ToUpper(flaginlist)
	if flaginlist != inListStr {
		inListChanged = true
	}

	// handle split length options
	flagmin, flagmax = minmaxSplit("inlen", flaginlen)
	flagcgmin, flagcgmax = minmaxSplit("cglen", flagcglen)
	flagpremin, flagpremax = minmaxSplit("prelen", flagprelen)
	flagsufmin, flagsufmax = minmaxSplit("suflen", flagsuflen)
	flaglessonstart, flaglessonend = minmaxSplit("lesson", flaglesson)
	flagDMmin, flagDMmax = minmaxSplit("DM", flagDM)

	//
	// verify valid options
	//

	if flag.NArg() > 0 {
		fmt.Printf("\nError processing the command line.\n\nYou may have:\n   forgotten a \"-\" before an option\n   or followed a \"-\" with a space\n   or added extra input\n   or put spaces around the \"=\"\n")
		os.Exit(1)
	}

	if flagpermute != "" {
		if flagpermute != "p" && flagpermute != "t" && flagpermute != "b" {
			fmt.Printf("\nError: values for option <permute> are: p (pairs), t (triples), b (both).\n")
			os.Exit(0)
		}
	}

	if flaghelp != "" {
		doHelp()
	}

	if flagNR == false {
		wordArray = nil // save space we're using the map not array
	}

	if flagdisplayFormat != "" {
		// see if we have arguments after the last flag
		flagdisplayFormat = strings.ToUpper(flagdisplayFormat)

		switch flagdisplayFormat {
		case "LF":
			flagLF = true
		case "LF_TAB", "TAB_LF":
			flagLF = true
			flagTAB = true
		case "TAB":
			flagTAB = true
		default:
			fmt.Printf("\nError: option <displayFormat> < %s > is invalid. Use: TAB, LF, LF_TAB, TAB_LF\n", flagdisplayFormat)
			os.Exit(77)
		}
	}

	//
	// out of range checks
	if flagDMmax > maxDelimChars {
		fmt.Printf("\nError: DM, delimiter multiple min >= 0, max <= %d.\n", maxDelimChars)
		os.Exit(1)
	} else if flagDMmin >= 1 {
		// ok DM is in range
		// split into fields if any

		// if prosigns are in a delimiter field
		runeMap['<'] = struct{}{}
		runeMap['>'] = struct{}{}
		runeMap[' '] = struct{}{}
		runeMap['*'] = struct{}{}
		runeMap['^'] = struct{}{}

		// first make sure any prosign is valid format
		m := regexp.MustCompile(`\s*\^[a-zA-Z]{2}\s*|\s*<[a-zA-Z]{2}>\s*`)
		tStr := flagdelimit // in case we need original later

		for _, field := range strings.Split(tStr, "|") {

			if field == "" {
				continue
			}

			if m.MatchString(field) {
				field = strings.TrimSpace(field)

				if !ckProsign(field) {
					fmt.Printf("\nError: option <delimiter> has more than a prosign in a field: %s \n", field)
					os.Exit(77)
				}

				processDelimiter(field)
				continue
			} else {
				if len(field) > 0 {
					if strings.Contains(field, "<") || strings.Contains(field, ">") || strings.Contains(field, "^") {
						fmt.Printf("\nError: option <delimiter> contains invalid prosign format < %s >.\n", field)
						os.Exit(77)
					}
				}
			}

			processDelimiter(field)

		}

		delete(runeMap, ' ')
		delete(runeMap, '<')
		delete(runeMap, '>')
		delete(runeMap, '*')
		delete(runeMap, '^')
	}

	if flagmin > flagmax {
		fmt.Printf("\nError: min must <= max <%d>, in -inlen.\n", flagmax)
		os.Exit(1)
	}

	if flagmax < flagmin {
		fmt.Printf("\nError: max must >= min <%d>, in -inlen.\n", flagmin)
		os.Exit(1)
	}

	if flagmax > maxWordLen {
		fmt.Printf("\nError: max in -inlen must <= <%d>, the system max.\n", maxWordLen)
		os.Exit(1)
	}

	if flagcgmax < flagcgmin {
		fmt.Println("\nError: cgmax must >= cgmin, in -cglen.")
		os.Exit(1)
	}

	if flagMixedMode < 0 || flagMixedMode == 1 || flagMixedMode > maxMixedMode {
		fmt.Printf("\nError: mixedMode X Where X  min=2, max=%d, default 0=off.\n", maxMixedMode)
		os.Exit(1)
	}

	if flagskip < 0 || flagskip > maxSkips {
		fmt.Printf("\nError: skip x  minimum 0, maximum %d, default 0.\n", maxSkips)
		os.Exit(1)
	}

	if flagskip >= 1 {
		// we will be skipping some words
		localSkipFlag = true
		localSkipCount = flagskip
	}

	if flagnum < 1 || flagnum > maxUserWords {
		fmt.Printf("\nError: num, number of output words desired. minimum 1, maximum %d, default 100.\n", maxUserWords)
		os.Exit(1)
	}

	if flaglen < 1 || flaglen > maxUserWords {
		fmt.Printf("\nError: len max output line length, default 80, maximum %d.\n", maxLineLen)
		os.Exit(1)
	}

	if flagsufmax > maxSuffix {
		fmt.Printf("\nError: suflen, 0 = no suffix, max number of characters is %d.\n", maxSuffix)
		os.Exit(1)
	}

	if flagCG == true {
		if flagMixedMode > 0 {
			fmt.Printf("\nError: mixedMode is mutually exclusive with codeGroups option.\n")
			os.Exit(1)
		}

		flaginlist = "" // incompatible
	}

	if flagsufmin > 0 {
		if flagsuflist != "" {
			flagSuflistRune = ckValidInString(flagsuflist, "suflist")
		} else {
			fmt.Printf("\nError: if suffix > 0, the suflist must contain characters, its empty.\n")
			os.Exit(1)
		}
	}

	if flagpremax > maxPrefix {
		fmt.Printf("\nError: prelen, 0=no prefix, max number of characters is %d.\n", maxPrefix)
		os.Exit(1)
	}

	if flagpremin > 0 {
		if flagprelist != "" {
			// return expanded
			flagPrelistRune = ckValidInString(flagprelist, "prelist")
		} else {
			fmt.Printf("\nError: if prelen > 0, the prelist must contain characters, its empty.\n")
			os.Exit(1)
		}
	}

	if flagrandom && (flagsufmin == 0 && flagpremin == 0) {
		fmt.Printf("\nError: random requires either prelen > 0 or suflen > 0.\n")
		os.Exit(1)
	}

	if flagrepeat < 1 || flagrepeat > maxRepeat {
		fmt.Printf("\nError: repeat (default 1) must be between 1 and %d.\n", maxRepeat)
		os.Exit(1)
	}

	if flagNR == true && (flagunique || flagCG) {
		fmt.Printf("\nError: NR is mutually exclusive with unique and codeGroups options.\n")
		os.Exit(1)
	}

	if flagoutput != "" {
		if flagoutput == flaginput {
			fmt.Printf("\nError: -out can't equal -in, or the input file would be over written.\n")
			os.Exit(1)
		}

		// check for existance first
		_, err := os.Stat(flagoutput)
		if err == nil {
			fmt.Printf("\nWarning: out file: <%s> exists!\n\nEnter \"y\" to overwrite it: ", flagoutput)
			ans := ""
			fmt.Scanf("%s", &ans)
			if ans != "y" {
				fmt.Printf("\nNo output as requested.\n")
				os.Exit(0)
			}
		}

		fp, err = os.Create(flagoutput)
		if err != nil {
			fmt.Println(err)
			os.Exit(9)
		}
		fmt.Printf("\nWriting to file: %s\n", flagoutput)
	}

	// LCWO options
	// hard code some values since they are arbitrary

	if flagLCWOnum < 0 || flagLCWOnum > 30 {
		fmt.Printf("\nError: LCWO_num number of speed values must be >= 0 and <= 30.\n")
		os.Exit(0)
	}

	if flagLCWOnum > 0 {
		if flagLCWOslow > 0 || flagLCWOfast > 0 {
			fmt.Printf("\nError: LCWO_num must = 0(off) if LCWO_slow > 0.\n")
			os.Exit(0)
		}
	}

	if flagLCWOlow < 5 {
		fmt.Printf("\nError: LCWO_low lowest speed must be at least 5 wpm.\n")
		os.Exit(0)
	}

	if flagLCWOeff > 0 {
		if flagLCWOeff >= flagLCWOlow {
			fmt.Printf("\nError: LCWO_effective speed must be < LCWO_low wpm.\n")
			os.Exit(0)
		}

		// set delta since eblow and eff are set
		effDelta = flagLCWOlow - flagLCWOeff
	}

	if flagLCWOstep < 0 || flagLCWOstep > 20 {
		fmt.Printf("\nError: LCWO_step speed incremental step must be >= 0 and <= 30 wpm, 0(off).\n")
		os.Exit(0)
	}

	if flagLCWOstep*flagLCWOnum > 50 {
		fmt.Printf("\nError: Speed change in session is excessive at <%d> wpm, adjust LCWO_numor LCWO_step.\n", flagLCWOstep*flagLCWOnum)
		os.Exit(0)
	}

	if flagLCWOslow < 0 || flagLCWOfast < 0 {
		fmt.Printf("\nError: LCWO_fast and LCWO_slow must be >= 0, 0(off).\n")
		os.Exit(0)
	}

	if flagLCWOrepeat > 1 && flagLCWOfast > 0 {
		fmt.Printf("\nError: LCWO_repeat is mutually exclusive with LCWO_slow.\n")
		os.Exit(0)
	}

	if flagLCWOfast > 0 && flagLCWOslow == 0 {
		fmt.Printf("\nError: LCWO_fast > 0 requires LCWO_slow to be specified.\n")
		os.Exit(0)
	}

	// we want LCWO, lots of exclusions to try
	if flagLCWOslow > 0 {
		if flagLCWOfast == 0 {
			flagLCWOfast = flagLCWOslow
		}

		if flagLCWOstep < 1 {
			fmt.Printf("\nError: LCWO_step must be >=1 with LCWO_slow/LCWO_fast.\n")
			os.Exit(0)
		}

		if flagLCWOrepeat >= 1 {
			fmt.Printf("\nError: LCWO_repeat is mutually exclusive with LCWO_slow and LCWO_fast options.\n")
			os.Exit(1)
		}
	}

	if flagLCWOrepeat < 0 || flagLCWOrepeat > 30 {
		fmt.Printf("\nError: LCWO_repeat must be >=2 and <= 30 for word speed repeat, 0(off).\n")
		os.Exit(0)
	}

	if flagLCWOramp {
		if flagLCWOfast > 0 {
			fmt.Printf("\nError: LCWO_ramp is mutually exclusive with LCWO_slow and LCWO_fast options.\n")
			os.Exit(1)
		}

		if flagLCWOnum == 0 {
			fmt.Printf("\nError: LCWO_ramp requires LCWO_num > 0.\n")
			os.Exit(1)
		}

		if flagLCWOstep == 0 {
			fmt.Printf("\nError: LCWO_ramp requires LCWO_step > 0.\n")
			os.Exit(1)
		}
	}

	if flagLCWOefframp {
		if flagLCWOrepeat > 0 {
			fmt.Printf("\nError: LCWO_effective_ramp is mutually exclusive with LCWO_repeat.\n")
			os.Exit(1)
		}

		if flagLCWOnum == 0 {
			fmt.Printf("\nError: LCWO_effective_ramp requires LCWO_num > 0.\n")
			os.Exit(1)
		}

		if flagLCWOstep < 1 {
			fmt.Printf("\nError: LCWO_effective_ramp requires LCWO_step >= 1 and its less than LCWO_low.\n")
			os.Exit(1)
		}

		if flagLCWOramp {
			fmt.Printf("\nError: LCWO_effective_ramp is mutually exclusive with LCWO_ramp.\n")
			os.Exit(1)
		}

		if flagLCWOlow < 1 {
			fmt.Printf("\nError: LCWO_effective_ramp requires LCWO_low >= 5.\n")
			os.Exit(1)
		}
	}

	flagtutor = strings.ToUpper(flagtutor)

	// expand now before we reuse (its UC)
	flagcglist = strRangeExpand(flagcglist, "cglist")

	// may need to concatenate
	origInlist := ""
	if flagtext != "" || flaginput != "" {
		origInlist = strings.ToUpper(flaginlist)
		origInlist = strRangeExpand(origInlist, "inlist")
	}

	if flaglessonend >= 1 {

		if flagtutor == "LCWO" {
			kochChars = "KMURESNAPTLWI.JZ=FOY,VG5/Q92H38B?47C1D60X"
		} else if flagtutor == "JUSTLEARNMORSECODE" || flagtutor == "JLMC" {
			kochChars = "KMRSUAPTLOWI.NJEF0YV,G5/Q9ZH38B?427C1D6X@=+"
		} else if flagtutor == "G4FON" {
			kochChars = "KMRSUAPTLOWI.NJEF0Y,VG5/Q9ZH38B?427C1D6X"
		} else if flagtutor == "MORSEELMER" || flagtutor == "ME" {
			kochChars = "KMRSUAPTLOWI.NJEF0Y,VG5/Q9ZH38B?427C1D6X=+"
		} else if flagtutor == "MORSECODENINJA" || flagtutor == "MCN" {
			kochChars = "TAENOIS14RHDL25CUMW36?FYPG79/BVKJ80=XQZ!."
		} else if flagtutor == "HAMMORSE" || flagtutor == "HM" {
			kochChars = "KMRSUAPTLOWI.NJEF0Y,VG5/Q9ZH38B?427C1D6X=+"
		} else if flagtutor == "LOCKDOWNMORSE" || flagtutor == "LDM" {
			flagtutor = "LOCKDOWNMORSE"
			kochChars = "EOAIUYZQJXKVBPGWFCLDMHRSNT0156273849.,/?"
		} else if flagtutor == "MFJ418" || flagtutor == "MFJ" {
			flagtutor = "MFJ418"
			kochChars = "WBMHATJSNIODELKZGCUQRVFPYX5.7/9,168?2043"
		} else {
			fmt.Printf("\nError: Your tutor name is invalid. Names are NOT case sensitive, and without any spaces, see the help.\n")
			os.Exit(1)
		}

		if (flaglessonend+1 > len(kochChars)) && (flagtutor == "LCWO" || flagtutor == "G4FON" || flagtutor == "JLMC") {
			fmt.Printf("\nError: Lesson value <%d> exceeds the max <%d>, for tutor <%s>.\n", flaglessonend, 40, flagtutor)
			os.Exit(1)
		}

		if flaglessonend > len(kochChars) {
			fmt.Printf("\nError: Lesson value <%d> exceeds the max <%d>, for tutor <%s>.\n", flaglessonend, len(kochChars), flagtutor)
			os.Exit(1)
		}

		if flaglessonstart <= 0 {
			fmt.Printf("\nError: Lesson values start at 1. (see -help=tutors)\n")
			os.Exit(1)
		}

		if flaglessonstart == flaglessonend {
			flaglessonstart = 1
		}
		flaglessonstart-- // because strings start at 0

		if flagtutor == "LCWO" || flagtutor == "G4FON" || flagtutor == "JLMC" {
			flaglessonend++
		}

		if flaglessonend <= len(kochChars) {
			flagcglist = kochChars[flaglessonstart:flaglessonend]
		}

		tmp_c := ""
		for _, c := range flagcglist {
			if c >= 'A' && c <= 'Z' {
				tmp_c += string(c)
			}
		}

		flaginlist = tmp_c
	}

	// must follow other cglist manipulation
	// either case lets get cglist evaluated now
	if flagCG || flagMixedMode > 0 {

		// make sure we have chars to work with
		if len(flagcglist) < 1 {
			fmt.Printf("\nError: you requested codeGroups or mixedMode, so cglist must have at least 1 characters.\n")
			os.Exit(1)
		} else {
			if flagcglist != "" {
				// return expanded
				flagCglistRune = ckValidInString(flagcglist, "cglist")
			}
		}
	}

	// no longer needed save space
	runeMap = nil

	if flaginput == "" && flagtext == "" && flagCG == false && flagpermute == "" && flagcallsigns == false && (flagsend == "" && flagsendcheck == "") {
		nm := filepath.Base(os.Args[0])
		if strings.HasSuffix(nm, ".exe") {
			nm = strings.ReplaceAll(nm, ".exe", "")
		}

		fmt.Printf("%s", message)

		os.Exit(99)
	}

	if flagtutor == "" && (flaglesson != "0" || flaglesson != "0:0") {
		fmt.Printf("\nError: <lesson> requires option <tutor> to have a valid value.\n")
		os.Exit(98)
	}

	if (flaglesson == "0" || flaglesson == "0:0") && flagCG {
		if flagcglist == "" {
			fmt.Printf("\nError: <codeGroups> requires option <lesson> greater than zero OR cglist must be used.\n")
			os.Exit(98)
		}
	}

	// if text or words we can have inlist added to lesson
	if flagtext != "" || flaginput != "" {
		if inListChanged {
			flaginlist += origInlist
		}
	}

	if flagmust != "" {
		flagmust = strings.ToUpper(flagmust)
	}

	//
	// major flow decision - WORD_MODE , CODE_GROUPS, PERMUTE, CALLSIGN ?
	//
	if flagCG || (flagpermute != "" && flagprosign != "") {
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

		if flagpermute == "" {
			makeGroups(fp)
			os.Exit(0) // program done
		}
	}

	if flagpermute != "" {
		permute(flagpermute, fp)
		os.Exit(0) // program done
	}

	if flagcallsigns == true {
		if len(flagcglist) == 0 {
			flagcglist = kochChars
		}

		doCallSigns(flagcglist, fp)
		os.Exit(0) // program done
	}

	if flagtext != "" {
		readStringsFile(localSkipFlag, localSkipCount, fp)
		doOutput(wordArray, fp)
		os.Exit(0) // program done
	}

	// does -send or -sendCheck
	if flagsend != "" || flagsendcheck != "" {
		doSendOpts(fp)
		os.Exit(0) // program done
	}

	// word mode default
	readFileMode(localSkipFlag, localSkipCount, fp)
	doOutput(wordArray, fp)
}

//
// make sure the string can be expanded into visable ASCII since all Morse is limited to that
func ckValidInString(ck string, whoAmI string) []rune {

	// need to see if shell or os did a path substitution
	if strings.Contains(ck, ":/") {
		fmt.Printf("\nWarning:\n\nIf you entered \":/ or / at the start of the string\" in option <%s> as seen here <%s>, or you used the \"/\" character,\nyour operating system is incorrectly changing that into a PATH variable. Try to move the position of the \"/\"", whoAmI, ck)
		fmt.Printf("\nor put the option in an option file to prevent this.\n")
		os.Exit(99)
	}

	str := strRangeExpand(ck, whoAmI)

	// check each rune to make sure its in the runeMap
	// also build a new string that is in the proper case
	newRune := []rune{}

	for _, runeRead := range []rune(str) {
		if whoAmI != "delimiter" && runeRead == '*' {
			fmt.Printf("\nError: Invalid character <%s>, in option <%s>.\nOnly used in delimiter option, as a special case delay for LCWO users.\n", string(runeRead), whoAmI)
			os.Exit(98)
		}

		if _, ok := runeMap[runeRead]; ok {
			newRune = append(newRune, unicode.ToUpper(runeRead))
		} else {
			if whoAmI == "delimiter" && runeRead == '^' {

			} else {
				fmt.Printf("\nError: Invalid entry in string <%s> for option <%s>.\n", str, whoAmI)
				os.Exit(99)
			}
		}
	}

	if whoAmI == "delimiter" {
		s := string(newRune)
		s = strings.ReplaceAll(s, "*", "|S250 ")
		delimiterSlice = append(delimiterSlice, s)
		newRune = nil
	}

	return newRune
}

/*
** walk the users string for prelist or suflist to see if we need to expand char ranges
 */
func strRangeExpand(inStr string, whoAmI string) string {

	outStr := ""
	last := ""
	gotDash := false

	// take care of special case of dash in beginning
	if inStr[0] == '-' {
		fmt.Printf("\nError: \"-\" character(s) at start of list option: %s=%s\n", whoAmI, inStr)
		os.Exit(7)
	}

	for _, char := range strings.Split(inStr, "") {

		if char == "-" {
			if gotDash == true {
				fmt.Printf("\nError: sequencial \"-\" characters in a list option: %s=%s\n", whoAmI, inStr)
				os.Exit(7)
			}

			gotDash = true
			continue
		}

		if gotDash == false {
			last = char
			outStr += last
		} else {
			outStr += expandIt(last, char, whoAmI)
			gotDash = false
			last = ""
		}
	}

	//  detect error
	if gotDash && len(inStr) > 1 {
		fmt.Printf("\nError: trailing \"-\" characters in a list option: %s=%s\n", whoAmI, inStr)
		os.Exit(7)
	}

	return outStr
}

// expandIt expands a char range into the individual chars
func expandIt(lower string, upper string, whoAmI string) string {
	outStr := ""

	low := rune(lower[0])
	up := rune(upper[0])

	if up < low {
		fmt.Printf("\nError: range in an option list is not in ASCII/UTF-8 order: i.e. C-A (invalid) vs. A-C (correct) \n")
		fmt.Printf("       Delimiters support ONLY a single range in a field. i.e. ^A-D^ or ^0-3^\n")
		os.Exit(7)
	}

	if whoAmI != "delimiter_simple" {
		low++ // we did low already for all other list options
	}

	for i := low; i <= up; i++ {
		if whoAmI == "delimiter_simple" {
			delimiterSlice = append(delimiterSlice, string(i))
		} else {
			outStr += string(i)
		}
	}
	return outStr
}

//
// prints the bufStr adjusting the length per flaglen
//
func printStrBuf(strBuf string, fp *os.File) {

	re := regexp.MustCompile(`!`) // for MorseNinja
	index := 0
	if flagtutor == "LOCKDOWNMORSE" || flagtutor == "LDM" && flaglessonend > 14 {
		strBuf = "<KA>\n" + strBuf + "\n<AR>"
		index = -5
	}

	// done processing now output it
	res := ""
	for _, r := range strBuf {

		if index <= flaglen {
			res = res + string(r)
			index++
			continue
		}

		if index >= flaglen {
			if r != ' ' && r != '\n' {
				res = res + string(r)
				index++
				continue
			} else {
				res = res + "\n"
				// extra space if no -out
				if flagLF {
					res += "\n"
				}
				index = 0
			}
		} else {
			res = ""
		}
	}

	// for MorseNinja
	if flagtutor == "MORSECODENINJA" || flagtutor == "MCN" && flaglessonend >= 40 {
		res = re.ReplaceAllString(res, "<BK>")
	}

	if flagTAB {
		res = strings.TrimRight(res," ")
		res = strings.ReplaceAll(res, " ", "	")
	}

	if flagoutput == "" {
		fmt.Printf("%s", res)
		os.Exit(0)
	} else {
		if res == "" {
			fmt.Printf("%s", message)
			fp.Close()
			os.Exit(0)
		}

		res += string('\u0008') // marked as from mcpt
		_, err := fp.WriteString(res)
		if err != nil {
			fmt.Println(err)
			fp.Close()
			os.Exit(0)
		}
		os.Exit(0)
	}
} // end

// simple random true or false
func flipFlop() bool {
	return rng.Int()%2 == 0
}

//
// called for each field of a delimiter option
func processDelimiter(inStr string) {
	// eliminate special case of simple range
	m := regexp.MustCompile("^([0-9]-[0-9])|([a-z]-[a-z])|([A-Z]-[A-Z])$")
	if m.MatchString(inStr) {
		expandIt(string(inStr[0]), string(inStr[2]), "delimiter_simple")
		return
	}

	ckValidInString(inStr, "delimiter")
}

// to split min and max from combined input
func minmaxSplit(name string, value string) (int, int) {
	var op []string
	var max int

	value = strings.TrimSpace(value)
	op = strings.Split(value, ":")

	if len(op) > 2 || len(op) == 0 {
		fmt.Printf("\nError: invalid format for option <%s>\n", name)
		os.Exit(99)
	}

	min, error := strconv.Atoi(op[0])

	if error != nil {
		fmt.Printf("\nError: invalid format for option <%s>\n", name)
		os.Exit(99)
	}

	if len(op) == 2 && op[1] != "" {
		max, error = strconv.Atoi(op[1])
		if error != nil {
			fmt.Printf("\nError: invalid format for option <%s>\n", name)
			os.Exit(99)
		}
	} else {
		max = min
	}

	if min < 0 || max < 0 {
		fmt.Printf("\nError: invalid value for option <%s>, min and max <%s> must be >= 0.\n", name, value)
		os.Exit(99)
	}

	if max < min {
		fmt.Printf("\nError: invalid value for option <%s>, max must be >= min <%s>.\n", name, value)
		os.Exit(99)
	}

	return min, max
}

// does slice of string contain a value
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
