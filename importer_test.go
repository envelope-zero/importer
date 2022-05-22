package main

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name   string
		osArgs []string
	}{
		{"With file parameter", []string{"cmd", "test.csv"}},
		{"Type specified", []string{"cmd", "--type=comdirect", "test.csv"}},
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

func Test_getInput(t *testing.T) {
	tests := []struct {
		name    string
		want    input
		wantErr bool
		osArgs  []string
	}{
		// Test UUID: bbbfb572-1e23-439c-a5ab-c47fce2baa1f
		// uuid.UUID{0xbb, 0xbf, 0xb5, 0x72, 0x1e, 0x23, 0x43, 0x9c, 0xa5, 0xab, 0xc4, 0x7f, 0xce, 0x2b, 0xaa, 0x1f}
		{"No parameters", input{}, true, []string{"cmd"}},
		{"Type specified", input{path: "test.csv", dryRun: true, t: "comdirect", accountId: uuid.UUID{0xbb, 0xbf, 0xb5, 0x72, 0x1e, 0x23, 0x43, 0x9c, 0xa5, 0xab, 0xc4, 0x7f, 0xce, 0x2b, 0xaa, 0x1f}}, false, []string{"cmd", "--type=comdirect", "--account-id=bbbfb572-1e23-439c-a5ab-c47fce2baa1f", "test.csv"}},
		{"Invalid type specified", input{}, true, []string{"cmd", "--type=n26", "test.csv"}},
		{"No dry run", input{path: "test.csv", dryRun: false, t: "comdirect", accountId: uuid.UUID{0xbb, 0xbf, 0xb5, 0x72, 0x1e, 0x23, 0x43, 0x9c, 0xa5, 0xab, 0xc4, 0x7f, 0xce, 0x2b, 0xaa, 0x1f}}, false, []string{"cmd", "--dry-run=false", "--account-id=bbbfb572-1e23-439c-a5ab-c47fce2baa1f", "--type=comdirect", "test.csv"}},
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
