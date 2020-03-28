package rinex_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-gnss/rinex/rinex3"
)

func TestParseFile(t *testing.T) {
	file, err := os.Open("data")
	if err != nil {
		t.Error("Failed to open test data file")
	}

	fmt.Println(rinex3.Parse(file))
}
