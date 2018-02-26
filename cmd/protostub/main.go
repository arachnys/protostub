package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}

// The second method is Set(value string) error
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)

	return nil
}

func convert(r io.Reader, w io.Writer) error {
	return nil
}

func main() {
	// Both shorthand and longhand.
	var files stringSlice
	flag.Var(&files, "file", "Pass the protobuf file to generate stubs for.")
	flag.Var(&files, "f", "Pass the protobuf file to generate stubs for.")

	var stdin = flag.Bool("stdin", false, "Reads from stdin instead of a file if true.")
	flag.BoolVar(stdin, "s", false, "Reads from stdin instead of a file if true.")

	var output = flag.String("output", "", "The generated mypi stub output.")
	flag.StringVar(output, "o", "", "The generated mypy stub output.")

	flag.Parse()

	// if reading from stdin, we can make some different assumptions
	// namely, if an output isn't specified, then the user is going to want
	// to write to stdout
	// also there's no point in checking for files from -f
	if *stdin {
		reader := os.Stdin
		writer := os.Stdout

		if len(*output) > 0 {
			f, err := os.Create(*output)

			if err != nil {
				panic(err)
			}

			writer = f
		}

		err := convert(reader, writer)

		if err != nil {
			panic(err)
		}

		// we're done.
		return
	}

	for _, file := range files {
		// just in case
		if len(file) < 0 {
			continue
		}

		fmt.Println("Working on: ", file)
		path := *output

		if len(*output) == 0 {
			path = strings.Replace(file, ".proto", "_pb2.pyi", -1)
		}

		f, err := os.Create(path)

		if err != nil {
			panic(err)
		}

		defer func() {
			f.Close()
		}()

		fmt.Println("Created: ", path)

		proto, err := os.Open(file)

		if err != nil {
			panic(err)
		}

		defer func() {
			proto.Close()
		}()

		err = convert(proto, f)

		if err != nil {
			panic(err)
		}
	}
}
