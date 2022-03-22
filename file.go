package rinex

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/go-gnss/rinex/header"
	"github.com/go-gnss/rinex/rinex3"
)

type RinexFile interface {
	GetFormatVersion() float64
	GetFileType() string
}

func ParseRinexFile(data io.Reader) (file RinexFile, err error) {
	r := bufio.NewReader(data)
	header, err := ParseHeader(r)
	if err != nil {
		return file, err
	}

	// TODO: RINEX 2
	// TODO: NavigationHeader and MeteorologicalHeader
	switch header.FileType {
	case "O":
		obsHeader := rinex3.NewObservationHeader(header)
		err = rinex3.ParseObservationHeader(r, &obsHeader)
		obsFile := rinex3.ObservationFile{ObservationHeader: obsHeader, Epochs: []rinex3.EpochRecord{}}
		obsFile.Epochs, err = rinex3.ParseEpochRecords(r, obsFile.ObservationTypes)
		return obsFile, err
	default:
		return file, fmt.Errorf("invalid or unsupported file type %q", header.GetFileType())
	}
}

// TODO: Check for empty strings / missing required values?
func ParseHeader(r *bufio.Reader) (rinexHeader header.Header, err error) {
	hr, err := header.ParseHeaderRecord(r)
	if err != nil {
		return rinexHeader, fmt.Errorf("error parsing header record: %e", err)
	}

	if hr.Key != "RINEX VERSION / TYPE" {
		return rinexHeader, errors.New("first line of header must be \"RINEX VERSION / TYPE\"")
	}

	h := header.Header{}
	err = header.HeaderRecordParsers[hr.Key](r, &h, hr)
	if err != nil {
		return rinexHeader, fmt.Errorf("error parsing header record: %e", err)
	}

	return h, nil
}
