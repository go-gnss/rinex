package rinex

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/go-gnss/rinex/header"
	"github.com/go-gnss/rinex/rinex3"
	"github.com/go-gnss/rinex/scanner"
)

// TODO: Implement differentiation between RINEX 2 and 3
// TODO: Implement RinexFileName

type RinexHeader interface {
	GetFormatVersion() float64
	GetFileType() string // TODO: FileType type
}

type RinexFile struct {
	scanner *scanner.Scanner
	Header  RinexHeader
}

// TODO: Header gives RinexVersion and FileType, consider implementation
// of Rinex3ObservationFile, Rinex2NavigationFile, etc

// TODO: Can scan through observations like so, or parse them into a map
//func (r RinexFile) NextEpoch() (epoch EpochRecord, err error) {
//	return epoch, err
//}

func ParseRinexFile(data io.Reader) (file RinexFile, err error) {
	scanner := &scanner.Scanner{bufio.NewReader(data), 0}
	header, err := ParseHeader(scanner)
	file = RinexFile{
		scanner: scanner,
		Header:  header,
	}
	if err != nil {
		return file, err
	}

	for err == nil {
		//var epoch rinex3.EpochRecord
		_, err = rinex3.ParseEpochRecord(scanner, file.Header.(rinex3.ObservationHeader).ObservationTypes)
		//fmt.Println(epoch, err)
	}
	if err != io.EOF {
		return file, err
	}
	return file, nil
}

// TODO: Check for empty strings / missing required values?
func ParseHeader(scanner *scanner.Scanner) (rinexHeader RinexHeader, err error) {
	hr, err := header.ParseHeaderRecord(scanner)
	if err != nil {
		return rinexHeader, header.NewHeaderRecordParsingError(err, scanner.Line)
	}

	// TODO: This isn't true for CRX files, but that is not reflected in the format
	// description - though it does mention .crx extensions are allowed in the
	// filename
	if hr.Key != "RINEX VERSION / TYPE" {
		return rinexHeader, errors.New("first line of header must be \"RINEX VERSION / TYPE\"")
	}

	h := header.Header{}
	err = header.HeaderRecordParsers[hr.Key](scanner, &h, hr)
	if err != nil {
		return rinexHeader, header.NewHeaderRecordParsingError(err, scanner.Line)
	}

	// TODO: NavigationHeader and MeteorologicalHeader
	// TODO: RINEX 2 and 3
	switch h.FileType {
	case "O":
		obsHeader := rinex3.NewObservationHeader(h)
		err = rinex3.ParseObservationHeader(scanner, &obsHeader)
		return obsHeader, err
	default:
		return rinexHeader, errors.New(fmt.Sprintf("invalid header type \"%v\"", h.FileType))
	}
}
