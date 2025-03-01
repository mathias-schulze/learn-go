package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
)

var (
	flagInFolder  = flag.String("in", "./", "input folder")
	flagOutFolder = flag.String("out", "", "output folder")
	flagSize      = flag.String("size", "500x500", "output size")
)

func main() {
	flag.Parse()

	size, err := parseSize(*flagSize)
	if err != nil {
		fmt.Println("cannot create picSize: ", err)
		os.Exit(1)
	}

	outFolder := *flagSize
	if *flagOutFolder != "" {
		outFolder = *flagOutFolder
	}

	err = resizeFolderImages(*flagInFolder, outFolder, size)
	if err != nil {
		fmt.Println(err)
		os.Exit(10)
	}
}

type errorList struct {
	errs []error
}

func (e *errorList) add(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

func (e *errorList) hasError() bool {
	return len(e.errs) > 0
}

func (e *errorList) Error() string {
	if !e.hasError() {
		return ""
	}
	out := fmt.Sprintf("number of errors: %d", len(e.errs))
	for i, err := range e.errs {
		out = fmt.Sprintf("%s\n%d: %s", out, i, err.Error())
	}
	return out
}

func resizeFolderImages(inFolder, outFolder string, size picSize) error {
	err := os.MkdirAll(outFolder, 0777)
	if err != nil {
		return fmt.Errorf("cannot create output folder: %w", err)
	}

	dir, err := os.ReadDir(inFolder)
	if err != nil {
		return fmt.Errorf("cannot read from input folder: %w", err)
	}

	wg := &sync.WaitGroup{}
	errList := &errorList{}
	errChan := make(chan error, 1)
	resizeChan := make(chan resizeArgs)
	wg.Add(3)
	go resizer(wg, resizeChan, errChan)
	go resizer(wg, resizeChan, errChan)
	go resizer(wg, resizeChan, errChan)

	go func(err *errorList, errChan chan error) {
		for err := range errChan {
			errList.add(err)
		}
	}(errList, errChan)

	for _, de := range dir {
		if de.IsDir() || !useFile(de.Name()) {
			continue
		}

		inPath := filepath.Join(inFolder, de.Name())
		outPath := filepath.Join(outFolder, de.Name())
		resizeChan <- resizeArgs{inPath, outPath, size}
	}

	close(resizeChan)
	close(errChan)
	wg.Wait()

	if errList.hasError() {
		return errList
	}

	return nil
}

type picSize struct {
	width, height int
}

func parseSize(s string) (picSize, error) {
	var ps picSize
	var err error

	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return ps, fmt.Errorf("%s does not fit to widthxheight", s)
	}

	ps.width, err = strconv.Atoi(parts[0])
	if err != nil {
		return ps, fmt.Errorf("parseSize: ps.x: %w", err)
	}

	ps.height, err = strconv.Atoi(parts[1])
	if err != nil {
		return ps, fmt.Errorf("parseSize: ps.y: %w", err)
	}

	return ps, nil
}

type resizeArgs struct {
	inPath  string
	outPath string
	size    picSize
}

func resizer(wg *sync.WaitGroup, c chan resizeArgs, errChan chan error) {
	for a := range c {
		log.Println("resize: ", a.inPath)
		inFile, err := os.Open(a.inPath)
		if err != nil {
			errChan <- fmt.Errorf("error opening file: %w", err)
			continue
		}
		outFile, err := os.OpenFile(a.outPath, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			errChan <- err
			inFile.Close()
			continue
		}
		err = resizeClose(a.size, inFile, outFile)
		if err != nil {
			errChan <- err
		}
	}
	wg.Done()
}

func resizeClose(ps picSize, r io.ReadCloser, w io.WriteCloser) error {
	defer r.Close()
	defer w.Close()
	return resize(ps, r, w)
}

func resize(ps picSize, r io.Reader, w io.Writer) error {
	img, format, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("error decoding: %w", err)
	}
	if format != "jpeg" {
		return fmt.Errorf("just jpeg supported, no %s", format)
	}

	resized := imaging.Fit(img, ps.width, ps.height, imaging.Lanczos)
	return jpeg.Encode(w, resized, nil)
}

func useFile(filename string) bool {
	allowed := []string{".jpg", ".jpeg"}
	ext := filepath.Ext(filename)
	for _, e := range allowed {
		if strings.EqualFold(ext, e) {
			return true
		}
	}
	return false
}
