package rinex_test

import (
	"testing"

	"github.com/go-gnss/rinex"
)

// TODO: Test RinexFilename attributes
// TODO: Format logs better

func TestObservationLongName(t *testing.T) {
	_, err := rinex.NewRinexFilename("SITE00AUS_R_20183280000_01D_30S_MO.rnx")
	if err != nil {
		t.Error(err)
	}

	_, err = rinex.NewRinexFilename("SITE00AUS_R_20183280000_01D_30S_MO.rnx.gz")
	if err != nil {
		t.Error(err)
	}
}

func TestNavigationLongName(t *testing.T) {
	_, err := rinex.NewRinexFilename("SITE00AUS_R_20183280000_01D_MN.rnx")
	if err != nil {
		t.Error(err)
	}

	_, err = rinex.NewRinexFilename("SITE00AUS_R_20183280000_01D_MN.rnx.gz")
	if err != nil {
		t.Error(err)
	}
}

func TestInvalidFilename(t *testing.T) {
	_, err := rinex.NewRinexFilename("SITE00AUS_R_201832800000_01D_MN.rnx.gz")
	if err.Error() != "invalid RINEX filename" {
		t.Error("NewRinexFilename did not return invalid RINEX filename error")
	}
}
