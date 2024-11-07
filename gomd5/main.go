package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	flagFile = flag.String("file", "", "file used for input")
	flagURL  = flag.String("url", "", "url used for input")
)

var input io.Reader

func main() {
	flag.Parse()
	switch {
	case *flagFile != "":
		f, err := os.Open(*flagFile)
		if err != nil {
			fmt.Println("error opening file:", *flagFile, err)
			os.Exit(123)
		}
		defer f.Close()
		input = f
	case *flagURL != "":
		resp, err := http.Get(*flagURL)
		if err != nil {
			fmt.Println("error getting url:", *flagURL, err)
			os.Exit(124)
		}
		defer resp.Body.Close()
		input = resp.Body
	default:
		input = os.Stdin
	}
	printMD5(input, os.Stderr)
}

func printMD5(r io.Reader, w io.Writer) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		fmt.Println("error copying stdin to hash")
		os.Exit(1)
	}

	fmt.Fprintf(w, "%x", h.Sum(nil))
}
