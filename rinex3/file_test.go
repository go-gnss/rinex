package rinex3_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-gnss/rinex/rinex3"
)

func TestParseFile(t *testing.T) {
	file, err := os.Open("ALBY00AUS_R_20183280000_01D_30S_MO.rnx")
	if err != nil {
		t.Error("Failed to open test data file")
	}

	fmt.Println(rinex3.Parse(file))
}
