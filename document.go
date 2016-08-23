package gopsd

import (
	"encoding/json"
	"errors"
	"image"
	"io/ioutil"

	"github.com/solovev/gopsd/types"
	"github.com/solovev/gopsd/util"
)

// TODO all INT -> INT64 (**PSB**)
// TODO make([]interface{}, 0) -> var name []interface{}
// TODO Set Flags bits (ints now) to separated booleans
// TODO Remove panic, add error to return
// Replace New*** to Read*** if reader object passed as parameter
//	- Example:
//		Flags. bit 1 = invert, bit 2 = not link, bit 3 = disable
//		From:
//			Flags int
//		To:
//			IsInverted, IsNotLinked, IsDisabled bool
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

func (d *Document) GetLayersByName(name string) []*Layer {
	var layers []*Layer
	for _, layer := range d.Layers {
		if layer.Name == name {
			layers = append(layers, layer)
		}
	}
	return layers
}

func (d *Document) GetLayerByID(id int) *Layer {
	for _, layer := range d.Layers {
		if layer.ID == int32(id) {
			return layer
		}
	}
	return nil
}

func (d *Document) GetLayer(index int) *Layer {
	if index >= len(d.Layers) {
		return nil
	}
	return d.Layers[index]
}

func (d *Document) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Document) GetTreeRepresentation() *Layer {
	root := new(Layer)
	root.ID = -1
	root.Name = "RootLayer"
	root.Rectangle = types.CreateRectangle(0, 0, d.Width, d.Height)

	current := root
	for i := len(d.Layers) - 1; i >= 0; i-- {
		entry := d.Layers[i]
		if entry.IsFolder {
			entry.Parent = current
			current.Children = append(current.Children, entry)
			current = entry
		} else if entry.IsSectionDivider {
			current = current.Parent
		} else {
			entry.Parent = current
			current.Children = append(current.Children, entry)
		}
	}
	return root
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

func ParseFromPath(path string) (*Document, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	doc, err := ParseFromBuffer(data)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
