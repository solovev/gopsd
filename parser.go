package gopsd

import "io/ioutil"

type Document struct {
	Header *HeaderSection
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

	return doc, nil
}
