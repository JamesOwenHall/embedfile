# embedfile

embedfile allows you to statically embed the contents of a file as a string in your Go code.

## Usage

    embedfile [options...] file1 [file2 [file3 ...]]

### Command Line Arguments

Syntax | Description
:------|:-----------
__-b__ | (*Not yet implemented*) Restricts the string to only use byte syntax (e.g. `"\x97"` instead of `"a"`)
__-f__ *filename* | Outputs the generated code to a file
__-p__ *package* | Sets the package name of the output.  Defaults to "main"
__-v__ | (*Not yet implemented*) Sets the output to be `var` instead of the default `const`