package gopsd

import (
	"encoding/json"
	"errors"
	"image"
	"io/ioutil"

	"github.com/solovev/gopsd/util"
)

// TODO all INT -> INT64 (**PSB**)
// TODO make([]interface{}, 0) -> var name []interface{}
type Document struct {
	IsLarge bool `json:"-"`

	Channels  int16 `json:"-"`
	Height    int32
	Width     int32
	Depth     int16       `json:"-"`
	ColorMode string      `json:"-"`
	Image     image.Image `json:"-"`

	Resources map[int16]interface{} `json:"-"`
	Layers    []*Layer
}

var (
	reader *util.Reader
)

func (d *Document) ToJSON() ([]byte, error) {
	return json.Marshal(d)
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
