package rinex_test

import (
	"os"
	"testing"

	"github.com/go-gnss/rinex"
	"github.com/go-gnss/rinex/rinex3"
)

func TestParseObservationFile(t *testing.T) {
	//file, err := os.Open("fixtures/ALIC00AUS_S_20220760315_15M_01S_MO.rnx")
	file, err := os.Open("fixtures/ALBY00AUS_R_20183280000_01D_30S_MO.rnx")
	if err != nil {
		t.Fatal("failed to open test observation file")
	}

	rinexFile, err := rinex.ParseRinexFile(file)
	if err != nil {
		t.Fatal(err.Error())
	}

	if ft := rinexFile.GetFileType(); ft != "O" {
		t.Errorf("incorrect RINEX File Type: %s", ft)
	}

	if fv := rinexFile.GetFormatVersion(); fv != 3.03 {
		t.Errorf("incorrect RINEX Format Version: %f", fv)
	}

	if _, ok := rinexFile.(rinex3.ObservationFile); !ok {
		t.Errorf("couldn't cast RinexFile interface to rinex3.ObservationFile")
	}

	// TODO: Test header attributes
}
