package util

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNormalizeGerman(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Amount without separators", args{"1000"}, "1000", false},
		{"Amount with thousands separator", args{"2.713"}, "2713", false},
		{"Amount with thousands separator and decimal separator", args{"1.337,57"}, "1337.57", false},
		{"Amount with decimal comma", args{"13,37"}, "13.37", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, err := decimal.NewFromString(tt.want)
			if err != nil {
				t.Errorf("NormalizeGerman() test error: could not parse want to decimal: %v", err)
			}

			got, err := NormalizeGerman(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeGerman() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("NormalizeGerman() = %v, want %v", got, want)
			}
		})
	}
}
