package rinex3

import (
	"bufio"
)

type RinexFile struct {
	Header Header
	Epochs []ObservationEpochRecord
}

func Parse(data bufio.Reader) (file *RinexFile, err error) {
	return file, err
}
