package rinex3

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type ObservationFile struct {
	ObservationHeader
	Epochs []EpochRecord
}

func ParseEpochRecords(r *bufio.Reader, observationTypes map[string][]string) (records []EpochRecord, err error) {
	for {
		record, err := ParseEpochRecord(r, observationTypes)
		if err == io.EOF {
			return records, nil
		} else if err != nil {
			return records, err
		}
		records = append(records, record)
	}
}

// Multiple epoch observation data records with identical time tags are not allowed (exception: Event records).
// Epochs MUST appear ordered in time.

type EpochRecord struct {
	Time               time.Time
	Flag               int
	NumSatellites      int
	ClockOffset        float64
	ObservationRecords []ObservationRecord
}

type ObservationRecord struct {
	Constellation   string
	SatelliteNumber int
	Observations    []Observation
}

type Observation struct {
	Value          float64
	LLI            int
	SignalStrength int
}

func ParseEpochRecord(r *bufio.Reader, observationTypes map[string][]string) (epoch EpochRecord, err error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return epoch, err
	}

	if string(line[0]) != ">" {
		return epoch, fmt.Errorf("invalid epoch record: %s", line)
	}

	t, err := time.Parse("2006 01 02 15 04", line[2:18])
	if err != nil {
		return epoch, err
	}
	seconds, err := strconv.ParseInt(strings.TrimSpace(line[19:21]), 10, 8)
	if err != nil {
		return epoch, err
	}
	milliseconds, err := strconv.ParseInt(strings.TrimSpace(line[22:29]), 10, 64)
	if err != nil {
		return epoch, err
	}
	t = t.Add(time.Duration(seconds) * time.Second)
	t = t.Add(time.Duration(milliseconds) * time.Millisecond)
	epoch.Time = t

	flag, err := strconv.ParseInt(line[31:32], 10, 8)
	if err != nil {
		return epoch, err
	}
	epoch.Flag = int(flag)

	numSats, err := strconv.ParseInt(strings.TrimSpace(line[32:35]), 10, 8)
	if err != nil {
		return epoch, err
	}
	epoch.NumSatellites = int(numSats)

	offset := strings.TrimSpace(line[35 : len(line)-1])
	if offset != "" {
		epoch.ClockOffset, err = strconv.ParseFloat(offset, 64)
		if err != nil {
			return epoch, err
		}
	}

	// Parse each ObservationRecord within EpochRecord
	for i := 0; i < int(numSats); i++ {
		line, err = r.ReadString('\n')
		if err != nil {
			return epoch, err
		}
		line = line[:len(line)-1] + "  " // Cheating because for some reason the fixture data doesn't have space for LLI or Signal strength for the last record (not optional in spec...)

		record, err := ParseObservationRecord(line, observationTypes)
		if err != nil {
			return epoch, err
		}

		epoch.ObservationRecords = append(epoch.ObservationRecords, record)
	}

	return epoch, err
}

func ParseObservationRecord(line string, observationTypes map[string][]string) (record ObservationRecord, err error) {
	sat, err := strconv.ParseInt(strings.TrimSpace(line[1:3]), 10, 64)
	if err != nil {
		return record, err
	}

	record.Constellation = line[0:1]
	record.SatelliteNumber = int(sat)

	for i := 0; i < len(observationTypes[record.Constellation]); i++ {
		// account for line ending early if not all signals present for satellite
		if len(line) < (19 + (16 * i)) {
			break
		}

		observation, err := ParseObservation(line[3+(16*i) : 19+(16*i)])
		if err != nil {
			return record, err
		}
		record.Observations = append(record.Observations, observation)
	}

	return record, nil
}

func ParseObservation(data string) (obs Observation, err error) {
	if data[:14] != "              " {
		obs.Value, err = strconv.ParseFloat(strings.TrimSpace(data[:14]), 64)
		if err != nil {
			return obs, err
		}
	}

	if data[14:15] != " " {
		lli, err := strconv.ParseInt(data[14:15], 10, 8)
		if err != nil {
			return obs, err
		}
		obs.LLI = int(lli)
	}

	if data[14:15] != " " {
		strength, err := strconv.ParseInt(data[15:16], 10, 8)
		if err != nil {
			return obs, err
		}
		obs.SignalStrength = int(strength)
	}

	return obs, nil
}
