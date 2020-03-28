package rinex3

import (
	"io"
    "github.com/alecthomas/participle"
)

type RinexFile struct {
	Header Header
	Observation ObservationDataRecord `@@*`
}

// TODO: This will move up to ParseFile or something
func Parse(data io.Reader) (file *RinexFile, err error) {
    parser, err := participle.Build(&RinexFile{})
	if err != nil {
		return file, err
	}

    file = &RinexFile{}
    err = parser.Parse(data, file)
	return file, err
}
