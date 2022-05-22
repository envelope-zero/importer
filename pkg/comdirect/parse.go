package comdirect

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/envelope-zero/backend/pkg/models"
	"github.com/envelope-zero/importer/internal/util"
	"github.com/envelope-zero/importer/pkg/types"
	"golang.org/x/text/encoding/charmap"
)

const (
	Buchungstag uint = iota
	Wertstellung
	Vorgang
	Buchungstext
	Umsatz
)

// This function parses the comdirect CSV files
func Parse(f *os.File, ch chan<- types.ResourceCerate) {
	// Set up the CSV reader
	r := csv.NewReader(charmap.ISO8859_15.NewDecoder().Reader(f))
	r.Comma = ';'

	// Do not check the fields per record as it varies
	r.FieldsPerRecord = -1

	// The comdirect CSV files have unneeded data in the beginning.
	// We loop over it until we find the "Buchungstag" as the first element,
	// which is the first column of the header line
	var h []string
	for {
		var err error
		h, err = r.Read()

		if err != nil {
			// The file does not contain the comdirect header
			if err == io.EOF {
				util.ExitGracefully(errors.New("The file does not contain a header line starting with \"Buchungstext\". It does not seem to be a valid comdirect CSV file."))
			}
			util.ExitGracefully(err)
		}

		if h[Buchungstag] == "Buchungstag" {
			break
		}
	}

	for {
		l, err := r.Read()

		// When we arrived at the end of the file
		if err == io.EOF {
			close(ch)
			break
		} else if err != nil {
			util.ExitGracefully(err)
		}

		d, err := parseLine(h, l)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("DEBUG: Adding to queue: %v\n", d)
		ch <- types.ResourceCerate{
			Resource: d,
		}
	}
}

// parseLine converts a single line to a transaction
func parseLine(h []string, l []string) (interface{}, error) {
	// Check if the line is a valid transaction by matching field amounts
	if len(l) != len(h) {
		return types.Transaction{}, fmt.Errorf("DEBUG: Not a valid transaction: %s\n", l)
	}

	// Parse the date of the transaction
	date, err := time.Parse("02.01.2006", l[Wertstellung])
	if err != nil {
		return types.Transaction{}, fmt.Errorf("DEBUG: Could not parse date, skipping: %s\n", l)
	}

	// Parse the amount of the transaction form German number formatting
	amount, err := util.NormalizeGerman(l[Umsatz])
	if err != nil {
		return types.Transaction{}, fmt.Errorf("WARNING: Amount not parseable, skipping: %s\n", l)
	}

	// Get the reference field
	note, err := extractField(l[Buchungstext], fields["Buchungstext"])
	if err != nil {
		return types.Transaction{}, err
	}

	// Extract the initator (Auftraggeber)
	auftraggeber, err := extractField(l[Buchungstext], fields["Auftraggeber"])
	if err != nil {
		return types.Transaction{}, err
	}

	// Extract the recipient (Empfänger)
	empfaenger, err := extractField(l[Buchungstext], fields["Empfänger"])
	if err != nil {
		return types.Transaction{}, err
	}

	// Get the Empfänger or Auftraggeber. If both are defined, the format changed,
	// in this case we throw an error
	o := auftraggeber
	if auftraggeber != "" && empfaenger != "" {
		return types.Transaction{}, fmt.Errorf("WARNING: Transaction has both Auftraggeber and Empfänger defined, this indicates a change of CSV format. Skipping: %s\n", l)
	} else if empfaenger != "" {
		o = empfaenger
	}

	r := types.Transaction{
		Create: models.TransactionCreate{
			Date: date,
			Note: note,
		},
	}

	// When the amount is negative, the destination account is the Auftraggeber/Empfänger
	r.Create.Amount = amount.Abs()
	if amount.IsNegative() {
		r.DestinationAccount = o
	} else {
		r.SourceAccount = o
	}

	// Returning our generated map
	return r, nil
}

// fields contains all strings that start a field in the Buchungstext CSV
var fields = map[string]string{
	"Buchungstext":        "Buchungstext",
	"Empfänger":           "Empfänger",
	"Auftraggeber":        "Auftraggeber",
	"Zahlungspflichtiger": "Zahlungspflichtiger",
	"Kto/IBAN":            "Kto/IBAN",
	"BLZ/BIC":             "BLZ/BIC",
	"Ref":                 "Ref",
}

// extractField extracts a specific field from the Buchungstext CSV.
// This code mostly a go reimplementation of some code by leolabs: https://github.com/leolabs/you-need-a-parser/blob/49bd27c34dcc3243bbbbba15599982385ebb2fa7/packages/ynap-parsers/src/de/comdirect/comdirect.ts#L52-L78
func extractField(s, t string) (string, error) {
	_, ok := fields[t]
	if !ok {
		return "", errors.New("ERROR: Tried to extract field that does not exist")
	}

	// Split by the searched field, this will either
	// tell us that the field does not exist or give
	// us the field with other garbage at the end
	split := strings.Split(s, t)
	if len(split) < 2 {
		return "", nil
	}

	// Remove all other fields behind the one we're looking for
	re := regexp.MustCompile("(Buchungstext|Empfänger|Auftraggeber|Zahlungspflichtiger|Kto\\/IBAN|BLZ\\/BIC|Ref\\.)")
	split = re.Split(split[1], -1)

	// Remove whitespace and other leftovers
	re = regexp.MustCompile("(^[:.\\s]+|\\s+$)")

	return re.ReplaceAllString(split[0], ""), nil
}
