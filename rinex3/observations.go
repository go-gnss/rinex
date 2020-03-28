package rinex3

// HeaderLines can appear anywhere in the file according to spec
// However, (though the spec is not specific on this) only the
// Header has an END OF HEADER line

type ObservationDataRecord struct {
	ObservationEpochRecord `@@`
	HeaderLine string `| @@`
}

type ObservationEpochRecord struct {
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
