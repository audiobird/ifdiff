package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: ifdiff [IN FILE] [OUT FILE] 
        Write [IN FILE] to [OUT FILE] if [OUT FILE] doesn't exist or if the content of [OUT FILE] != [IN FILE]

        With no [IN FILE], or when [IN FILE] is -, read standard input`)
	}

	flag.Parse()
}

func main() {
	parseFlags()

	cerr := log.New(os.Stderr, "", 0)

	inpF := os.Stdin
	outF := ""

	switch len(os.Args) {
	case 2:
		outF = os.Args[1]
	case 3:
		outF = os.Args[2]
		if os.Args[1] == "-" {
			break
		}
		t, err := os.Open(os.Args[1])
		if err != nil {
			cerr.Fatal(err)
		}
		defer t.Close()
		inpF = t
	default:
		flag.Usage()
		os.Exit(1)
	}

	in, err := io.ReadAll(inpF)
	if err != nil {
		cerr.Fatal(err)
	}

	prevOut, err := os.ReadFile(outF)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			cerr.Fatal(err)
		}
	}

	if slices.Compare(in, prevOut) == 0 {
		os.Exit(0)
	}

	outFile, err := os.Create(outF)
	if err != nil {
		cerr.Fatal(err)
	}
	defer outFile.Close()

	_, err = outFile.Write(in)
	if err != nil {
		cerr.Fatal(err)
	}
}
