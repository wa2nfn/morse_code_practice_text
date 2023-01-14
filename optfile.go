package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func doOptFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	ignore := regexp.MustCompile("^\\s*#|^\\s*$|^\\s*//")
	doneEnd := regexp.MustCompile("^\\s*#\\s*(END|DONE)$")
	blkStart := regexp.MustCompile("^\\s*/\\*")
	blkEnd := regexp.MustCompile("^\\s*\\*")
	lineNum := 1
	inBlk := false

	for scanner.Scan() {
		str := scanner.Text()

		// start block comment
		if inBlk == false && blkStart.MatchString(str) {
			lineNum++
			inBlk = true
			continue
		}

		// end block comment
		if inBlk == true {
			lineNum++

			if blkEnd.MatchString(str) {
				inBlk = false
				continue
			} else {
				continue
			}
		}

		if doneEnd.MatchString(str) {
			return
		}

		if ignore.MatchString(str) {
			lineNum++
			continue
		}

		// cuts EOL comments off
		dex := strings.Index(str, "#")
		if dex != -1 {
			str = str[:dex]
		}

		// we assume its an option
		// trim down to get the string at the end
//		str = strings.TrimSpace(str)
		str = strings.TrimLeft(str, " -")

		/*
		if str[len(str)-1] == '=' {
			fmt.Printf("\nError: Invalid format for option <%v> on line <%d> of file <%s>. Appears to be missing a value after \"=\".\n", str, lineNum, flagopt)
			os.Exit(7)

		}
		*/

		// = sep?
		arr := strings.SplitN(str, "=", 2)
		if len(arr) == 2 {
			if flag.Lookup(arr[0]) == nil {
				fmt.Printf("\nError: Invalid option <%s> on line <%d> of file <%s>.\n", arr[0], lineNum, flagopt)
				os.Exit(7)
			}

			arr[1] = strings.TrimLeft(arr[1], "'\"")
			arr[1] = strings.TrimRight(arr[1], "'\"")

			if arr[0] == "opt" {
				fmt.Printf("\nWarning: option <opt> can't be set in the options file <%s> on line <%d>. Ignoring it and continuing.\n\n", flagopt, lineNum)
				continue
			}

			if arr[0] == "help" {
				fmt.Printf("\nWarning: option <help> can't be set in the options file <%s> on line <%d>. Ignoring it and continuing.\n\n", flagopt, lineNum)
				continue
			}
			flag.Set(arr[0], arr[1])
			continue
		}

		// space sep?
		arr = strings.SplitN(str, " ", 2)
		if len(arr) == 2 {
			if flag.Lookup(arr[0]) == nil {
				fmt.Printf("\nError: Invalid option <%s> on line <%d> of file <%s>.\n", arr[0], lineNum, flagopt)
				os.Exit(7)
			}

			arr[1] = strings.TrimLeft(arr[1], "'\"")
			arr[1] = strings.TrimRight(arr[1], "'\"")

			if arr[0] == "opt" {
				fmt.Printf("\nWarning: option \"opt\" can't be reset in the options file <%s> on line <%d>. Ignoring it and continuing.\n", flagopt, lineNum)
				continue
			}

			flag.Set(arr[0], arr[1])
		} else if len(arr) == 1 {
			flag.Set(arr[0], "true")
			continue
		} else {
			if flag.Lookup(arr[0]) == nil {
				fmt.Printf("\nError: Invalid option <%s> on line <%d> of file <%s>.\n", arr[0], lineNum, flagopt)
				os.Exit(7)
			}
		}

		flag.Set(arr[0], "")

	}
}
