package rinex3

import (
	"bufio"
	"io"
)

type RinexFile struct {
	scanner *Scanner
	Name    string // TODO: RinexFileName
	Header  RinexHeader
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
	scanner := &Scanner{bufio.NewReader(data), 0}
	header, _ := ParseHeader(scanner)
	file = RinexFile{
		scanner: scanner,
		Header:  header,
	}
	return file, err
}
