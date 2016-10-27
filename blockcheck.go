package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func readBlockContents(scanner *bufio.Scanner) (string, error) {
	blockStarted := false

	contents := ""

	for scanner.Scan() {
		line := scanner.Text()

		if blockStarted {
			if strings.HasPrefix(line, "```") {
				return contents + "\n", nil
			}

			if contents == "" {
				contents = line
			} else {
				contents = contents + "\n" + line
			}
		} else if strings.HasPrefix(line, "```") {
			blockStarted = true
		} else {
			return "", errors.New("Expected a code block (```) after BLOCKCHECK comment")
		}
	}

	return "", errors.New("Reached EOF before finding end of code block (```)")
}

func checkFile(markdownFile string) int {
	f, err := os.Open(markdownFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	re := regexp.MustCompile("<!-- MUSTMATCH: ([^ -]+) ?-->")

	checks := 0
	for scanner.Scan() {
		line := scanner.Text()

		if match := re.FindStringSubmatch(line); match != nil {
			blockContents, err := readBlockContents(scanner)
			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}

			relativePath := filepath.Dir(markdownFile)

			compareFile := filepath.Join(relativePath, string(os.PathSeparator), match[1])

			compareContents, err := ioutil.ReadFile(compareFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if string(compareContents) != blockContents {
				fmt.Printf("Block in %s does not match %s\n", markdownFile, compareFile)
				os.Exit(1)
			}
			checks++
		}
	}

	return checks
}

func main() {
	verbose := flag.Bool("v", false, "Print verbose output")

	flag.Parse()

	var fileNames []string
	if flag.NArg() == 0 {
		// Read file names from stdin.
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fileNames = strings.Split(strings.Trim(string(input), "\n"), "\n")
	} else {
		fileNames = flag.Args()
	}

	for _, f := range fileNames {
		if *verbose {
			fmt.Printf("Checking %s: ", f)
		}
		checks := checkFile(f)
		if *verbose {
			fmt.Printf("%d passed\n", checks)
		}
	}
}
