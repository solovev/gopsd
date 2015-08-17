package gopsd

import (
	"errors"
	"image"
	"io/ioutil"

	"github.com/solovev/gopsd/util"
)

// TODO all INT -> INT64 (**PSB**)
// TODO make([]interface{}, 0) -> var name []interface{}
type Document struct {
	IsLarge bool

	Channels  int16
	Height    int32
	Width     int32
	Depth     int16
	ColorMode string
	Image     image.Image

	Resources map[int16]interface{}
	Layers    []*Layer
}

var (
	reader *util.Reader
)

func (Document) ToJson() ([]byte, error) {
	return nil, nil
}

func ParseFromBuffer(buffer []byte) (doc *Document, err error) {
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

	reader = util.NewReader(buffer)

	doc = new(Document)
	readHeader(doc)
	readColorMode(doc)
	readResources(doc)
	readLayers(doc)
	readImageData(doc)

	return doc, nil
}

func ParseFromPath(path string) (doc *Document, err error) {
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
	reader = util.NewReader(data)

	doc = new(Document)
	readHeader(doc)
	readColorMode(doc)
	readResources(doc)
	readLayers(doc)
	readImageData(doc)

	return doc, nil
}
