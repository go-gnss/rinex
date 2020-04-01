package rinex3

import (
	"bufio"
	"fmt"
	"github.com/alexflint/go-restructure"
)

// Multiple epoch observation data records with identical time tags are not allowed (exception: Event records).
// Epochs MUST appear ordered in time.

type ObservationEpochRecord struct {
	Epoch              Epoch
	ObservationRecords []ObservationRecord
}

type Epoch struct {
	_             struct{} `^> `
	Year          int      `\d{4} `
	Month         int      `\d{2} `
	Day           int      `\d{2} `
	Hour          int      `\d{2} `
	Minute        int      `\d{2} `
	Second        float64  `.{11}  `
	EpochFlag     bool     `[01] `
	NumSatellites int      `...`
	ClockOffset   float64  ` {6}?.{15}?$`
}

type ObservationRecord struct {
	SatelliteNumber string `^[a-zA-z][ 0-9][0-9]`
	Observations    []Observation
	_               struct{} `$`
}

type Observation struct {
	Value          float64 `.{14}`
	LLI            bool    `[01 ]`
	SignalStrength bool    `[01]`
}

func NextObservationEpochRecord(reader *bufio.Reader) (obs ObservationEpochRecord, err error) {
	line, err := reader.ReadString('\n')
	fmt.Println(line, err)
	fmt.Println(restructure.Find(&obs, line))
	line, err = reader.ReadString('\n')
	fmt.Println(line, err)
	return obs, err
}
