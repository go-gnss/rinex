package rinex_test

import (
	"os"
	"testing"

	"github.com/go-gnss/rinex"
	"github.com/go-gnss/rinex/rinex3"
)

func TestParseObservationFile(t *testing.T) {
	file, err := os.Open("data/ALBY00AUS_R_20183280000_01D_30S_MO.rnx")
	if err != nil {
		t.Error("Failed to open test data file")
	}

	rinexFile, err := rinex.ParseRinexFile(file)
	if err != nil {
		t.Errorf(err.Error())
	}

	if ft := rinexFile.Header.GetFileType(); ft != "O" {
		t.Errorf("incorrect RINEX File Type: %s", ft)
	}

	if fv := rinexFile.Header.GetFormatVersion(); fv != 3.03 {
		t.Errorf("incorrect RINEX Format Version: %f", fv)
	}

	if _, ok := rinexFile.Header.(rinex3.ObservationHeader); !ok {
		t.Errorf("couldn't cast RinexHeader interface to ObservationHeader")
	}

	// TODO: Test header attributes
}
