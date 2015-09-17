package util

import "fmt"

type Descriptor struct {
	Name  string
	Class string
	Items map[string]*DescriptorEntity
}

type DescriptorEntity struct {
	Key   string
	Type  string
	Value interface{}
}

type DescriptorUnitFloat struct {
	Type  string
	Value float64
}

type DescriptorClass struct {
	Name  string
	Class string
}

type DescriptorEnum struct {
	Type string
	Enum string
}

type DescriptorProperty struct {
	Name  string
	Class string
	Key   string
}

type DescriptorReferenceEnum struct {
	Name  string
	Class string
	Type  string
	Enum  string
}

type DescriptorOffset struct {
	Name  string
	Class string
	Value int32
}

func NewDescriptor(reader *Reader) *Descriptor {
	value := new(Descriptor)

	value.Name = reader.ReadUnicodeString()
	value.Class = reader.ReadDynamicString()
	value.Items = newDescriptorList(value, reader)

	return value
}

func newDescriptorUnitFloat(reader *Reader) *DescriptorUnitFloat {
	unit := new(DescriptorUnitFloat)
	unit.Type = reader.ReadString(4)
	unit.Value = reader.ReadFloat64()
	return unit
}

func newDescriptorClass(reader *Reader) *DescriptorClass {
	class := new(DescriptorClass)
	class.Name = reader.ReadUnicodeString()
	class.Class = reader.ReadDynamicString()
	return class
}

func newDescriptorEnum(reader *Reader) *DescriptorEnum {
	enum := new(DescriptorEnum)
	enum.Type = reader.ReadDynamicString()
	enum.Enum = reader.ReadDynamicString()
	return enum
}

func newDescriptorList(descriptor *Descriptor, reader *Reader) map[string]*DescriptorEntity {
	value := make(map[string]*DescriptorEntity)
	count := reader.ReadInt32()
	for i := 0; i < int(count); i++ {
		entity := new(DescriptorEntity)
		if descriptor != nil {
			entity.Key = reader.ReadDynamicString()
		}
		entity.Type = reader.ReadString(4)
		switch entity.Type {
		case "obj ":
			entity.Value = newDescriptorReference(reader)
		case "Objc", "GlbO":
			entity.Value = NewDescriptor(reader)
		case "VlLs":
			entity.Value = newDescriptorList(nil, reader)
		case "doub":
			entity.Value = reader.ReadFloat64()
		case "UntF":
			entity.Value = newDescriptorUnitFloat(reader)
		case "TEXT":
			entity.Value = reader.ReadUnicodeString()
		case "enum":
			entity.Value = newDescriptorEnum(reader)
		case "long":
			entity.Value = reader.ReadInt32()
		case "bool":
			entity.Value = reader.ReadByte() == 1
		case "type", "GlbC":
			entity.Value = newDescriptorClass(reader)
		case "alis": // TODO
			reader.Skip(reader.ReadInt32())
		case "tdta": // TODO
			reader.Skip(reader.ReadInt32())
		default:
			panic(fmt.Sprintf("Unknown OSType key [%s] in entity [%s]", entity.Type, entity.Key))
		}
		value[entity.Key] = entity
	}
	return value
}

func newDescriptorReference(reader *Reader) map[string]*DescriptorEntity {
	value := make(map[string]*DescriptorEntity)
	count := reader.ReadInt32()
	for i := 0; i < int(count); i++ {
		entity := new(DescriptorEntity)
		entity.Type = reader.ReadString(4)
		switch entity.Type {
		case "prop":
			entity.Value = newDescriptorProperty(reader)
		case "Clss":
			entity.Value = newDescriptorClass(reader)
		case "Enmr":
			entity.Value = newDescriptorReferenceEnum(reader)
		case "rele":
			entity.Value = newDescriptorOffset(reader)
		case "Idnt":
		case "indx":
		case "name":
		default:
			panic(fmt.Sprintf("Unknown OSType key [%s] in entity [%s]", entity.Type, entity.Key))
		}
		value[entity.Key] = entity
	}
	return value
}

func newDescriptorProperty(reader *Reader) *DescriptorProperty {
	property := new(DescriptorProperty)
	property.Name = reader.ReadUnicodeString()
	property.Class = reader.ReadDynamicString()
	property.Key = reader.ReadDynamicString()
	return property
}

func newDescriptorReferenceEnum(reader *Reader) *DescriptorReferenceEnum {
	enum := new(DescriptorReferenceEnum)
	enum.Name = reader.ReadUnicodeString()
	enum.Class = reader.ReadDynamicString()
	enum.Type = reader.ReadDynamicString()
	enum.Enum = reader.ReadDynamicString()
	return enum
}

func newDescriptorOffset(reader *Reader) *DescriptorOffset {
	offset := new(DescriptorOffset)
	offset.Name = reader.ReadUnicodeString()
	offset.Class = reader.ReadDynamicString()
	offset.Value = reader.ReadInt32()
	return offset
}

func (d Descriptor) String(indent int) string {
	sm := new(StringMixer)

	sm.AddIndent(indent).Add("Descriptor [", fmt.Sprint(len(d.Items)), "]: ", d.Class).NewLine()
	sm.AddIndent(indent).Add("{").NewLine()
	sm.Add(stringList(d.Items, indent))
	sm.AddIndent(indent).Add("}")

	return sm.String()
}

func stringList(items map[string]*DescriptorEntity, indent int) string {
	sm := new(StringMixer)

	for _, item := range items {
		sm.AddIndent(indent+1).Add("[", item.Type, "] ", item.Key, ": ")
		switch value := item.Value.(type) {
		case map[string]*DescriptorEntity: // Reference, List
			if item.Type == "obj " {
				sm.AddIndent(indent+2).Add("Reference [", fmt.Sprint(len(items)), "]").NewLine()
			} else {
				sm.AddIndent(indent+2).Add("List [", fmt.Sprint(len(items)), "]").NewLine()
			}
			sm.AddIndent(indent + 2).Add("{").NewLine()
			sm.Add(stringList(value, indent+2))
			sm.AddIndent(indent + 2).Add("}")
		case *Descriptor:
			sm.NewLine()
			sm.Add(value.String(indent + 2))
		case float64, int32, bool:
			sm.Add(fmt.Sprint(value))
		case *DescriptorUnitFloat:
			sm.Add("[Type: ", value.Type, ", Value: ", fmt.Sprint(value.Value), "]")
		case string:
			sm.Add(value)
		case *DescriptorEnum:
			sm.Add("[Type: ", value.Type, ", Enum: ", value.Enum, "]")
		case *DescriptorClass:
			sm.Add("[Name: ", value.Name, ", Class: ", value.Class, "]")
		case *DescriptorProperty:
			sm.Add("[Key: ", value.Key, " Name: ", value.Name, ", Class: ", value.Class, "]")
		case *DescriptorOffset:
			sm.Add("[Name: ", value.Name, ", Class: ", value.Class, " Value: ", fmt.Sprint(value.Value), "]")
		case *DescriptorReferenceEnum:
			sm.Add("[Type: ", value.Type, ", Enum: ", value.Enum, ", Class: ", value.Class, " Name: ", value.Name, "]")
		default:
			sm.Add("?")
		}
		sm.NewLine()
	}

	return sm.String()
}

type Rectangle struct {
	top    int32 `json:"-"`
	left   int32 `json:"-"`
	bottom int32 `json:"-"`
	right  int32 `json:"-"`

	X, Y, Width, Height int32
}

func NewRectangle(reader *Reader) *Rectangle {
	r := new(Rectangle)

	r.top = reader.ReadInt32()
	r.Y = r.top

	r.left = reader.ReadInt32()
	r.X = r.left

	r.bottom = reader.ReadInt32()
	r.right = reader.ReadInt32()

	r.Width = r.right - r.left
	r.Height = r.bottom - r.top

	return r
}

func (r Rectangle) ToString() string {
	return fmt.Sprintf("[X: %d, Y: %d, Width: %d, Height: %d]", r.X, r.Y, r.Width, r.Height)
}

type Color struct {
	red, green, blue, alpha int16
}

func (c Color) Red() int16 {
	return c.red
}

func (c Color) Green() int16 {
	return c.green
}

func (c Color) Blue() int16 {
	return c.blue
}

func (c Color) Alpha() int16 {
	return c.alpha
}

func NewRGBAColor(reader *Reader) *Color {
	return &Color{reader.ReadInt16(), reader.ReadInt16(), reader.ReadInt16(), reader.ReadInt16()}
}
