package header

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-gnss/rinex/scanner"
)

// Header Record Descriptors in columns 61-80 are mandatory

type HeaderRecordParser func(*scanner.Scanner, *Header, HeaderRecord) error

var (
	// TODO: Consider making these attributes a part of RinexFile, instead of having Header
	// type - also Comments can appear anywhere in a file
	HeaderRecordParsers map[string]HeaderRecordParser = map[string]HeaderRecordParser{
		"RINEX VERSION / TYPE": func(_ *scanner.Scanner, h *Header, hr HeaderRecord) (err error) {
			h.FormatVersion, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[:9]), 64)
			h.FileType = string(hr.Value[20])
			h.SatelliteSystem = string(hr.Value[40])
			return err
		},
		"PGM / RUN BY / DATE": func(_ *scanner.Scanner, h *Header, hr HeaderRecord) error {
			h.Program = strings.TrimSpace(hr.Value[:20])
			h.RunBy = strings.TrimSpace(hr.Value[20:40])
			h.CreationDate = strings.TrimSpace(hr.Value[40:])
			return nil
		},
		"COMMENT": func(_ *scanner.Scanner, h *Header, hr HeaderRecord) error {
			h.Comments = append(h.Comments, HeaderComment{hr.Value[:60], hr.Line})
			return nil
		},
		"END OF HEADER": func(s *scanner.Scanner, h *Header, hr HeaderRecord) error {
			return nil
		},
	}
)

type HeaderRecord struct {
	Value string
	Key   string
	Line  int
}

func ParseHeaderRecord(scanner *scanner.Scanner) (hr HeaderRecord, err error) {
	line, err := scanner.ReadLine()
	if err != nil {
		return hr, err
	}

	line = strings.TrimRight(line, " \n")
	if len(line) < 61 || len(line) > 80 {
		return hr, errors.New(fmt.Sprintf("invalid header line \"%s\"", line))
	}

	return HeaderRecord{line[:60], line[60:], scanner.Line}, err
}

type HeaderRecordParsingError error

func NewHeaderRecordParsingError(err error, lineNumber int) HeaderRecordParsingError {
	return errors.New(fmt.Sprintf("failed to parse header record at line %d with reason: %s", lineNumber, err.Error()))
}
