package rinex

// Multiple epoch observation data records with identical time tags are not allowed (exception: Event records).
// Epochs MUST appear ordered in time.

type EpochRecord struct {
	Epoch              Epoch
	ObservationRecords []ObservationRecord
}

type Epoch struct {
	Year          int
	Month         int
	Day           int
	Hour          int
	Minute        int
	Second        float64
	EpochFlag     bool
	NumSatellites int
	ClockOffset   float64
}

type ObservationRecord struct {
	SatelliteNumber string
	Observations    []Observation
}

type Observation struct {
	Value          float64
	LLI            bool
	SignalStrength bool
}

func ParseEpochRecord(s *Scanner) (epoch EpochRecord, err error) {
	// TODO:
	return epoch, err
}
