package rinex3

// TODO: NewHeader function to provide non nil values for fields like Antenna

type Header struct {
	FormatVersion float64
	FileType string
	SatelliteSystem string
	Program string
	RunBy string
	CreationDate string // TODO: time.Time
	Marker Marker
	Observer string
	Agency string
	Receiver Receiver
	Antenna Antenna
	CenterOfMass struct {
		X float64
		Y float64
		Z float64
	}
	ObservationTypes map[string][]string // TODO: map[SatelliteSystem][]ObservationType
	Comments []HeaderComment
	records []HeaderRecord // TODO: Does this need to exist?
}

func NewHeader() Header {
	return Header{
		Marker: Marker{},
		Receiver: Receiver{},
		Antenna: Antenna{},
		ObservationTypes: map[string][]string{},
	}
}

type Marker struct {
	Name string
	Number string
	Type string
	ApproxPosition struct {
		X float64
		Y float64
		Z float64
	}
}

type Receiver struct {
	Number string
	Type string
	Version string
}

type Antenna struct {
	Number string
	Type string
	Height float64
	East float64
	North float64
	// TODO: What should this be called?
	DeltaXYZ struct {
		X float64
		Y float64
		Z float64
	}
	PhaseCenter struct {
		SatelliteSystem string
		ObservationCode string
		Position struct {
			// TODO: Some triple type - spec says can be XYZ or NEU
		}
	}
	BSight struct {
		X float64
		Y float64
		Z float64
	}
	ZeroDirection struct {
		Azimuth float64
		X float64
		Y float64
		Z float64
	}
}

type HeaderComment struct {
	Comment string
	Line int
}
