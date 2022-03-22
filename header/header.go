package header

// Free ordering of Header section, with Exceptions:
// RINEX VERSION / TYPE record MUST always be the first record in a file
// SYS / # / OBS TYPES record(s) should precede any SYS / DCBS APPLIED and SYS / SCALE FACTOR records
// # OF SATELLITES record (if present) should be immediately followed by the corresponding number of PRN / # OF OBS records.
// 		These records may be handy for documentary purposes. However, since they may only be created after having read the whole raw data file, we define them to be optional
// END OF HEADER record MUST be the last record in the header
type Header struct {
	FormatVersion   float64
	FileType        string
	SatelliteSystem string
	Program         string
	RunBy           string
	CreationDate    string // TODO: time.Time
	Comments        []string
}

func (h Header) GetFormatVersion() float64 {
	return h.FormatVersion
}

func (h Header) GetFileType() string {
	return h.FileType
}
