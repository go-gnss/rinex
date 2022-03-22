package rinex3

import (
	"bufio"

	"github.com/go-gnss/rinex/header"
)

type ObservationHeader struct {
	header.Header
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
}

func NewObservationHeader(header header.Header) ObservationHeader {
	return ObservationHeader{
		Header:               header,
		ObservationTypes:     map[string][]string{},
		PhaseShifts:          map[string][]float64{},
		GLONASSCodePhaseBias: map[string]float64{},
	}
}

func ParseObservationHeader(r *bufio.Reader, header *ObservationHeader) (err error) {
	hr, err := ParseObservationHeaderRecord(r, header)
	for err != nil || hr.Key != "END OF HEADER" {
		hr, err = ParseObservationHeaderRecord(r, header)
	}
	return err
}
