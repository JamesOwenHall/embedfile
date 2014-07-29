// package embedfile creates Go source code with the contents of files
// embedded as strings.
package embedfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var hexDigits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

type GoWriter struct {
	packageName string
	output      io.Writer
}

// NewGoWriter creates and returns an instance of GoWriter.
func NewGoWriter(packageName string, output io.Writer) GoWriter {
	return GoWriter{
		packageName: packageName,
		output:      output,
	}
}

// Open prepares the output for writing.
func (g *GoWriter) Open() error {
	_, err := fmt.Fprintln(g.output, "package ", g.packageName)
	return err
}

// Writes the contents of the file to the output.
func (g *GoWriter) WriteFile(varName string, file *os.File) error {
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanAllRunes)

	_, err := fmt.Fprint(g.output, "\nvar ", varName, " = \"")
	if err != nil {
		return err
	}

	for scanner.Scan() {
		token := scanner.Bytes()

		if len(token) == 1 && token[0] >= 0x80 || token[0] <= 31 {
			_, err = fmt.Fprint(g.output, `\x`+toHex(token[0]))
		} else if token[0] == '\\' {
			_, err = fmt.Fprint(g.output, `\\`)
		} else if token[0] == '"' {
			_, err = fmt.Fprint(g.output, `\"`)
		} else {
			_, err = g.output.Write(token)
		}

		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(g.output, "\"\n")
	return err
}

// toHex converts the byte into a string with its hexadecimal notation.
func toHex(b byte) string {
	return hexDigits[b>>4] + hexDigits[b&15]
}
