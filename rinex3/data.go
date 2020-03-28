package rinex3

import (
    "fmt"
    "github.com/alecthomas/participle"
//    "github.com/alecthomas/participle/lexer"
    "io"
)

// Probably need a custom lexer github.com/alecthomas/participle/lexer/text_scanner.go
// Might need custom scanner too text/scanner

type ObservationData struct {
    ObservationDataEpochs []ObservationDataEpoch `@@*`
}

type ObservationDataEpoch struct {
    Epoch Epoch `">" @@`
    ObservationRecords []ObservationRecord `@@*`
}

type Epoch struct {
    Year int `@Int`
    Month int `@Int`
    Day int `@Int`
    Hour int `@Int`
    Minute int `@Int` // These are padded with a 0 which scanner.Int interprets to mean the value is octal, making it fail for 08 and 09
    Second float64 `@Float`
    EpochFlag int `@Int`
    NumSatellites int `@Int`
    ClockOffset float64 `@("-"? Float)?`
}

type ObservationRecord struct {
    SatelliteInt string `@Ident`
    Observations []Observation `@@*`
}

type Observation struct {
    Value float64 `@("-"? Float)`
    LliSignalStrength int `@Int?` // There won't be space between Value, LLI, and Signal Strength if all three are present
}

// TODO: This will move up to ParseFile or something
func Parse(data io.Reader) (obs *ObservationData) {
    parser, err := participle.Build(&ObservationData{})

    //r := &ObservationData{}
    err = parser.Parse(data, obs)
    fmt.Println(obs, err)
	return obs
}
