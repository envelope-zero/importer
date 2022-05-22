package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type input struct {
	filepath  string
	separator string
}

func main() {
	_, err := getInput()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getInput() (input, error) {
	if len(os.Args) < 2 {
		return input{}, errors.New("A filepath argument is required")
	}

	s := flag.String("separator", ";", "Column separator")

	// Parse all flags
	flag.Parse()

	p := flag.Arg(0)

	return input{p, *s}, nil
}

func checkIfValidFile(filename string) (bool, error) {
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("File %s does not exist", filename)
	}

	return true, nil
}
