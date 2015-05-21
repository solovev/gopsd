package gopsd

import (
	"errors"
	"io/ioutil"
)

type Document struct {
	Header *HeaderSection
}

var (
	reader *Reader
)

func Parse(path string) (doc *Document, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch value := r.(type) {
			case string:
				err = errors.New(value)
			case error:
				err = value
			}
			doc = nil
		}
	}()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reader = NewReader(data)

	doc = new(Document)
	doc.Header = newHeader()

	return doc, nil
}
