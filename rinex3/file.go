package rinex3

import (
	"bufio"
	"io"
)

// TODO: Observation / Navigation / Meteorological RinexFile (currently only parsing Obs)
type RinexFile struct {
	scanner *Scanner
	Name    string // TODO: RinexFileName
	Header  Header
}

// TODO: Can scan through observations like so, or parse them into a map
func (r RinexFile) NextEpoch() (epoch EpochRecord, err error) {
	return epoch, err
}

type Scanner struct {
	*bufio.Reader
	line int
}

func (s *Scanner) ReadLine() (line string, err error) {
	s.line += 1
	return s.ReadString('\n')
}

func ParseRinexFile(data io.Reader) (file RinexFile, err error) {
	file = RinexFile{
		scanner: &Scanner{bufio.NewReader(data), 0},
		Header:  NewHeader(),
	}
	err = ParseHeader(file.scanner, &file.Header)
	return file, err
}
