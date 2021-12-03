package main

import (
	"fmt"
	"strings"
	"flag"
	"os"
)

func doHelp() {
	flaghelp = strings.ToUpper(flaghelp)

	if flaghelp == "LCWO" {
		fmt.Println("\n\t\tLCWO Help Information")
		fmt.Printf("\nRelationships and compatibilities for LCWO Options (those prefixed by \"LCWO\".\n\nOption names starting with LCWO are only for users that will input the generated practice text\nto LCWO's \"Convert text to CW\" screen, or for input to LCWO (see Note).\n\nOther non-LCWO options can be used as well.\n\nLCWO_slow: (alternating slow/fast word groups)\n\t\trequires: LCWO_slow, LCWO_low, LCWO_fast, LCWO_step\n\t\toptional: LCWO_sf, LCWO_fs\n\t\tnon-compatible: LCWO_ramp, LCWO_effective_ramp, LCWO_num, LCWO_repeat\n")

		fmt.Printf("\nLCWO_ramp: (steady increase of character speed sections of words)\n\t\trequires: LCWO_ramp, LCWO_low, LCWO_step, LCWO_num\n\t\toptional: LCWO_sf, LCWO_repeat, LCWO_effective\n\t\tnon-compatible: LCWO_effective_ramp\n")

		fmt.Printf("\nLCWO_effective_ramp: (steady increase of effective speed sections, with fixed character speed)\n\t\trequires: LCWO_effective_ramp, LCWO_low, LCWO_step, LCWO_num, LCWO_effective\n\t\toptional: LCWO_sf\n\t\tnon-compatible: LCWO_ramp, LCWO_repeat\n")

		fmt.Printf("\nLCWO_repeat: (each word repeated at increasing character speeds)\n\t\trequires: LCWO_repeat, LCWO_low, LCWO_step\n\t\toptional: LCWO_effective, LCWO_ramp (LCWO_num if LCWO_ramp)\n\t\tnon-compatible: LCWO_effective_ramp\n")

		fmt.Printf("\nLCWO_random: LCWO_low, LCWO_num, LCWO_step (no other LCWO options): (\n\t\teach word has random character speed, between LCWO_low and LCWO_low+(LCWO_num * LCWO_num)\n\t\trequires: LCWO_low, LCWO_step, LCWO_num\n\t\toptional:\n\t\tnon-compatible: LCWO_effective_ramp, LCWO_ramp, LCWO_repeat, LCWO_slow,\n\t\tLCWO_fast, LCWO_sf, LCWO_fs\n")

		fmt.Printf("\nNote: LCWO_sf and LCWO_fs MAY be used outside of LCWO. If used with option mixedMode,\nthese will insert the specified string immediately before and/or after the codeGroup.\nIf LCWO_low and LCWO_step, are used with both the options <in> and <mixedMode>, then plain text is sent at the \nfaster speed of LCWO_low + LCWO_step, while the codeGroup will be sent at the slower speed of LCWO_low.")
	} else if flaghelp == "TUTORS" {
		fmt.Println("\n\t\tTUTORS Help Information")
		fmt.Printf(`
The option <tutor> and <lesson> are for user convience. By choosing the pair, you are prepopluating the
option <inlist> which reads words from the option <in> file, and <cglist> which is used to create code groups.

The generated practice text therefore can be given to ANY tutor. In some cases, a tutor will teach a ProSign,
but these two options only function with single characters.

The option <inlist> will be populated with both uppercase and lowercase for each alpha character.

The term "lesson" may not be used in every tutor, but its just the order that the character was taught.
Lockdown Morse videos are numbered, that is NOT the lesson number below. Simply find the characters covered
by the video (e.g. video 005 teaches B and P) that correlates to lesson 13 and 14 below.

Lesson is cumulative, that is if you enter lesson=5, all the characters from lesson 1 through 5 are used.

If the value for <tutor>, the names are NOT case sensitive, and for ease of typing you can use the abbreviated 
name given on the second line below.

Lesson  LCWO   JustLearnMC  G4FON  MorseElmer  MorseCodeNinja  HamMorse  LockdownMorse  MFJ418  PCWTutor  CWOPTs
Ch/Cnt           (JLMC)               (ME)          (MCN)        (HM)        (LDM)       (MFJ)   (PCWT)
------  ----   -----------  -----  ----------  --------------  --------  -------------  ------  --------  ------
1       K M                            T          T               E           E           W       Q S     T
2       U      K M          K M        M          A               M           O           B       E       E
3       R      R            R          R          E               R           A           M       M       A
4       E      S            S          S          N               S           I           H       T       N
5       S      U            U          U          O               U           U           A       A       O
6       N      A            A          A          I               A           Y           T       D       I
7       A      P            P          P          S               P           Z           J       J       S
8       P      T            T          T          1               T           Q           S       I       1
9       T      L            L          L          4               L           J           N       R       4
10      L      0            O          O          R               O           X           I       C       R

Lesson  LCWO   JustLearnMC  G4FON  MorseElmer  MorseCodeNinja  HamMorse  LockdownMorse  MFJ418  PCWTutor  CWOPTs
------  ----   -----------  -----  ----------  --------------  --------  -------------  ------  --------  ------
11      W      W            W          W          H              W            K (1)       O       5       H
12      I      I            I          I          D              I            V           D       N       D
13      .      .            .          .          L              .            B           E       L       L
14      J      N            N          N          2              N            P           L       G       U
15      Z      J            J          J          5              J            G           K       0       C
16      = <BT> E            E          E          C              G            W           Z       U       2
17      F      F            F          F          U              F            F           G       B       5
18      O      0            0          0          M              0            C           C       4       M
19      Y      Y            Y          Y          W              Y            L           U       1       W
20      ,      V            ,          ,          3              ,            D           Q       H       3

Lesson  LCWO   JustLearnMC  G4FON  MorseElmer  MorseCodeNinja  HamMorse  LockdownMorse  MFJ418  PCWTutor  CWOPTs
------  ----   -----------  -----  ----------  --------------  --------  -------------  ------  --------  ------
21      V      ,            v          V          6              V            M           R       O       6
22      G      G            G          G          ?              G            H           V       Z       ?
23      5      5            5          5          F              5            R           F       Y       F
24      /      /            /          /          Y              /            S           P       6       Y
25      Q      Q            Q          Q          P              Q            N           Y       9       3
26      9      9            9          9          G              9            T           X       K       6
27      2      Z            Z          Z          7              Z            0           5       W       P
28      H      H            H          H          9              H            1           .       2       G
29      3      3            3          3          /              3            5           7       7       7
30      8      8            8          8          B              8            6           /       F       9

Lesson  LCWO   JustLearnMC  G4FON  MorseElmer  MorseCodeNinja  HamMorse LockdownMorse  MFJ418  PCWTutor  CWOPTs
------  ----   -----------  -----  ----------  --------------  -------- -------------  ------  --------  ------
31      B      B            B          B          V              B            2          9        X       /
32      ?      ?            ?          ?          K              ?            7          ,        .       B
33      4      4            4          4          J              4            3          1        ?       v
34      7      2            2          2          8              2            8          6        3       J
35      C      7            7          7          0              7            4          8        8       K
36      1      C            C          C          = <BT>         C            9          ?        P       8
37      D      1            1          1          X              1            .          2        V       0
38      6      D            D          D          Q              6            ,          0        ,       <BT> =
39      0      6            6          6          Z              X            /          4        /       X
40      X      x            X          X          <BK>           = <BT>       ?          3        = <BT>  <SK>

Lesson  LCWO   JustLearnMC  G4FON  MorseElmer  MorseCodeNinja  HamMorse  LockdownMorse  MFJ418  PCWTutor  CWOPTs
------  ----   -----------  -----  ----------  --------------  --------  -------------  ------  --------  ------
41             @            = <BT>     .  
42             = <BT>       + <AR>(2)
43             + <AR> 

Notes: 
1- <KA> and <AR> introduced. Used only at start/end of sequences.

All the above prosigns, can be practiced with the prosign option, in delimiters, or as words
in an input file.

     **** Please scroll to the top of the screen. Use full screen to minimize line wraps. ****
			`) // end of table

	} else if flaghelp == "TOUR" {
		fmt.Println("\n\t\tWelcome To The MCPT Tour")
		fmt.Printf(`
This tour will give you an overview of the online help information. It will also demonstrate
a few features that you can experiment with to get started. (Many will NOT be covered.)

To take full advantage of the features, refer to the UserGuide.html.

Windows users, the delivered file will have the extension \".exe\", you can enter the command with or without the extension.

You have several approaches for a tour, from the least, to the most detail.

Approaches:

1 - Windows users, run "QuickStart.bat" (with or without the \".bat\" extension) in a terminal window.
    Non-window users, copy-and-paste, each command in "QuickStart.html" at the prompt.

2 - Read the file QuickStart.html, for a little more detail of most commands in number 1.

3 - Read the file UserGuide.html, for details on commands, files, options and examples.

    Some frequently used information is also available with the following help options:

    mcpt -help
    mcpt -help=tutors
    mcpt -help=files
    mcpt -help=LCWO
    mcpt -help=options

*** LCWO Users: The advanced, on-the-fly speed change features, controlled by the "LCWO"
*** options have NOT been demonstrated. They are fully explained in the UserGuide.html
*** in the examples, after #20 - you don't want to miss these!

New User? Where do you go from here? Copy-and-paste any of the commands listed in the QuickStart.html
and follow the instructions to try them. Modify the option values per the -help or user guide.
Provide the file created (remember to add -out=output.txt to commands without it) to your tutor and
listen to the code. Find more details in the UserGuide.html to help you.

Experiment. Enjoy. And hopefully see your comfort and speed with the Morse code increase.

     **** Please scroll to the top of the screen. Use full screen to minimize line wraps. ****
		`)
	} else if flaghelp == "FILES" {
		fmt.Println("\n\t\tFile Information: files that you might use")
		fmt.Printf(`

For certain features, or user convenience, you might specify a text file name in an option.
(the name is your choice, these are just to be descriptive, note a ".txt" lets you edit with Notepad)

-out=output.txt Add to any command to capture output in a file, instead of writing to the screen.
    This is the easiest way to get the text to the software that will play the code for you or
    for some tutors, like LCWO, create an MP3 file.

-prosign=prosign.txt A sample is in the download. Use it with word practice (i.e. -in=words.txt)
    or a command that is for code groups (i.e. -codeGroups). 
    A prosign (formally procedure signal) is a way communicate an instruction or short
    message, with a unique sound. In text, a prosign is indicated by two letters contained within
    a pair of angle brackets (i.e. <SK>, <BT>). Code tutors, know that this means the user wants
    the unique sound NOT the four characters, and NOT the two letters as if they were two letters of
    a word. The normal space that would be between the letters is eliminated.

-opt=option.txt An option.txt is provided. Any options that can be on the command line could 
    have been in an options file (with comments). Use of this can save time, and typos. Some users have 
    multiple option files, each for specific purposes.

    For example: instead of typing: 
    mcpt -codeGroups -num=200 -cgmin=3 -cgmax=6 -out=out.txt -lesson=12

    You could type: mcpt -opt=option.txt if you had created an file named option.txt, as below:

    **** everything below until the next starred line can be in the file

    #
    #  option.txt an option file short code groups
    #
    codeGroups
    num=200      # creates 200 groups
    cgmin=3      # the minimum group length is 3 characters
    cgmax=6      # the minimum group length is 6 characters
    out=out.txt  # write groups to this file
    lesson=12    # use all the characters for the first 12 lessons (using the LCWO as the default)

    ******

    There are several other features of this file, see the user guide.  A note worthy concept is that a command 
    line value overrides the files value. So after the user does the above  practice, as an example, he 
    could run:

    mcpt -opt=option.txt -lesson=15 
    Instead of editing the file, to change the lesson number.

    -in=word.txt This file is for those that want practice with "words". By "words", we mean a sequence of alpha 
       characters (i.e. name morse xyz IBM), and/or a prosign (i.e. <AR>). The prosign option (-prosign) is not 
       needed if the prosign is in the input file, just like words are.

       Words can be linked together with "~" (my~linked~words),  if the user wants those words to be treated 
       as a single word with other options (suffix, prefix, repeat; the tutor does not see the "~").

       No other punctuation, or numbers are preserved. The format of the input file is irrelevant,
       and by default so is the case of the letters.

     ***** below until the next starred line COULD be the format of an input file

     NASA is a United States space agency. They have had many programs:
       Mercury
       Apollo-11
       Apollo-12
       Gemini

     Many contractors: Boeing, IBM, AT&T, etc.
     1000 plus employees are on site 24 hours a day!

     ******
     Below are the unique words found above, with the previous conditions. 

     gemini are on is united had many plus site hours states space 
     they have etc nasa a agency boeing mercury ibm employees day 


     **** Please scroll to the top of the screen. Use full screen to minimize line wraps. ****
		`)
	} else if flaghelp == "OPTIONS" {
		fmt.Println("\n\t\tYour command line plus default values are shown below.\n")
		flag.VisitAll(func(f *flag.Flag) {
			if f.Name != "help" {
				fmt.Printf("%s=%v\n", f.Name, f.Value)
			}
		})
	} else {
		fmt.Printf("\nError: Invalid value for the option <help>, \n\tchoices are (case insensitive): TOUR, FILES, TUTORS, LCWO, or OPTIONS.\n")
	}
	os.Exit(1)
}
