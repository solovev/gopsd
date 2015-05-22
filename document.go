package gopsd

import (
	"errors"
	"io/ioutil"
)

type Document struct {
	IsLarge bool

	Channels  int16
	Height    int32
	Width     int32
	Depth     int16
	ColorMode string

	Resources map[int16]*ImageResource
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
	readHeader(doc)
	readColorMode(doc)
	readResources(doc)

	return doc, nil
}
