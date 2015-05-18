package gopsd

import "io/ioutil"

type Document struct {
	Header    *HeaderSection
	ColorMode *ColorModeDataSection
}

var (
	reader *Reader
)

func Parse(path string) (*Document, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	reader = NewReader(data)

	doc := new(Document)

	header, err := newHeader()
	if err != nil {
		return nil, err
	}
	doc.Header = header

	cm, err := newColorMode(header.ColorMode)
	if err != nil {
		return nil, err
	}
	doc.ColorMode = cm

	return doc, nil
}
