package rinex

import (
	"errors"
	"regexp"
)

type Filename struct {
	StationName string
	DataSource string // rune ?
	StartTime string // time.Time ?
	Duration string // time.Duration ?
	Frequency string // time.Duration ?
	FileType string
	FileFormat string // rnx, crx, 
	Compression string
}

var (
	// TODO: 
	//ShortNamePattern *regexp.Regexp = regexp.MustCompile(
	//	`^(?<siteId>[\w]{4})` +
    //    `(?<dayOfYear>\d{3})` +
    //    `((?<daily>0)|((?<hour>[a-x])(?<minute>(00|15|30|45)?)))` +
    //    `\.(?<year>\d{2})` +
    //    `(?<fileTypeCode>[abcdfghlmnopqs])` +
    //    `((\.(?<extension>((z)|(gz)|(bz2)|(zip))))?)$`)

	LongNamePattern *regexp.Regexp = regexp.MustCompile(
		`^(?P<stationName>[\w]{4}\d{2}[a-zA-Z]{3})_` +
        `(?P<dataSource>[RSU])_` +
        `(?P<startTime>\d{11})_` +
        `(?P<duration>\d{2}[SMHDYU])_` +
        `((((?P<navFrequency>\d{2}[CZSMHDU])_)?(?P<navFileTypeCode>[GREJCISM]N))` +
        `|((?P<obsFrequency>\d{2}[CZSMHDU])_(?P<obsFileTypeCode>[GREJCISM]O))` +
        `|(((?P<metFrequency>\d{2}[CZSMHDU])_)?(?P<metFileTypeCode>MM)))` +
        `\.(?P<format>((rnx)|(crx)))(\.(?P<compression>.*))?$`)
)

func NewRinexFilename(name string) (filename Filename, err error) {
	matchGroups := LongNamePattern.FindStringSubmatch(name)
	if len(matchGroups) == 0 {
		return filename, errors.New("invalid RINEX filename")
	}

	groups := map[string]string{}
	for i, name := range LongNamePattern.SubexpNames(){
		if name != "" {
			groups[name] = matchGroups[i]
		}
	}

	filename = Filename{
		StationName: groups["stationName"],
		DataSource: groups["dataSource"],
		StartTime: groups["startTime"],
		Duration: groups["duration"],
		FileFormat: groups["format"],
		Compression: groups["compression"],
	}

	switch {
		case groups["navFileTypeCode"] != "":
			filename.FileType = groups["navFileTypeCode"]
			filename.Frequency = groups["navFrequency"]
		case groups["obsFileTypeCode"] != "":
			filename.FileType = groups["obsFileTypeCode"]
			filename.Frequency = groups["obsFrequency"]
		case groups["metFileTypeCode"] != "":
			filename.FileType = groups["metFileTypeCode"]
			filename.Frequency = groups["metFrequency"]
	}

	return filename, nil
}
