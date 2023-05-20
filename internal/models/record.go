package models

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

var UnsupportedType error = errors.New("unsupported type")

type SerializableRecord interface {
	String() (string, error)
	Bytes() ([]byte, error)
}

type ComparableByField interface {
	LessThan(other ComparableByField, field string) (bool, error)
}

type Record []Field

func recordFromBytes(b []byte, rd RecordDescriptor) (*Record, error) {
	result := &Record{}
	buf := bytes.NewBuffer(b)
	for _, fd := range rd {
		f := Field{
			Descriptor: fd,
		}
		tmp := make([]byte, fd.Size)
		_, err := buf.Read(tmp)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize field: %w", err)
		}
		switch fd.Type {
		case "string":
			f.Value = string(tmp)
		case "uint32":
			f.Value = binary.LittleEndian.Uint32(tmp)
		default:
			return nil, UnsupportedType
		}
		*result = append(*result, f)
	}
	return result, nil
}

func SerializableRecordFromBytes(b []byte, rd RecordDescriptor) (SerializableRecord, error) {
	return recordFromBytes(b, rd)
}

func ComparableByFieldFromBytes(b []byte, rd RecordDescriptor) (ComparableByField, error) {
	return recordFromBytes(b, rd)
}

func (r *Record) Size() int {
	result := 0
	for _, fd := range *r {
		result += fd.Descriptor.Size
	}
	return result
}

func (r *Record) LessThan(other ComparableByField, field string) (bool, error) {
	for _, f := range *r {
		if f.Descriptor.Name != field {
			continue
		}
		switch f.Descriptor.Type {
		case "string":
			oval, err := other.(*Record).Value(f.Descriptor.Name)
			if err != nil {
				return false, err
			}
			return f.Value.(string) < oval.(string), nil
		case "uint32":
			oval, err := other.(*Record).Value(f.Descriptor.Name)
			if err != nil {
				return false, err
			}
			return f.Value.(uint32) < oval.(uint32), nil
		default:
			return false, UnsupportedType
		}
	}
	return false, fmt.Errorf("unknown field %[1]s", field)
}

func (r *Record) String() (string, error) {
	result := strings.Builder{}
	for idx, f := range *r {
		if idx > 0 {
			err := result.WriteByte(' ')
			if err != nil {
				return "", err
			}
		}
		_, err := result.WriteString(fmt.Sprintf("%v", f.Value))
		if err != nil {
			return "", err
		}
		if f.Descriptor.Type == "string" {
			strlen := len(f.Value.(string))
			for i := strlen; i < f.Descriptor.Size; i++ {
				err := result.WriteByte(' ')
				if err != nil {
					return "", err
				}
			}
		}
	}
	return result.String(), nil
}

func (r *Record) Bytes() ([]byte, error) {
	b := make([]byte, 0, r.Size())
	buf := bytes.NewBuffer(b)
	for _, f := range *r {
		switch f.Descriptor.Type {
		case "string":
			val := []byte(f.Value.(string))
			strlen := len(val)
			if f.Descriptor.Size < strlen {
				strlen = f.Descriptor.Size
			}
			_, err := buf.Write(val[:strlen])
			if err != nil {
				return nil, err
			}
			for i := strlen; i < f.Descriptor.Size; i++ {
				err := buf.WriteByte(' ')
				if err != nil {
					return nil, err
				}
			}
		case "uint32":
			intbuf := make([]byte, 4)
			binary.LittleEndian.PutUint32(intbuf, f.Value.(uint32))
			_, err := buf.Write(intbuf)
			if err != nil {
				return nil, err
			}
		default:
			return nil, UnsupportedType
		}
	}
	return buf.Bytes(), nil
}

func (r *Record) Value(fieldname string) (interface{}, error) {
	for _, f := range *r {
		if f.Descriptor.Name != fieldname {
			continue
		}
		return f.Value, nil
	}
	return nil, fmt.Errorf("invalid field name")
}
