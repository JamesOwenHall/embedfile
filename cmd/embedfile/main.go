package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/JamesOwenHall/embedfile"
	"os"
	"path"
	"regexp"
	"strings"
)

// Command line options
var (
	bytemode    bool
	outFilename string
	packageName string
	variable    bool
)

var whiteSpace = regexp.MustCompile(`\w+`)

func main() {
	var err error

	flag.BoolVar(&bytemode, "b", false, "sets the type of the data to []byte")
	flag.StringVar(&outFilename, "o", "", "the name of the file to write the output to")
	flag.StringVar(&packageName, "p", "main", "the name of the package of the output")
	flag.BoolVar(&variable, "v", false, "sets the embedded data to var (default is const)")
	flag.Parse()

	// Check for input files
	if flag.NArg() == 0 {
		fmt.Println("error: no input file specified")
		return
	}

	// Get the output
	var output *os.File
	if len(outFilename) == 0 {
		output = os.Stdout
	} else {
		output, err = os.Create(outFilename)
		if err != nil {
			fmt.Println("error: unable to create file", outFilename)
			return
		}
	}
	defer output.Close()

	// Create the GoWriter
	outWriter := bufio.NewWriter(output)
	goWriter := embedfile.NewGoWriter(packageName, outWriter)
	goWriter.Open()

	// Write the files
	for i := 0; i < flag.NArg(); i++ {
		fileName := flag.Arg(i)
		varName := sanitizeVarName(fileName)

		file, err := os.Open(fileName)
		if err == os.ErrNotExist {
			fmt.Println("error: file", fileName, "does not exist")
			return
		} else if err != nil {
			fmt.Println("error: cannot open file", fileName)
			return
		}
		defer file.Close()

		err = goWriter.WriteFile(varName, file, variable, bytemode)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Flush the output buffer
	err = outWriter.Flush()
	if err != nil {
		fmt.Println("error:", err)
	}
}

// sanitizeVarName converts a file path into a suitable Go variable name.
func sanitizeVarName(filename string) string {
	// Trim the path
	filename = path.Base(filename)

	periodIndex := strings.LastIndex(filename, ".")
	if periodIndex == 0 {
		// File name begins with a period (e.g. ".htaccess"), just delete it
		filename = filename[1:]
	} else {
		// Trim the file extension
		filename = filename[:periodIndex]
	}

	// Limit the name to alphanumeric characters
	var result []rune
	for _, r := range filename {
		if alphabetic(r) || numeric(r) || r == '_' {
			result = append(result, r)
		}
	}

	// Names can't start with a number
	if numeric(result[0]) {
		result = prependRune('_', result)
	}

	return string(result)
}

// alphabetic returns true if r is a letter in the english alphabet.
func alphabetic(r rune) bool {
	if ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z') {
		return true
	}

	return false
}

// numeric returns true if r is a digit.
func numeric(r rune) bool {
	if '0' <= r && r <= '9' {
		return true
	}

	return false
}

// prependRune returns a slice starting with r and followed by
// the elements of s.
func prependRune(r rune, s []rune) []rune {
	result := make([]rune, 0, len(s)+1)

	result = append(result, r)
	for _, run := range s {
		result = append(result, run)
	}

	return result
}
