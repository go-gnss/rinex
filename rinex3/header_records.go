package rinex3

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Header Record Descriptors in columns 61-80 are mandatory

type HeaderRecord struct {
	Value string
	Key   string
	Line  int
}

type HeaderRecordParser func(*Scanner, *Header, HeaderRecord) error

var (
	HeaderRecordPatternError error = errors.New("header record did not match pattern")
	// TODO: Consider moving function definitions out of the map
	// TODO: Consider using a proper parser library again, or just regex groups w/
	// Getters and Setters on Header
	HeaderRecordParsers map[string]HeaderRecordParser = map[string]HeaderRecordParser{
		"RINEX VERSION / TYPE": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			h.FormatVersion, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[:9]), 64)
			h.FileType = string(hr.Value[20])
			h.SatelliteSystem = string(hr.Value[40])
			return err
		},
		"PGM / RUN BY / DATE": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Program = strings.TrimSpace(hr.Value[:20])
			h.RunBy = strings.TrimSpace(hr.Value[20:40])
			h.CreationDate = strings.TrimSpace(hr.Value[40:])
			return nil
		},
		"COMMENT": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Comments = append(h.Comments, HeaderComment{hr.Value[:60], hr.Line})
			return nil
		},
		"MARKER NAME": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Marker.Name = strings.TrimSpace(hr.Value)
			return nil
		},
		"MARKER NUMBER": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Marker.Number = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"MARKER TYPE": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Marker.Type = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"OBSERVER / AGENCY": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Observer = strings.TrimSpace(hr.Value[:20])
			h.Agency = strings.TrimSpace(hr.Value[20:])
			return nil
		},
		"REC # / TYPE / VERS": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Receiver.Number = strings.TrimSpace(hr.Value[:20])
			h.Receiver.Type = strings.TrimSpace(hr.Value[20:40])
			h.Receiver.Version = strings.TrimSpace(hr.Value[40:])
			return nil
		},
		"ANT # / TYPE": func(_ *Scanner, h *Header, hr HeaderRecord) error {
			h.Antenna.Number = strings.TrimSpace(hr.Value[:20])
			h.Antenna.Type = strings.TrimSpace(hr.Value[20:40])
			return nil
		},
		"APPROX POSITION XYZ": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
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
		"ANTENNA: DELTA H/E/N": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
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
		"ANTENNA: DELTA X/Y/Z": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err
		},
		"ANTENNA: PHASECENTER": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"ANTENNA: B.SIGHT XYZ": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"ANTENNA: ZERODIR AZI": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"ANTENNA: ZERODIR XYZ": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"CENTER OF MASS: XYZ": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / # / OBS TYPES": func(s *Scanner, h *Header, hr HeaderRecord) (err error) {
			linePattern := regexp.MustCompile(`^([A-Z])..([ 0-9][ 0-9][0-9])|( ([A-Z0-9]{3}))`)
			continuationLinePattern := regexp.MustCompile(`^( {6})|( ([A-Z0-9]{3}))`)

			matchLine := linePattern.FindAllStringSubmatch(hr.Value, -1)
			if len(matchLine) < 3 {
				return HeaderRecordPatternError
			}

			system := matchLine[0][1]
			totalObs, err := strconv.ParseInt(strings.TrimSpace(matchLine[0][2]), 10, 64)
			if err != nil {
				return HeaderRecordPatternError
			}

			for { // Handle continuation lines
				for _, obs := range matchLine[1:] {
					h.ObservationTypes[system] = append(h.ObservationTypes[system], strings.TrimSpace(obs[0]))
				}
				if len(h.ObservationTypes[system]) < int(totalObs) {
					line, err := ParseHeaderLine(s)
					if err != nil {
						return err
					}
					if line.Key != "SYS / # / OBS TYPES" {
						return HeaderRecordPatternError
					}
					matchLine = continuationLinePattern.FindAllStringSubmatch(line.Value, -1)
					// TODO: Check match result
				} else {
					return nil
				}
			}
		},
		"SIGNAL STRENGTH UNIT": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			h.SignalStrength = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"INTERVAL": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			h.Interval, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[:10]), 64)
			return err
		},
		"TIME OF FIRST OBS": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			h.TimeOfFirstObs, err = ParseTimeRecord(hr.Value)
			return err
		},
		"TIME OF LAST OBS": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			h.TimeOfLastObs, err = ParseTimeRecord(hr.Value)
			return err
		},
		"RCV CLOCK OFFS APPL": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / DCBS APPLIED": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / PCVS APPLIED": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / SCALE FACTOR": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / PHASE SHIFT": func(s *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err
		},
		"GLONASS SLOT / FRQ #": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"GLONASS COD/PHS/BIS": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			if strings.TrimSpace(hr.Value) == "" {
				return nil // Spec states line must be defined, but can be blank
			}
			linePattern := regexp.MustCompile(`(( [A-Z0-9]{3}) ([ 0-9]{3}[0-9].[0-9]{3}))`)
			match := linePattern.FindAllStringSubmatch(hr.Value, -1)
			if len(match) != 4 {
				return HeaderRecordPatternError
			}

			for _, cpbMatch := range match {
				code := strings.TrimSpace(cpbMatch[2])
				correction, err := strconv.ParseFloat(strings.TrimSpace(cpbMatch[3]), 64)
				if err != nil {
					return err
				}
				h.GLONASSCodePhaseBias[code] = correction
			}

			return nil
		},
		"LEAP SECONDS": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"# OF SATELLITES": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"PRN / # OF OBS": func(_ *Scanner, h *Header, hr HeaderRecord) (err error) {
			return err // TODO:
		},
		"END OF HEADER": func(s *Scanner, h *Header, hr HeaderRecord) error {
			return nil
		},
	}
)

func ParseTimeRecord(line string) (t Time, err error) {
	t.Year, err = strconv.ParseInt(strings.TrimSpace(line[:6]), 10, 64)
	if err != nil {
		return t, err
	}
	t.Month, err = strconv.ParseInt(strings.TrimSpace(line[6:12]), 10, 64)
	if err != nil {
		return t, err
	}
	t.Day, err = strconv.ParseInt(strings.TrimSpace(line[12:18]), 10, 64)
	if err != nil {
		return t, err
	}
	t.Hour, err = strconv.ParseInt(strings.TrimSpace(line[18:24]), 10, 64)
	if err != nil {
		return t, err
	}
	t.Minute, err = strconv.ParseInt(strings.TrimSpace(line[24:30]), 10, 64)
	if err != nil {
		return t, err
	}
	t.Second, err = strconv.ParseFloat(strings.TrimSpace(line[30:43]), 64)
	if err != nil {
		return t, err
	}
	t.System = strings.TrimSpace(line[48:51])
	return t, nil
}

func ParseHeaderLine(scanner *Scanner) (hr HeaderRecord, err error) {
	line, err := scanner.ReadLine()
	if err != nil {
		return hr, err
	}

	line = strings.TrimRight(line, " \n")
	if len(line) < 61 || len(line) > 80 {
		return hr, errors.New(fmt.Sprintf("invalid header line \"%s\"", line))
	}

	return HeaderRecord{line[:60], line[60:], scanner.line}, err
}

func ParseHeaderRecord(scanner *Scanner, header *Header) (hr HeaderRecord, err error) {
	hr, err = ParseHeaderLine(scanner)
	if err != nil {
		return hr, err
	}

	if parser, ok := HeaderRecordParsers[hr.Key]; ok {
		err = parser(scanner, header, hr)
		if err == nil {
			return hr, nil
		}
	} else {
		err = errors.New(fmt.Sprintf("invalid header label \"%s\"", hr.Key))
	}

	return hr, NewHeaderRecordParsingError(err, scanner.line)
}

type HeaderRecordParsingError error

func NewHeaderRecordParsingError(err error, lineNumber int) HeaderRecordParsingError {
	return errors.New(fmt.Sprintf("failed to parse header record at line %d with reason: %s", lineNumber, err.Error()))
}
