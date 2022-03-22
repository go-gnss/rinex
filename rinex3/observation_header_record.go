package rinex3

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-gnss/rinex/header"
)

type ObservationHeaderRecordParser func(*bufio.Reader, *ObservationHeader, header.HeaderRecord) error

var (
	HeaderRecordPatternError error = errors.New("header record did not match pattern")
	// TODO: Consider moving function definitions out of the map
	// TODO: Consider using a proper parser library again, or just regex groups w/
	// Getters and Setters on Header
	ObservationHeaderRecordParsers map[string]ObservationHeaderRecordParser = map[string]ObservationHeaderRecordParser{
		"MARKER NAME": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) error {
			h.Marker.Name = strings.TrimSpace(hr.Value)
			return nil
		},
		"MARKER NUMBER": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) error {
			h.Marker.Number = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"MARKER TYPE": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) error {
			h.Marker.Type = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"OBSERVER / AGENCY": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) error {
			h.Observer = strings.TrimSpace(hr.Value[:20])
			h.Agency = strings.TrimSpace(hr.Value[20:])
			return nil
		},
		"REC # / TYPE / VERS": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) error {
			h.Receiver.Number = strings.TrimSpace(hr.Value[:20])
			h.Receiver.Type = strings.TrimSpace(hr.Value[20:40])
			h.Receiver.Version = strings.TrimSpace(hr.Value[40:])
			return nil
		},
		"ANT # / TYPE": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) error {
			h.Antenna.Number = strings.TrimSpace(hr.Value[:20])
			h.Antenna.Type = strings.TrimSpace(hr.Value[20:40])
			return nil
		},
		"APPROX POSITION XYZ": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
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
		"ANTENNA: DELTA H/E/N": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
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
		"ANTENNA: DELTA X/Y/Z": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err
		},
		"ANTENNA: PHASECENTER": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"ANTENNA: B.SIGHT XYZ": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"ANTENNA: ZERODIR AZI": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"ANTENNA: ZERODIR XYZ": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"CENTER OF MASS: XYZ": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / # / OBS TYPES": func(s *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
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
					line, err := header.ParseHeaderRecord(s)
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
		"SIGNAL STRENGTH UNIT": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			h.SignalStrength = strings.TrimSpace(hr.Value[:20])
			return nil
		},
		"INTERVAL": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			h.Interval, err = strconv.ParseFloat(strings.TrimSpace(hr.Value[:10]), 64)
			return err
		},
		"TIME OF FIRST OBS": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			h.TimeOfFirstObs, err = ParseTimeRecord(hr.Value)
			return err
		},
		"TIME OF LAST OBS": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			h.TimeOfLastObs, err = ParseTimeRecord(hr.Value)
			return err
		},
		"RCV CLOCK OFFS APPL": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / DCBS APPLIED": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / PCVS APPLIED": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / SCALE FACTOR": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"SYS / PHASE SHIFT": func(s *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err
		},
		"GLONASS SLOT / FRQ #": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"GLONASS COD/PHS/BIS": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
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
		"LEAP SECONDS": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"# OF SATELLITES": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
		"PRN / # OF OBS": func(_ *bufio.Reader, h *ObservationHeader, hr header.HeaderRecord) (err error) {
			return err // TODO:
		},
	}
)

type Time struct { // TODO: time.Time
	Year   int64
	Month  int64
	Day    int64
	Hour   int64
	Minute int64
	Second float64
	System string
}

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

func ParseObservationHeaderRecord(r *bufio.Reader, obsHeader *ObservationHeader) (hr header.HeaderRecord, err error) {
	hr, err = header.ParseHeaderRecord(r)
	if err != nil {
		return hr, err
	}

	if parser, ok := header.HeaderRecordParsers[hr.Key]; ok {
		err = parser(r, &obsHeader.Header, hr)
		if err == nil {
			return hr, nil
		}
	}

	if parser, ok := ObservationHeaderRecordParsers[hr.Key]; ok {
		err = parser(r, obsHeader, hr)
		if err == nil {
			return hr, nil
		}
	} else {
		err = errors.New(fmt.Sprintf("invalid header label \"%s\"", hr.Key))
	}

	return hr, fmt.Errorf("failed to parse header record: %e", err)
}
