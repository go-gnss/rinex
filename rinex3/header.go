package rinex3

type HeaderComment struct {
	Comment string
	Line    int // TODO: This might not be useful for reconstructing Headers if additional lines are added
}

// Free ordering of Header section, with Exceptions:
// RINEX VERSION / TYPE record MUST always be the first record in a file
// SYS / # / OBS TYPES record(s) should precede any SYS / DCBS APPLIED and SYS / SCALE FACTOR records
// # OF SATELLITES record (if present) should be immediately followed by the corresponding number of PRN / # OF OBS records.
// 		These records may be handy for documentary purposes. However, since they may only be created after having read the whole raw data file, we define them to be optional
// END OF HEADER record MUST be the last record in the header

// TODO: Obs / Nav / Met have different headers, but all share "RINEX
// VERSION / TYPE", "PGM / RUN BY / DATE", and "COMMENT" HeaderRecords
// - Obs and Met also share "MARKER NAME" and "MARKER NUMBER" - could
// implement as Interface which can be cast to specific type
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
	// TODO: Probably don't want to define any of these structs inline
	Marker struct {
		Name           string
		Number         string
		Type           string
		ApproxPosition struct {
			X float64
			Y float64
			Z float64
		}
	}
	Observer string
	Agency   string
	Receiver struct {
		Number  string
		Type    string
		Version string
	}
	Antenna struct {
		Number string
		Type   string
		Height float64
		East   float64
		North  float64
		// TODO: Figure out how to deal with body-fixed vs fixed station
		//X      float64
		//Y      float64
		//Z      float64
		//PhaseCenter map[string]map[string]struct {
		//	// TODO: Some triple type - spec says can be XYZ (body-fixed) or NEU (fixed station)
		//}
		//BSight struct {
		//	X float64
		//	Y float64
		//	Z float64
		//}
		//ZeroDirection struct {
		//	Azimuth float64
		//	X       float64
		//	Y       float64
		//	Z       float64
		//}
	}
	//CenterOfMass struct {
	//	X float64
	//	Y float64
	//	Z float64
	//}
	ObservationTypes     map[string][]string // TODO: map[SatelliteSystem][]ObservationType
	SignalStrength       string
	Interval             float64
	TimeOfFirstObs       Time
	TimeOfLastObs        Time
	PhaseShifts          map[string][]float64
	GLONASSCodePhaseBias map[string]float64 // TODO: map[Signal]float64
	Comments             []HeaderComment
}

type Time struct { // TODO: time.Time
	Year   int64
	Month  int64
	Day    int64
	Hour   int64
	Minute int64
	Second float64
	System string
}

func NewHeader() Header {
	return Header{
		ObservationTypes:     map[string][]string{},
		PhaseShifts:          map[string][]float64{},
		GLONASSCodePhaseBias: map[string]float64{},
	}
}

// TODO: Check for empty strings / missing required values?
func ParseHeader(scanner *Scanner, header *Header) (err error) {
	// TODO: Check if first line parsed is RINEX VERSION / TYPE
	hr, err := ParseHeaderRecord(scanner, header)
	for ; err == nil && hr.Key != "END OF HEADER"; hr, err = ParseHeaderRecord(scanner, header) {
	}
	return err
}
