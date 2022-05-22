package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func Test_getInput(t *testing.T) {
	tests := []struct {
		name    string
		want    input
		wantErr bool
		osArgs  []string
	}{
		{"Default parameters", input{"test.csv", ";"}, false, []string{"cmd", "test.csv"}},
		{"No parameters", input{}, true, []string{"cmd"}},
		{"Separator specified", input{"test.csv", ";"}, false, []string{"cmd", "--separator=;", "test.csv"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Saving the original os.Args reference
			actualOsArgs := os.Args
			// This defer function will run after the test is done
			defer func() {
				// Restoring the original os.Args reference
				os.Args = actualOsArgs
				// Reseting the Flag command line. So that we can parse flags again
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tt.osArgs // Setting the specific command args for this

			got, err := getInput()
			if (err != nil) != tt.wantErr {
				t.Errorf("getInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name   string
		osArgs []string
	}{
		{"With file parameter", []string{"cmd", "test.csv"}},
		{"No parameters", []string{"cmd"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Saving the original os.Args reference
			actualOsArgs := os.Args
			// This defer function will run after the test is done
			defer func() {
				// Restoring the original os.Args reference
				os.Args = actualOsArgs
				// Reseting the Flag command line. So that we can parse flags again
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tt.osArgs // Setting the specific command args for this

			main()
		})
	}
}

func Test_checkIfValidFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		{"File does exist", "importer.go", true, false},
		{"File does not exist", "nowhere/test.csv", false, true},
		{"File is not csv", "test.txt", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfValidFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfValidFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkIfValidFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
