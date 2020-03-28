package rinex3

type Header struct {
	RinexVersion string
	RinexType string
	Program string
	RunBy string
	CreationDate string // yyyymmddhhmmss zone
	MarkerName string
	MarkerNumber string
	MarkerType string
	Observer string
	Agency string
	ReceiverNumber string
	ReceiverType string
	ReceiverVersion string
}
