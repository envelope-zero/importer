package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/envelope-zero/importer/internal/util"
	"github.com/envelope-zero/importer/pkg/api"
	"github.com/envelope-zero/importer/pkg/comdirect"
	"github.com/envelope-zero/importer/pkg/types"
	"github.com/google/uuid"
)

type input struct {
	path      string
	dryRun    bool
	accountId uuid.UUID
	t         string // The type. Can't be named type because that is a reserved keyword
}

func main() {
	// Showing useful information when the user enters the --help option
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <file>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse input flags
	i, err := getInput()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Check if the file exists and is accessible
	if _, err := os.Stat(i.path); err != nil && os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", i.path)
		return
	}

	// Open file to parse
	f, err := os.Open(i.path)
	if err != nil {
		util.ExitGracefully(err)
	}
	defer f.Close()

	ch := make(chan types.ResourceCerate)
	done := make(chan bool)

	go api.CreateResources(i.accountId, ch, done)

	switch i.t {
	case "comdirect":
		go comdirect.Parse(f, ch)
	}

	<-done
}

// getInput parses the input and validates it
func getInput() (input, error) {
	if len(os.Args) < 2 {
		return input{}, errors.New("A filepath argument is required")
	}

	t := flag.String("type", "", "Type of file you want to import, e.g. 'comdirect'")
	dryRun := flag.Bool("dry-run", true, "Don’t import transactions to Envelope Zero, only show what would be done.")
	accountId := flag.String("account-id", "", "The ID of the account to import to")

	// Parse all flags
	flag.Parse()
	p := flag.Arg(0)

	if !(*t == "comdirect") {
		return input{}, errors.New("type must be one of 'comdirect'")
	}

	id, err := uuid.Parse(*accountId)
	if err != nil {
		return input{}, errors.New("the account ID must be a valid uuid")
	}

	return input{path: p, dryRun: *dryRun, t: *t, accountId: id}, nil
}
