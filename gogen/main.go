package main

import (
	"fmt"
	"os"
	"text/template"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("gogen needs exactly 2 arguments!")
		os.Exit(1)
	}

	templateFileName := os.Args[1]
	typeName := os.Args[2]

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		fmt.Println("error parsing template: ", err)
		os.Exit(2)
	}

	outName := fmt.Sprintf("gogen_%s_gen.go", typeName)
	fd, err := os.OpenFile(outName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error creating file:", err)
		os.Exit(3)
	}
	defer fd.Close()

	data := struct {
		T string
	}{
		typeName,
	}

	t.Execute(fd, data)
}
