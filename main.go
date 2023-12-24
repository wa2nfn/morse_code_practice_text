//

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	program       = "mcpt"
	version       = "2.2.5" // 03/20/2023
	maxWordLen    = 60
	maxUserWords  = 5000
	maxLineLen    = 500
	maxSuffix     = 20
	maxPrefix     = 20
	maxDelimChars = 20
	maxRepeat     = 20
	maxMixedMode  = 20
	inListStr     = "A-Za-z"
)

var (
	kochChars      = ""
	seed           = time.Now().UTC().UnixNano()
	rng            = rand.New(rand.NewSource(seed))
	wordMap        = make(map[string]struct{})
	iwrWordMap     = make(map[string]struct{})
	wordArray      = make([]string, 0, 0)
	delimiterSlice []string
	effDelta       int
	proSign        []string
	runeMap        = make(map[rune]struct{})
	ccMap          = make(map[rune]rune)
	ps2runeMap     = make(map[string]rune)
	rune2psMap     = make(map[rune]string)
	/*
		// cannot use lc e or w, conflict with LCWO
	*/
	ps2charReplacer = strings.NewReplacer(
		"<AS>", "(",
		"<AR>", ")",
		"<BT>", "{",
		"<KA>", "}",
		"<SK>", "[",
		"<SN>", "]",
		"<HH>", "#",
		"\u00D8", "0",
		"+", ")",
		"=", "{",
		"<VE>", "$",
		"<DU>", "%",
		"<VA>", "!",
		"<SOS>", "&")
	char2psReplacer = strings.NewReplacer(
		"(", "<AS>",
		")", "<AR>",
		"{", "<BT>",
		"}", "<KA>",
		"[", "<SK>",
		"]", "<SN>",
		"#", "<HH>",
		"0", "\u00D8",
		"$", "<VE>",
		"!", "<VE>",
		"%", "<DU>",
		"&", "<SOS>")
)

var (
	flagLF            bool
	flagcc            bool
	flagTAB           bool
	flagdisplayFormat string
	flagmax           int
	flagcgmax         int
	flaglen           int
	flagmin           int
	flagcgmin         int
	flagrepeat        int
	flagrepeatStr     string
	flagnum           int
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
	flagLCWOiwrfile   string
	flagLCWOiwrwpm    int
	flagheader        string
	flagfooter        string
	flagprelist       string
	flagPrelistRune   []rune
	flagReview        bool
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
	flagordered       bool
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
	flagheadcopy      int
	flagsend          string
	flagsendcheck     string
	flagFL            bool
	flaglc            bool
	flagll            bool
	flagslashedzero          bool
	flagPhraseLen     int
	isLicw            bool
	scrambleWord      bool

	message string = `

  Error: 

  Either you forgot a required option (listed below), or your values were so 
  restrictive that there was nothing to show.

  If you are a new user, you might want to run: mcpt -help=tour or 
  run: mcpt -help, to review options.
  Examples are in the MCPT User Guide, as well.

  One of these is required:
  =========================
  -textFile   (for strings/words from a file)
  -inFile     (for words from a file)
  -codeGroups (code groups)
  -permute    (permutations of characters based on tutor/lesson)
  -help       (option listing)
  -review     (concentrated practice with random groups)
  -send       (create sending practice by specific groups of characters)
  -sendCheck  (to verify your accuracy of sending)
  -callSigns  (simple generated call signs based on tutor/lesson)
  -lessonList (use lesson chars or cgList for ever increasing code group)

` // end message
)

func init() {
	fmt.Fprintf(os.Stderr, "\n                         MCPT - Morse Code Practice Text\n                                     v %s\n                                    by WA2NFN\n\n", version)

	flag.StringVar(&flaginlen, "inLen", "1:5", "# characters in a word. inLen=min:max")
	flag.StringVar(&flagcglen, "cgLen", "5:5", "# characters in a code group. cgLen=min:max.")
	flag.StringVar(&flagrepeatStr, "repeat", "1", "Number of times to repeat word sequentially.")
	flag.IntVar(&flagnum, "num", 100, fmt.Sprintf("Number of words (or code groups) to output. Min 1, max %d.\n", maxUserWords))
	flag.IntVar(&flaglen, "len", 80, fmt.Sprintf("Length characters in an output line (max %d).", maxLineLen))
	flag.StringVar(&flagsuflen, "sufLen", "0:0", "The number of suffix characters to append. sufLen=min:max")
	flag.StringVar(&flagprelen, "preLen", "0:0", "The number of prefix characters to prefix. preLen=min:max")
	flag.BoolVar(&flagrandom, "random", false, "If preLen|sufLen is set, determines if it is used on a\nword-by-word basis. (default false)")
	flag.StringVar(&flagsuflist, "sufList", "0-9/,.?=", "Characters for a suffix.")
	flag.StringVar(&flagprelist, "preList", "0-9/,.?=", "Characters for a prefix.")
	flag.StringVar(&flaginlist, "inList", inListStr, "Characters to define an input word.")
	flag.StringVar(&flaginput, "inFile", "", "Input file, for words (including extension).")
	flag.StringVar(&flagtext, "textFile", "", "Input text file, for any strings in file (including extension).")
	flag.StringVar(&flagoutput, "outFile", "", "Output file for generated material.")
	flag.StringVar(&flagopt, "optFile", "", "Specify an option file to read or create.")
	flag.StringVar(&flagprosign, "prosignFile", "", "ProSign file name. One ProSigns per line. i.e. <BT>")
	flag.StringVar(&flagdelimit, "delimiter", "", "Output an inter-word delimiter string. A \"|\" separates delimiters e.g. <SK>|abc|123.\nA blank field e.g. aa| |bb, is valid to get a space. ")
	flag.BoolVar(&flagunique, "unique", false, "Each output word is sent only once (num option quantity may be reduced).\n (default false)")
	flag.StringVar(&flagtutor, "tutor", "LCWO", "Only with -lessons. Sets order and # of characters by tutor type.\nLCWO, JustLearnMorseCode, G4FON, MorseElmer, MorseCodeNinja, HamMorse, LockdownMorse, MFJ418, PCWTutor, CWOPTS, FARNSWORTH, BC1(S|C), BC2(S|C).\nUse -help=tutors for more info.")
	flag.StringVar(&flagDM, "delimiterNum", "0:0", "(If delimiter is used.) The number of delimiter\nstrings to add together. delimiterNum=min:max")
	flag.StringVar(&flaglesson, "lesson", "0:0", "Given the lesson number by <tutor>, populates options inlist and cglist with appropriate characters.")
	flag.BoolVar(&flagDR, "delimiterRandom", false, "Delimiter random, if delimiterNum > 0. (default false)")
	flag.IntVar(&flagMixedMode, "mixedMode", 0, fmt.Sprintf("mixedMode X, If X > 1 and  X < %d , a code group will print every X words.", maxMixedMode))
	flag.BoolVar(&flagreverse, "reverse", false, "Reverses the spelling of words from inlist file. (default false)")
	flag.BoolVar(&flagCG, "codeGroups", false, "Random code groups from cglist characters (default false).")
	flag.BoolVar(&flagordered, "ordered", false, "Ordered (Non-Randomized) output words read from input (num=0 limits output to qualified matching words) (default false).")
	flag.BoolVar(&flagMMR, "MMR", false, "Mixed-Mode-Random, randomize code group occurance in mixed mode. (default false)")
	flag.StringVar(&flagcglist, "cgList", "A-Z0-9/.,?=", "Set of characters to make code groups.")
	flag.StringVar(&flagheader, "header", "", "A string copied verbatim to head of output.")
	flag.StringVar(&flagfooter, "footer", "", "A string copied verbatim to foot of output.")
	flag.IntVar(&flagheadcopy, "headCopy", 0, "1: increment length by 1 for each word/group.\n2: each word spelled letter by letter, i.e. c co cod code")
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
	flag.StringVar(&flagLCWOiwrfile, "LCWO_iwrFile", "", "file of \"words\" to play at higher speed for IWR.")
	flag.IntVar(&flagLCWOiwrwpm, "LCWO_iwr_wpm", 0, "wpm speed to play IWR \"words\" from <LCWO_iwrFile>.")
	flag.StringVar(&flaghelp, "help", "", "[TOUR|FILES|LCWO|OPTIONS|TUTORS] more help of given topics.")
	flag.StringVar(&flagpermute, "permute", "", "Selected permutations of current \"lesson\" characters [p,t,b([pairs,triples,both)].")
	flag.BoolVar(&flagcallsigns, "callSigns", false, "Call signs with current lesson's characters.")
	flag.StringVar(&flagdisplayFormat, "displayFormat", "", "LF, TAB, TAB_LF, LF_TAB. Cosmetic output options to give more whitespace for easier screen reading.")
	flag.StringVar(&flagmust, "must", "", "A string of characters. Each output codeGroup/string/word, MUST get one character from this string.")
	flag.StringVar(&flagsend, "send", "", "A string of group numbers (0-7) to make sending practice groups.")
	flag.StringVar(&flagsendcheck, "sendCheck", "", "Two files: <mcptSend.txt,yourSent.txt>, one is output of -send, the other from you CW practice.")
	flag.BoolVar(&flagFL, "favorLast", false, "Favor last characters learned in code Groups.")
	flag.BoolVar(&flagReview, "review", false, "Concentrated random groups building on previous char.")
	flag.BoolVar(&flaglc, "lowercase", false, "Make output lowercase.")
	flag.BoolVar(&flagcc, "cc", false, "Emphasize confused character pairs in code groups.")
	flag.BoolVar(&flagll, "lessonList", false, "Use lesson chars or cglist for ever increasing group length.")
	flag.IntVar(&flagPhraseLen, "phraseLen", 0, "Number of words/groups per line as a phrase.")
	flag.BoolVar(&flagslashedzero, "slashedZero", false, "Default false, if set (-slashedZero), ZERO is displayed WITH a slash: \u00D8.")

	// rune map, validate option string like: cglist, prelist, delimiter
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
	runeMap['"'] = struct{}{}  // added at bottom of LCWO
	runeMap['\''] = struct{}{} // added at bottom of LCWO
	runeMap['('] = struct{}{}  // added at bottom of LCWO
	runeMap[')'] = struct{}{}  // added at bottom of LCWO
	runeMap['-'] = struct{}{}  // added at bottom of LCWO
	runeMap[':'] = struct{}{}  // added at bottom of LCWO
	runeMap[';'] = struct{}{}  // added at bottom of LCWO
	runeMap['*'] = struct{}{}  // DUMMY value for delimiter and LCWO users
}

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nIf above was NOT from -help, scroll to the top for more details.\nUsually: a misspelling, missing \"-\", missing space before \"-\", illegal space after \"-\", unmatched \".")

	}

	var fp *os.File
	kochChars = "KMURESNAPTLWI.JZ=FOY,VG5/Q92H38B?47C1D60X" // default for LCWO

	flag.Parse() // first parse to see if we had -opt

	if flagopt != "" || (flag.NFlag() == 0 && flag.NArg() == 1) {
		if flagopt == "" {
			flagopt = os.Args[1]
		}
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
					fmt.Printf("\n%s\nFile name <%s>.\n", err, flagopt)
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

				re := regexp.MustCompile("^opt|^help")
				for _, arg := range os.Args[1:] {
					arg = strings.TrimLeft(arg, "-")

					if re.MatchString(arg) {
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

	// save time
	hadLCWO := false
	markIt := func(f *flag.Flag) {

		if hadLCWO == false && strings.HasPrefix(f.Name, "LCWO_") {
			hadLCWO = true
		}
	}
	flag.Visit(markIt) 

	if flaginlist != inListStr {
		flaginlist = strings.ToUpper(flaginlist)
	}

	flagtutor = strings.ToUpper(flagtutor)
	flaglesson = strings.ToUpper(flaglesson)

	// handle split length options
	flagmin, flagmax = minmaxSplit("inLen", flaginlen)
	flagcgmin, flagcgmax = minmaxSplit("cgLen", flagcglen)
	flagpremin, flagpremax = minmaxSplit("preLen", flagprelen)
	flagsufmin, flagsufmax = minmaxSplit("sufLen", flagsuflen)
	flagDMmin, flagDMmax = minmaxSplit("DM", flagDM)

	if flagtutor == "BC1C" || flagtutor == "BC2C" || flagtutor == "BC1S" || flagtutor == "BC2S" {
		flagcglist = licw()
		isLicw = true
	} else {
		flaglessonstart, flaglessonend = minmaxSplit("lesson", flaglesson)
	}

	//
	// verify valid options
	//

	if flag.NArg() > 0 && flagsend == "" && flagopt == "" {
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

	if flagordered == false {
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
	} else if flagDMmin >= 1 || flagReview {
		// ok DM is in range
		// split into fields if any

		// if prosigns are in a delimiter field
		runeMap['<'] = struct{}{}
		runeMap['>'] = struct{}{}
		runeMap[' '] = struct{}{}
		runeMap['^'] = struct{}{}

		// first make sure any prosign is valid format
		m := regexp.MustCompile(`\s*\^[a-zA-Z]{2}\s*|\s*<[a-zA-Z]{2}>\s*|\s*\^SOS\s*|\s*<SOS>\s*|\s*^sos\s*|\s*<sos>\s*`)
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
		delete(runeMap, '^')
	}

	if flagheadcopy == 1 && flagtext != "" {
		fmt.Printf("\nError: option <headCopy=1> cannot be used with option <textFile> use <inFile>.\n")
		os.Exit(1)
	}

	if flagheadcopy == 2 && flaginput == "" && flagtext == "" {
		fmt.Printf("\nError: option <headCopy=2> requires <inFile> or <textFile>, cannot be used with option <CodeGroups>.\n")
		os.Exit(1)
	}

	if flagmin > flagmax {
		fmt.Printf("\nError: min must <= max <%d>, in -inLen.\n", flagmax)
		os.Exit(1)
	}

	if flagmax < flagmin {
		fmt.Printf("\nError: max must >= min <%d>, in -inLen.\n", flagmin)
		os.Exit(1)
	}

	if flagmax > maxWordLen {
		fmt.Printf("\nError: max in -inLen must <= <%d>, the system max.\n", maxWordLen)
		os.Exit(1)
	}

	if flagcgmax < flagcgmin {
		fmt.Println("\nError: max must >= min, in -cgLen.")
		os.Exit(1)
	}

	if flagMixedMode < 0 || flagMixedMode == 1 || flagMixedMode > maxMixedMode {
		fmt.Printf("\nError: mixedMode X Where X  min=2, max=%d, default 0=off.\n", maxMixedMode)
		os.Exit(1)
	}

	if flagordered {
		if flagnum < 0 || flagnum > maxUserWords {
			fmt.Printf("\nError: number of output words desired. minimum 0, maximum %d, default 100.\n", maxUserWords)
			fmt.Printf("\n       num = 0 means, match the number of qualified matching words\n")
			os.Exit(1)

		}
	} else {
		if flagnum < 1 || flagnum > maxUserWords {
			fmt.Printf("\nError: number of output words desired. minimum 0, maximum %d, default 100.\n", maxUserWords)
			os.Exit(1)
		}
	}

	if flaglen < 1 || flaglen > maxUserWords {
		fmt.Printf("\nError: len max output line length, default 80, maximum %d.\n", maxLineLen)
		os.Exit(1)
	}

	if flagsufmax > maxSuffix {
		fmt.Printf("\nError: sufLen, 0 = no suffix, max number of characters is %d.\n", maxSuffix)
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
			flagSuflistRune = ckValidListString(flagsuflist, "suflist")
		} else {
			fmt.Printf("\nError: if suffix > 0, the suflist must contain characters, its empty.\n")
			os.Exit(1)
		}
	}

	if flagpremax > maxPrefix {
		fmt.Printf("\nError: preLen, 0=no prefix, max number of characters is %d.\n", maxPrefix)
		os.Exit(1)
	}

	if flagpremin > 0 {
		if flagprelist != "" {
			// return expanded
			flagPrelistRune = ckValidListString(flagprelist, "prelist")
		} else {
			fmt.Printf("\nError: if preLen > 0, the preList must contain characters, its empty.\n")
			os.Exit(1)
		}
	}

	if flagrandom && (flagsufmin == 0 && flagpremin == 0 && flagll == false) {
		fmt.Printf("\nError: random requires: either preLen > 0 or sufLen > 0, or lessonList.\n")
		os.Exit(1)
	}

	flagrepeatStr = strings.ToUpper(flagrepeatStr)
	if strings.HasPrefix(flagrepeatStr, "R") {
		if flagtext != "" {
			// words file supports
			scrambleWord = true
		}
	}
	// always in case its done when we didn't want it
	flagrepeatStr = strings.TrimLeft(flagrepeatStr, "R")
	// make in int for everywhere else
	flagrepeat, _ = strconv.Atoi(flagrepeatStr)

	if flagrepeat < 1 || flagrepeat > maxRepeat {
		fmt.Printf("\nError: repeat (default 1) must be between 1 and %d.\n", maxRepeat)
		os.Exit(1)
	}

	if flagordered == true && (flagunique || flagCG) {
		fmt.Printf("\nError: Ordered is mutually exclusive with unique and codeGroups options.\n")
		os.Exit(1)
	}

	if flagoutput != "" {
		if flagoutput == flaginput {
			fmt.Printf("\nError: -outFile can't equal -inFile, or the input file would be over written.\n")
			os.Exit(1)
		}

		// check for existance first
		_, err := os.Stat(flagoutput)
		if err == nil {
			fmt.Printf("\nWarning: file: <%s> exists!\n\nEnter \"y\" to overwrite it: ", flagoutput)
			ans := ""
			fmt.Scanf("%s", &ans)
			if ans != "y" {
				fmt.Printf("\nNo change.\n")
				os.Exit(0)
			}
		}

		fp, err = os.Create(flagoutput)
		if err != nil {
			fmt.Println(err)
			os.Exit(9)
		}
		fmt.Printf("\nWriting file: %s\n", flagoutput)
	}

	// LCWO options
	// hard code some values since they are arbitrary
	if hadLCWO {

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

		if flagLCWOstep < -20 || flagLCWOstep > 20 {
			fmt.Printf("\nError: LCWO_step speed incremental step must be >= 0 and <= 30 wpm, 0(off).\n")
			os.Exit(0)
		}

		if flagLCWOstep*flagLCWOrepeat > 50 {
			fmt.Printf("\nError: Speed change in session is excessive at <%d> wpm, adjust LCWO_repeat or LCWO_step.\n", flagLCWOstep*flagLCWOrepeat)
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

		if flagLCWOiwrwpm != 0 {
			if flagLCWOiwrfile == "" {
				fmt.Printf("\nError: <LCWO_iwr_wpm> requires <LCWO_iwrFile>.\n")
				os.Exit(1)
			} else if flagLCWOrepeat > 0 || flagLCWOefframp || flagLCWOramp || flagLCWOrandom || flagLCWOslow > 0 || flagLCWOfast > 0 {
				fmt.Printf("\nError: IWR requires <LCWO_iwr_wpm>, <LCWO_iwrFile>, <LCWO_low>; optionally <LCWO_effective>.\n")
				fmt.Printf("\n       Any other LCWO options is an error.\n")
				os.Exit(1)
			}

			if flagLCWOiwrwpm <= flagLCWOlow {
				fmt.Println("\nError: <LCWO_iwr_wpm> must be greater than <LCWO_low>.")
				os.Exit(123)
			}

			if flagtext == "" && flaginput == "" {
				fmt.Println("\nError: The IWR feature requiers <LCWO_iwrFile> and <inFile> or <textFile>")
				os.Exit(123)
			}

			// read IWR file and store words in map
			readIwrFile(flagLCWOiwrfile)
		}
	}

	// expand now before we reuse (its UC)
	flagcglist = strRangeExpand(flagcglist, "cgList")
	flaginlist = strRangeExpand(flaginlist, "inList")

	// new special handling of LICW carousel
	if isLicw == false {
		if flaglessonend >= 1 {

			if flagtutor == "LCWO" {
				kochChars = "KMURESNAPTLWI.JZ=FOY,VG5/Q92H38B?47C1D60X"
			} else if flagtutor == "JUSTLEARNMORSECODE" || flagtutor == "JLMC" || flagtutor == "JUSTLEARNMC" {
				kochChars = "KMRSUAPTLOWI.NJEF0YV,G5/Q9ZH38B?427C1D6X@=+"
			} else if flagtutor == "G4FON" {
				kochChars = "KMRSUAPTLOWI.NJEF0Y,VG5/Q9ZH38B?427C1D6X"
			} else if flagtutor == "MORSEELMER" || flagtutor == "ME" {
				kochChars = "KMRSUAPTLOWI.NJEF0Y,VG5/Q9ZH38B?427C1D6X=+"
			} else if flagtutor == "MORSECODENINJA" || flagtutor == "MCN" {
				kochChars = "TAENOIS14RHDL25CUMW36?FYPG79/BVKJ80=XQZ!."
			} else if flagtutor == "CWOPTS" {
				kochChars = "TEANOIS14RHDL25UCMW36?FYPG79/BVKJ80=XQZ" //cwopts
			} else if flagtutor == "HAMMORSE" || flagtutor == "HM" {
				kochChars = "KMRSUAPTLOWI.NJEF0Y,VG5/Q9ZH38B?427C1D6X=+"
			} else if flagtutor == "LOCKDOWNMORSE" || flagtutor == "LDM" {
				flagtutor = "LOCKDOWNMORSE"
				kochChars = "EOAIUYZQJXKVBPGWFCLDMHRSNT0156273849.,/?"
			} else if flagtutor == "MFJ418" || flagtutor == "MFJ" {
				flagtutor = "MFJ418"
				kochChars = "WBMHATJSNIODELKZGCUQRVFPYX5.7/9,168?2043"
			} else if flagtutor == "PCWTUTOR" || flagtutor == "PCWT" {
				flagtutor = "PCWT"
				kochChars = "QSEMTADJIRC5NLG0UB41HOZY69KW27FX.?38PV,/="
			} else if flagtutor == "FARNSWORTH" || flagtutor == "FW" {
				flagtutor = "FW"
				kochChars = "TAEHCSNOL.BIFRW?DYMGUP,\"VKXQJZ(;:12345/-67890=d+"
			} else {
				fmt.Printf("\nError: Your tutor name is invalid. Names are NOT case sensitive, and without any spaces, see -helpi=tutors.\n")
				os.Exit(1)
			}

			if (flaglessonend+1 > len(kochChars)) && (flagtutor == "LCWO" || flagtutor == "G4FON" || flagtutor == "JLMC" || flagtutor == "PCWT") {
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

			if flagtutor == "LCWO" || flagtutor == "G4FON" || flagtutor == "JLMC" || flagtutor == "PCWT" {
				flaglessonend++
			}

			if flaglessonend <= len(kochChars) {
				flagcglist = kochChars[flaglessonstart:flaglessonend]
				if flaginput != "" {
					flaginlist = flagcglist
				}
			}

		}
	}

	// must follow other cglist manipulation
	// either case lets get cglist evaluated now
	if flagCG || flagMixedMode > 0 || flagll {

		// make sure we have chars to work with
		if len(flagcglist) < 1 {
			fmt.Printf("\nError: codeGroups, mixedMode, and lesson list\n        require cgList have at least 1 character.\n")
			os.Exit(1)
		} else {
			if flagcglist != "" {
				// return expanded
				flagCglistRune = ckValidListString(flagcglist, "cgList")
			}
		}
	}

	// ***** enhanced for char confusion
	// need to see which potential confusing chars are in the flagcglist
	if flagcc {
		// populate the map for each confusing-char pair

		if strings.ContainsRune(flagcglist, 'Q') && strings.ContainsRune(flagcglist, 'Y') {
			ccMap['Q'] = 'Y'
			ccMap['Y'] = 'Q'
		}
		if strings.ContainsRune(flagcglist, 'F') && strings.ContainsRune(flagcglist, 'L') {
			ccMap['F'] = 'L'
			ccMap['L'] = 'F'
		}
		if strings.ContainsRune(flagcglist, 'W') && strings.ContainsRune(flagcglist, 'G') {
			ccMap['W'] = 'G'
			ccMap['G'] = 'W'
		}
		if strings.ContainsRune(flagcglist, 'R') && strings.ContainsRune(flagcglist, 'K') {
			ccMap['R'] = 'K'
			ccMap['K'] = 'R'
		}
		if strings.ContainsRune(flagcglist, 'U') && strings.ContainsRune(flagcglist, 'D') {
			ccMap['U'] = 'D'
			ccMap['D'] = 'U'
		}
		if strings.ContainsRune(flagcglist, 'P') && strings.ContainsRune(flagcglist, 'X') {
			ccMap['P'] = 'X'
			ccMap['X'] = 'P'
		}
		if strings.ContainsRune(flagcglist, 'A') && strings.ContainsRune(flagcglist, 'N') {
			ccMap['A'] = 'N'
			ccMap['N'] = 'A'
		}
		if strings.ContainsRune(flagcglist, 'V') && strings.ContainsRune(flagcglist, 'B') {
			ccMap['B'] = 'V'
			ccMap['V'] = 'B'
		}
		if strings.ContainsRune(flagcglist, ',') && strings.ContainsRune(flagcglist, '?') {
			ccMap[','] = '?'
			ccMap['?'] = ','
		}
		if strings.ContainsRune(flagcglist, 'S') && strings.ContainsRune(flagcglist, 'O') {
			ccMap['S'] = 'O'
			ccMap['O'] = 'S'
		}
		if strings.ContainsRune(flagcglist, 'M') && strings.ContainsRune(flagcglist, 'I') {
			ccMap['M'] = 'I'
			ccMap['I'] = 'M'
		}
	}

	if flaginput == "" && flagtext == "" && flagCG == false && flagpermute == "" && flagcallsigns == false && (flagsend == "" && flagsendcheck == "") && flagReview == false && flagll == false {
		nm := filepath.Base(os.Args[0])
		if strings.HasSuffix(nm, ".exe") {
			nm = strings.ReplaceAll(nm, ".exe", "")
		}

		fmt.Printf("%s", message)

		os.Exit(99)
	}

	if flaglesson != "0" && flaglesson != "0:0" {
		if flagtutor == "" {
			fmt.Printf("\nError: <lesson> requires option <tutor> to have a valid value.\n")
			os.Exit(98)
		}
	}

	if (flaglesson == "0" || flaglesson == "0:0") && flagCG {
		if flagcglist == "" {
			fmt.Printf("\nError: <codeGroups> requires option <lesson> greater than zero OR cglist must be used.\n")
			os.Exit(98)
		}
	}

	if flagmust != "" {
		flagmust = string(ckValidListString(flagmust, "must"))
	}

	if flagheader != "" {
		if flaglc && strings.ContainsAny(flagheader, "ABCDEFGHIJKLMNOPRSTUVWXYZ") {
			fmt.Printf("\nWarning: <flagheader> contains uppercase letters, <lowercase> will change them.\n         You will need to edit them in the output if that will be a problem.\n\nEnter \"y\", for yes make all lowercase: ")
			ans := ""

			fmt.Scanf("%s", &ans)
			if ans != "y" {

				fmt.Println("\nNo output as requested.")
				os.Exit(0)
			}
			fmt.Println()
		}
	}

	if flagfooter != "" {
		if flaglc && strings.ContainsAny(flagfooter, "ABCDEFGHIJKLMNOPRSTUVWXYZ") {
			fmt.Printf("\nWarning: <flagfooter> contains uppercase letters, <owercase> will change them.\n         You will need to edit them in the output if that will be a problem.\n\nEnter \"y\", for yes make all lowercase: ")
			ans := ""

			fmt.Scanf("%s", &ans)
			if ans != "y" {

				fmt.Println("\nNo output as requested.")
				os.Exit(0)
			}
			fmt.Println()
		}
	}

	//
	// Major flow decision - WORD_MODE , CODE_GROUPS, PERMUTE, CALLSIGN ?
	//
	if flagFL == true && len(flagcglist) > 1 {
		// favorLast
		var last byte = flagcglist[len(flagcglist)-1]
		var nextToLast byte
		if len(flagcglist) >= 2 {
			nextToLast = flagcglist[len(flagcglist)-2]
			flagcglist = flagcglist + strings.Repeat(string(nextToLast), 3)
		}
		// increase last char usage
		flagcglist = flagcglist + strings.Repeat(string(last), 6)
	}

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

	// do -send or -sendCheck
	if flagsend != "" || flagsendcheck != "" {
		doSendOpts(fp)
		os.Exit(0) // program done
	}

	if flagtext != "" {
		readStringsFile(fp)
		doOutput(wordArray, fp)
		os.Exit(0) // program done
	}

	if flagReview {
		review(fp)
		os.Exit(0) // program done
	}

	if flagll {
		lessonlist(fp)
		os.Exit(0) // program done
	}

	// word mode default or codeGroups
	readFileMode(fp)
	doOutput(wordArray, fp)
}

// make sure the string can be expanded into visable ASCII since
// all Morse is limited to that
func ckValidListString(ck string, whoAmI string) []rune {

	// need to see if shell or os did a path substitution
	if strings.Contains(ck, ":\\") {
		fmt.Printf("\nWarning:\n\nIf you entered a single \"\\\" after an = in option <%s> this will likely\nconfuse the operating system. Change it to \\/", whoAmI)
		fmt.Printf(" or use an option file.\n")
		os.Exit(99)
	}

	if flagtutor == "FW" {
		// no for other cases or 'd' becomes <KA>
		ck = strings.Replace(ck, "d", "<KA>", -1) // needed for Farnsworth
	}
	str := strings.ToUpper(ck)
	str = ps2charReplacer.Replace(str)

	if strings.Contains(str, "<>") {
		fmt.Printf("\nError: Option <%s> contains an unsupported ProSign or the character \"<\" or\" >\".\n", whoAmI)
		os.Exit(99)
	}

	str = strRangeExpand(str, whoAmI)

	// check each rune to make sure its in the runeMap
	// also build a new string that is in the proper case

	for _, read := range str {
		if whoAmI != "delimiter" && read == '*' {
			fmt.Printf("\nError: Invalid character <%c>, in option <%s>.\nOnly used in delimiter option, as a special case delay for LCWO users.\n", read, whoAmI)
			os.Exit(98)
		}

		if _, ok := runeMap[read]; ok {
			continue
		} else {
			if whoAmI != "delimiter" && read == '^' {
				fmt.Printf("\nError: Invalid entry in string <%s> for option <%s>.\n", str, whoAmI)
				os.Exit(99)
			}
		}
	}

	if whoAmI == "delimiter" {
		str = strings.ReplaceAll(str, "*", "|S250 ")
		delimiterSlice = append(delimiterSlice, str)
	}

	return []rune(str)
}

/*
** walk the users string for prelist or suflist to see if we need to expand char ranges
 */
func strRangeExpand(inStr string, whoAmI string) string {

	outStr := ""
	last := ""
	gotDash := false

	if inStr == "" {
		return outStr
	}

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

// prints the bufStr adjusting the length per flaglen
// PRINTIT
func printStrBuf(strBuf string, fp *os.File) {

	index := 0
	spaceCount := 0

	// if doing IWR its LCWO so frame words with speed change
	if flagLCWOiwrwpm > 0 {

		// need start speed only once

		origSpd := ""
		iwrSpd := ""
		if flagLCWOeff > 0 {
			// deal with effective
			// this is just at the start
			strBuf = "|w" + strconv.Itoa(flagLCWOlow) + " |e" + strconv.Itoa(flagLCWOeff) + " " + strBuf
			// these two used many times
			origSpd = "|w" + strconv.Itoa(flagLCWOlow) + " |e" + strconv.Itoa(flagLCWOeff) + " "
			iwrSpd = "|e0 |w" + strconv.Itoa(flagLCWOiwrwpm) + " "
		} else {
			// no eff so make same as wpm
			strBuf = "|e0 |w" + strconv.Itoa(flagLCWOlow) + " " + strBuf
			origSpd = "|w" + strconv.Itoa(flagLCWOlow) + " "
			iwrSpd = "|w" + strconv.Itoa(flagLCWOiwrwpm) + " "
		}

		// since LCWO has a browser limit lets trim the buffer
		// probably 8000
		maxChars := 8000
		if len(strBuf) > maxChars {
			strBuf = strBuf[:maxChars]
		}

		for wd, _ := range iwrWordMap {
			pat := fmt.Sprintf("\\s%s\\s", regexp.QuoteMeta(wd))
			re := regexp.MustCompile(pat)

			strBuf = re.ReplaceAllString(strBuf, iwrSpd+wd+"%"+origSpd)
			strBuf = re.ReplaceAllString(strBuf, iwrSpd+wd+"%"+origSpd)
		}

		strBuf = strings.ReplaceAll(strBuf, "%", " ")
		// we'll let LCWO do final trimming
		if len(strBuf) > maxChars {
			strBuf = strBuf[:maxChars]
		}
	}

	if flaglessonend > 14 && (flagtutor == "LOCKDOWNMORSE" || flagtutor == "LDM") {
		strBuf = "<KA>\n" + strBuf + "\n<AR>"
		index = -5
	}

	if isLicw {
		strBuf = strings.ReplaceAll(strBuf, "@", " BK ")
	}

	strBuf = strings.ReplaceAll(strBuf, "  ", " ")

	// done processing now output it
	res := ""

	for _, r := range strBuf {
		// added for use with Precision PC Tutor
		if flagPhraseLen != 0 && r == ' ' {
			spaceCount++

			if spaceCount >= flagPhraseLen {
				res += string('\n')
				spaceCount = 0
				index = 0
				continue
			}
		}

		if index <= flaglen {
			res = res + string(r)
			if r == '\n' {
				index = 0
			} else {
				index++
			}
		} else {
			if r != ' ' && r != '\n' {
				res = res + string(r)
				index++
			} else {
				res += "\n"
				// extra space if no -out
				if flagLF {
					res += "\n"
				}
				index = 0
			}
		}
	}

	// for MorseNinja
	if flagtutor == "MORSECODENINJA" || flagtutor == "MCN" && flaglessonend >= 40 {
		res = strings.ReplaceAll(res, "!", "<BK>")
	}

	if flagTAB {
		res = strings.TrimRight(res, " ")
		res = strings.ReplaceAll(res, " ", "	")
	}

	// revert lc to ProSigns
	res = char2psReplacer.Replace(res)

	if flagfooter != "" {
		res += "\n" + flagfooter
	}

	if flaglc {
		res = strings.ToLower(res)
	}

	if flagslashedzero {
		fmt.Println("Use -slashedZero, to convert zeros to display without a slashes.")
	} else {
		res = strings.ReplaceAll(res, "\u00D8", "0")
	}

	if flagoutput == "" {
		fmt.Printf("%s\n", res)
		os.Exit(0)
	} else {
		if res == "" {
			fmt.Printf("%s", message)
			fp.Close()
			os.Exit(0)
		}

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

// called for each field of a delimiter option
func processDelimiter(inStr string) {
	// eliminate special case of simple range
	m := regexp.MustCompile("^([0-9]-[0-9])|([A-Z]-[A-Z])$")

	inStr = strings.ToUpper(inStr)
	if m.MatchString(inStr) {
		expandIt(string(inStr[0]), string(inStr[2]), "delimiter_simple")
		return
	}

	ckValidListString(inStr, "delimiter")
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

// read iwr file and store words in file
func readIwrFile(file string) {
	var wordLimit int = 1000

	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("\n%s reading IWR file name <%s>.\n", err, file)
		os.Exit(1)
	}

	// parse it for words

	line := string(content)
	line = strings.ToUpper(line)

	for _, wd := range strings.Fields(line) {
		cnt := 0
		wd = strings.TrimSpace(wd)
		if strings.HasPrefix(wd, "#") {
			continue
		}

		iwrWordMap[wd] = struct{}{}

		cnt++
		if cnt >= wordLimit {
			break
		}
	}

	lenMap := len(iwrWordMap)
	if lenMap == 0 {
		fmt.Printf("\nError: IWR file <%s> has no entries.\n", file)
		os.Exit(1)
	}
}
