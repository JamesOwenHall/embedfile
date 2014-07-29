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
	outFilename string
	packageName string
)

var whiteSpace = regexp.MustCompile(`\w+`)

func main() {
	var err error

	flag.StringVar(&outFilename, "f", "", "the name of the file to write the output to")
	flag.StringVar(&packageName, "p", "main", "the name of the package of the output")
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

		err = goWriter.WriteFile(varName, file)
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

	var result []rune
	for _, c := range filename {
		if ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') || c == '_' {
			result = append(result, c)
		}
	}

	return string(result)
}
