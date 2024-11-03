package main

import (
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		log.Println("min 1 arg")
		os.Exit(1)
	}
	for _, arg := range os.Args[1:] {
		fd, err := os.Open(arg)
		if err != nil {
			log.Println("error opening file", arg)
			os.Exit(2)
		}
		_, err = io.Copy(os.Stdout, fd)
		if err != nil {
			log.Println("error with io.Copy()", arg)
			os.Exit(3)
		}
		fd.Close()
	}
}
