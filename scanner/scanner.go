package scanner

import "bufio"

type Scanner struct {
	*bufio.Reader
	Line int
}

func (s *Scanner) ReadLine() (line string, err error) {
	s.Line += 1
	return s.ReadString('\n')
}
