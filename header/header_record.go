package header

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Header Record Descriptors in columns 61-80 are mandatory

type HeaderRecordParser func(*bufio.Reader, *Header, HeaderRecord) error

var (
	// TODO: Consider making these attributes a part of RinexFile, instead of having Header
	// type - also Comments can appear anywhere in a file
	HeaderRecordParsers map[string]HeaderRecordParser = map[string]HeaderRecordParser{
		"RINEX VERSION / TYPE": func(_ *bufio.Reader, h *Header, hr HeaderRecord) (err error) {
			h.FormatVersion, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[:9]), 64)
			h.FileType = string(hr.Value[20])
			h.SatelliteSystem = string(hr.Value[40])
			return err
		},
		"PGM / RUN BY / DATE": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Program = strings.TrimSpace(hr.Value[:20])
			h.RunBy = strings.TrimSpace(hr.Value[20:40])
			h.CreationDate = strings.TrimSpace(hr.Value[40:])
			return nil
		},
		"COMMENT": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Comments = append(h.Comments, hr.Value[:60])
			return nil
		},
		"END OF HEADER": func(s *bufio.Reader, h *Header, hr HeaderRecord) error {
			return nil
		},
	}
)

type HeaderRecord struct {
	Value string
	Key   string
}

func ParseHeaderRecord(r *bufio.Reader) (hr HeaderRecord, err error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return hr, err
	}

	line = strings.TrimRight(line, " \n")
	if len(line) < 61 || len(line) > 80 {
		return hr, errors.New(fmt.Sprintf("invalid header line \"%s\"", line))
	}

	return HeaderRecord{line[:60], line[60:]}, err
}
