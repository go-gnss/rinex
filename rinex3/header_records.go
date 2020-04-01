package rinex3

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Header record descriptors in columns 61-80 are mandatory

// Free ordering of Header section, with Exceptions:
// RINEX VERSION / TYPE record MUST always be the first record in a file
// SYS / # / OBS TYPES record(s) should precede any SYS / DCBS APPLIED and SYS / SCALE FACTOR records
// # OF SATELLITES record (if present) should be immediately followed by the corresponding number of PRN / # OF OBS records.
// 		These records may be handy for documentary purposes. However, since they may only be created after having read the whole raw data file, we define them to be optional
// END OF HEADER record MUST be the last record in the header

type HeaderRecord struct {
	Value string
	Key   string
	Line  int
}

type HeaderRecordParser func(*bufio.Reader, *Header, HeaderRecord) error

// TODO: Check for empty strings / missing required values?
var (
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
			h.Comments = append(h.Comments, HeaderComment{hr.Value[:60], hr.Line})
			return nil
		},
		"MARKER NAME": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Marker.Name = strings.TrimSpace(hr.Value)
			return nil
		},
		"MARKER NUMBER": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Marker.Number = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"MARKER TYPE": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Marker.Type = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"OBSERVER / AGENCY": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Observer = strings.TrimSpace(hr.Value[:20])
			h.Agency = strings.TrimSpace(hr.Value[20:])
			return nil
		},
		"REC # / TYPE / VERS": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Receiver.Number = strings.TrimSpace(hr.Value[:20])
			h.Receiver.Type = strings.TrimSpace(hr.Value[20:40])
			h.Receiver.Version = strings.TrimSpace(hr.Value[40:])
			return nil
		},
		"ANT # / TYPE": func(_ *bufio.Reader, h *Header, hr HeaderRecord) error {
			h.Antenna.Number = strings.TrimSpace(hr.Value[:20])
			h.Antenna.Type = strings.TrimSpace(hr.Value[20:40])
			return nil
		},
		"APPROX POSITION XYZ": func(_ *bufio.Reader, h *Header, hr HeaderRecord) (err error) {
			h.Marker.ApproxPosition.X, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[:14]), 64)
			if err != nil {
				return err
			}
			h.Marker.ApproxPosition.Y, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[14:28]), 64)
			if err != nil {
				return err
			}
			h.Marker.ApproxPosition.Z, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[28:42]), 64)
			return err
		},
		"ANTENNA: DELTA H/E/N": func(_ *bufio.Reader, h *Header, hr HeaderRecord) (err error) {
			h.Antenna.Height, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[:14]), 64)
			if err != nil {
				return err
			}
			h.Antenna.East, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[14:28]), 64)
			if err != nil {
				return err
			}
			h.Antenna.North, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[28:42]), 64)
			return err
		},
		"SYS / # / OBS TYPES": func(r *bufio.Reader, h *Header, hr HeaderRecord) (err error) {
			// TODO: Fix this, it's insane and will fail in the event of multiple continuation lines
			// Probably better off having HeaderRecordParser take the Reader (added) so it can read
			// ahead for any records with continuation lines, getting rid of Header.records
			// Still need to keep track of line numbers though - maybe need a Scanner Object
			pattern := regexp.MustCompile(`^([A-Z])..([ 0-9][ 0-9][0-9])|(?: ([A-Z0-9]{3}))`)
			matchLine := pattern.FindAllStringSubmatch(hr.Value, -1)
			if len(matchLine) > 0 {
				system := matchLine[0][1]
				if system != "" {
					for _, obsType := range matchLine[1:] {
						h.ObservationTypes[system] = append(h.ObservationTypes[system], obsType[3])
					}
					return nil
				} else {
					// continuation line
					matchPrevious := pattern.FindStringSubmatch(h.records[len(h.records)-2].Value)
					if len(matchPrevious) > 0 {
						system = matchPrevious[1]
						for _, obsType := range matchLine[0:] {
							h.ObservationTypes[system] = append(h.ObservationTypes[system], obsType[3])
						}
						return nil
					}
				}
			}
			return errors.New("Failed to parse SYS / # / OBS TYPES field")
		},
	}
)

func ParseHeaderRecords(reader *bufio.Reader) (header Header, err error) {
	header = NewHeader()
	currentLine := 1
	line, err := reader.ReadString('\n')
	// TODO: First line MUST be RINEX VERSION / TYPE
	for ; err == nil; line, err = reader.ReadString('\n') {
		line = strings.TrimRight(line, " \n")

		if len(line) < 61 || len(line) > 80 {
			return header, NewInvalidHeaderRecordError(line, currentLine)
		}

		hr := HeaderRecord{line[:60], line[60:], currentLine}
		header.records = append(header.records, hr)

		if parser, ok := HeaderRecordParsers[hr.Key]; ok {
			parser(reader, &header, hr)
		} else {
			return header, NewInvalidHeaderRecordError(line, currentLine)
		}

		if line[60:] == "END OF HEADER" {
			return header, err
		}
		currentLine++
	}
	return header, err
}

type InvalidHeaderRecord error

func NewInvalidHeaderRecordError(line string, lineNumber int) InvalidHeaderRecord {
	// TODO: Add reason for error, i.e. InvalidHeaderLabel for example
	return errors.New(fmt.Sprintf("Invalid header record found at line %d: \"%s\"", lineNumber, line))
}
