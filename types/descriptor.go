package types

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/solovev/gopsd/util"
)

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

func NewDescriptor(reader *util.Reader) *Descriptor {
	value := new(Descriptor)

	value.Name = reader.ReadUnicodeString()
	value.Class = reader.ReadDynamicString()
	value.Items = newDescriptorList(value, reader)

	return value
}

func newDescriptorUnitFloat(reader *util.Reader) *DescriptorUnitFloat {
	unit := new(DescriptorUnitFloat)
	unit.Type = reader.ReadString(4)
	unit.Value = reader.ReadFloat64()
	return unit
}

func newDescriptorClass(reader *util.Reader) *DescriptorClass {
	class := new(DescriptorClass)
	class.Name = reader.ReadUnicodeString()
	class.Class = reader.ReadDynamicString()
	return class
}

func newDescriptorEnum(reader *util.Reader) *DescriptorEnum {
	enum := new(DescriptorEnum)
	enum.Type = reader.ReadDynamicString()
	enum.Enum = reader.ReadDynamicString()
	return enum
}

func newDescriptorList(descriptor *Descriptor, reader *util.Reader) map[string]*DescriptorEntity {
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
		case "tdta":
			r := util.NewReader(reader.ReadBytes(reader.ReadInt32()))
			entity.Value = readTextData(r)
		default:
			panic(fmt.Sprintf("Unknown OSType key [%s] in entity [%s]", entity.Type, entity.Key))
		}
		value[entity.Key] = entity
	}
	return value
}

func readTextData(r *util.Reader) interface{} {
	r.SkipWhitespaces()
	c := r.ReadByte()
	switch c {
	case 60: // "<" - Starting of map
		r.Skip(1)
		collection := make(map[string]interface{})
		for {
			r.SkipWhitespaces()
			switch r.ReadByte() {
			case 47: // "/"
				var name []byte
				for {
					char := r.ReadByte()
					// If byte if letter (a-zA-Z)
					if (char >= 65 && char <= 90) || (char >= 97 && char <= 122) {
						name = append(name, char)
					} else {
						r.UnreadByte()
						break
					}
				}
				collection[string(name)] = readTextData(r)
			case 62: // ">"
				r.Skip(1)
				return collection
			}
		}
	case 40: // "(" - Starting of utf16 string
		r.Skip(2) // 254 & 255
		var buffer []byte
		for {
			b := r.ReadByte()
			n := len(buffer)
			if n > 0 && buffer[n-1] == 0 && b == 13 {
				buffer = buffer[0 : n-1]
				continue
			}
			// Break if byte is ")" and length of buffer is 0 or previous byte isn't "\"
			if b == 41 {
				if n == 0 || buffer[n-1] != 92 {
					break
				} else {
					buffer = append(buffer, 0)
				}
			}
			buffer = append(buffer, b)
		}
		return util.BytesToUTF16(buffer, binary.BigEndian)
	case 91: // - Starting of array
		var list []interface{}
		for {
			c = r.ReadByte()
			if c == 9 || c == 10 || c == 32 {
				continue
			}
			if c == 93 {
				return list
			}
			r.UnreadByte()
			list = append(list, readTextData(r))
		}
	default: // Primitive object (boolean/float)
		switch c {
		case 116: // "t"
			r.Skip(3) // Skip "rue"
			return true
		case 102: // "f"
			r.Skip(4) // Skip "alse"
			return false
		default:
			array := []byte{c}
			for {
				c = r.ReadByte()
				// if byte is "." or diggit
				if c == 46 || (c >= 48 && c <= 57) {
					array = append(array, c)
				} else {
					r.UnreadByte()
					break
				}
			}
			result, err := strconv.ParseFloat(string(array), 64)
			if err != nil {
				panic(err)
			}
			return result
		}
	}
}

func newDescriptorReference(reader *util.Reader) map[string]*DescriptorEntity {
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

func newDescriptorProperty(reader *util.Reader) *DescriptorProperty {
	property := new(DescriptorProperty)
	property.Name = reader.ReadUnicodeString()
	property.Class = reader.ReadDynamicString()
	property.Key = reader.ReadDynamicString()
	return property
}

func newDescriptorReferenceEnum(reader *util.Reader) *DescriptorReferenceEnum {
	enum := new(DescriptorReferenceEnum)
	enum.Name = reader.ReadUnicodeString()
	enum.Class = reader.ReadDynamicString()
	enum.Type = reader.ReadDynamicString()
	enum.Enum = reader.ReadDynamicString()
	return enum
}

func newDescriptorOffset(reader *util.Reader) *DescriptorOffset {
	offset := new(DescriptorOffset)
	offset.Name = reader.ReadUnicodeString()
	offset.Class = reader.ReadDynamicString()
	offset.Value = reader.ReadInt32()
	return offset
}

func (d *Descriptor) GetValue(path string) (interface{}, error) {
	return getValue(path, "Root", d.Items)
}

func getValue(path, collectionName string, collection map[string]*DescriptorEntity) (interface{}, error) {
	pathSplit := strings.Split(path, "->")
	pathSlice := strings.TrimSpace(pathSplit[0])
	pathIndex := -1

	if strings.HasPrefix(pathSlice, "#") { // CHECK: Or use [n] instead of #n?
		itemIndex, err := strconv.Atoi(pathSlice[1:])
		if err == nil {
			if itemIndex >= len(collection) || itemIndex < 0 {
				return nil, fmt.Errorf("Index %s out of %s's bounds", pathSlice, collectionName)
			}
			pathIndex = itemIndex
		}
	}
	i := 0
	for key, item := range collection {
		if i == pathIndex || key == pathSlice {
			if item.Type == "tdta" {
				if len(pathSplit) > 1 {
					return getTextDataValue(strings.Join(pathSplit[1:], "->"), item.Key, item.Value)
				}
				return item.Value, nil
			}
			switch instance := item.Value.(type) {
			case map[string]*DescriptorEntity:
				if len(pathSplit) > 1 {
					return getValue(strings.Join(pathSplit[1:], "->"), item.Key, instance)
				}
				return stringList(instance, 0), nil
			case *Descriptor:
				if len(pathSplit) > 1 {
					return getValue(strings.Join(pathSplit[1:], "->"), item.Key, instance.Items)
				}
				return instance.string(0), nil
			case float64, int32, bool, string:
				return instance, nil
			case *DescriptorUnitFloat:
				return instance.Value, nil
			default:
				return nil, fmt.Errorf("Can't get value of \"%s\". Unsupported type \"%s\"", pathSlice, item.Type)
			}
		}
		i++
	}
	if i == len(collection) {
		return nil, fmt.Errorf("Item \"%s\" does not exist in \"%s\"", pathSlice, collectionName)
	}
	return stringList(collection, 0), nil
}

func getTextDataValue(path, collectionName string, collection interface{}) (interface{}, error) {
	pathSplit := strings.Split(path, "->")
	pathSlice := strings.TrimSpace(pathSplit[0])
	pathIndex := -1

	if strings.HasPrefix(pathSlice, "#") {
		itemIndex, err := strconv.Atoi(pathSlice[1:])
		if err == nil {
			pathIndex = itemIndex
		}
	}

	switch instance := collection.(type) {
	case map[string]interface{}:
		if value, exist := instance[pathSlice]; exist {
			if len(pathSplit) > 1 {
				return getTextDataValue(strings.Join(pathSplit[1:], "->"), pathSlice, value)
			}
			return value, nil
		}
		return nil, fmt.Errorf("Item \"%s\" does not exist in \"%s\"", pathSlice, collectionName)
	case []interface{}:
		if pathIndex != -1 {
			if pathIndex >= len(instance) {
				return nil, fmt.Errorf("Index %s out of %s's bounds", pathSlice, collectionName)
			}
			for i, value := range instance {
				if i == pathIndex {
					if len(pathSplit) > 1 {
						return getTextDataValue(strings.Join(pathSplit[1:], "->"), pathSlice, value)
					}
					return value, nil
				}
			}
		}
		return nil, fmt.Errorf("%s is list. Specify element id, instead of \"%s\" key", collectionName, pathSplit)
	default:
		return instance, nil
	}
}

func (d *Descriptor) ToString() string {
	return d.string(0)
}

func (d Descriptor) string(indent int) string {
	sm := new(util.StringMixer)

	sm.Add("Descriptor [Class: ", d.Class, ", Length: ", fmt.Sprint(len(d.Items)), "]: ").NewLine()
	sm.AddIndent(indent).Add("{").NewLine()
	sm.Add(stringList(d.Items, indent))
	sm.AddIndent(indent).Add("}")

	return sm.String()
}

func stringList(items map[string]*DescriptorEntity, indent int) string {
	sm := new(util.StringMixer)

	for _, item := range items {
		sm.AddIndent(indent+1).Add("[", item.Type, "] \"", item.Key, "\": ")
		switch value := item.Value.(type) {
		case map[string]*DescriptorEntity: // Reference, List
			if item.Type == "obj " {
				sm.Add("Reference [Length: ", fmt.Sprint(len(items)), "]").NewLine()
			} else {
				sm.Add("List [Length: ", fmt.Sprint(len(items)), "]").NewLine()
			}
			sm.AddIndent(indent + 1).Add("{").NewLine()
			sm.Add(stringList(value, indent+2))
			sm.AddIndent(indent + 1).Add("}")
		case *Descriptor:
			sm.Add(value.string(indent + 1))
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
