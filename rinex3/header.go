package rinex3

import (
	"errors"
	"fmt"
)

// Free ordering of Header section, with Exceptions:
// RINEX VERSION / TYPE record MUST always be the first record in a file
// SYS / # / OBS TYPES record(s) should precede any SYS / DCBS APPLIED and SYS / SCALE FACTOR records
// # OF SATELLITES record (if present) should be immediately followed by the corresponding number of PRN / # OF OBS records.
// 		These records may be handy for documentary purposes. However, since they may only be created after having read the whole raw data file, we define them to be optional
// END OF HEADER record MUST be the last record in the header

type RinexHeader interface {
	GetFormatVersion() float64
	GetFileType() string // TODO: FileType type
}

// TODO: Consider having Header just be a slice of HeaderRecord, using
// Getters and Setters for each attribute - wouldn't need separate
// types for Obs, Nav, Met, but would need a lot of error checking
type Header struct {
	FormatVersion   float64
	FileType        string
	SatelliteSystem string
	Program         string
	RunBy           string
	CreationDate    string // TODO: time.Time
	Comments        []HeaderComment
}

func (h Header) GetFormatVersion() float64 {
	return h.FormatVersion
}

func (h Header) GetFileType() string {
	return h.FileType
}

// TODO: Check for empty strings / missing required values?
func ParseHeader(scanner *Scanner) (header RinexHeader, err error) {
	hr, err := ParseHeaderRecord(scanner)
	if err != nil {
		return header, NewHeaderRecordParsingError(err, scanner.line)
	}

	if hr.Key != "RINEX VERSION / TYPE" {
		return header, errors.New("first line of header must be \"RINEX VERSION / TYPE\"")
	}

	h := Header{}
	HeaderRecordParsers[hr.Key](scanner, &h, hr)
	switch h.FileType {
	case "O":
		obsHeader := NewObservationHeader(h)
		err = ParseObservationHeader(scanner, &obsHeader)
		return obsHeader, err
	// TODO: NavigationHeader and MeteorologicalHeader
	default:
		return header, errors.New(fmt.Sprintf("invalid header type \"%v\"", h.FileType))
	}
}

type HeaderComment struct {
	Comment string
	Line    int // TODO: This might not be useful for reconstructing Headers if additional lines are added
}
