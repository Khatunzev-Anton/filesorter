package models

import (
	"reflect"
	"testing"
)

func TestRecord_Value(t *testing.T) {
	type args struct {
		fieldname string
	}

	rec := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "John Smith",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(100000),
		},
	}

	tests := []struct {
		name    string
		r       *Record
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "name val",
			r:    rec,
			args: args{
				fieldname: "name",
			},
			want:    "John Smith",
			wantErr: false,
		},
		{
			name: "salary val",
			r:    rec,
			args: args{
				fieldname: "salary",
			},
			want:    uint32(100000),
			wantErr: false,
		},
		{
			name: "invalid field",
			r:    rec,
			args: args{
				fieldname: "invalidfield",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Value(tt.args.fieldname)
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_String(t *testing.T) {

	rec := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "John Smith",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(178464),
		},
	}

	tests := []struct {
		name    string
		r       *Record
		want    string
		wantErr bool
	}{
		{
			name:    "John smith",
			r:       rec,
			want:    "John Smith                                                       178464",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Record.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_Bytes(t *testing.T) {
	rec1 := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "John Smith",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(178464),
		},
	}

	rec2 := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "Adela Jäckle",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(142962),
		},
	}

	tests := []struct {
		name    string
		r       *Record
		want    []byte
		wantErr bool
	}{
		{
			name:    "bytes",
			r:       rec1,
			want:    []byte{74, 111, 104, 110, 32, 83, 109, 105, 116, 104, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 185, 2, 0},
			wantErr: false,
		},
		{
			name:    "bytes",
			r:       rec2,
			want:    []byte{65, 100, 101, 108, 97, 32, 74, 195, 164, 99, 107, 108, 101, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 114, 46, 2, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_LessThan(t *testing.T) {
	recjss1 := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "John Smith",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(178464),
		},
	}

	recjss0 := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "John Smith",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(100000),
		},
	}

	recaj := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "Adela Jäckle",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(142962),
		},
	}

	recajd := &Record{
		Field{
			Descriptor: FieldDescriptor{Name: "name", Size: 64, Type: "string"},
			Value:      "Agata Joanna Dachroth",
		},
		Field{
			Descriptor: FieldDescriptor{Name: "salary", Size: 4, Type: "uint32"},
			Value:      uint32(142962),
		},
	}

	type args struct {
		other ComparableByField
		field string
	}
	tests := []struct {
		name    string
		r       *Record
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "JS0 vs JS1 salary",
			r:    recjss0,
			args: args{
				other: recjss1,
				field: "salary",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "JS0 vs JS1 name",
			r:    recjss0,
			args: args{
				other: recjss1,
				field: "name",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "AJ vs AJD name",
			r:    recaj,
			args: args{
				other: recajd,
				field: "name",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "AJD vs AJ salary",
			r:    recajd,
			args: args{
				other: recaj,
				field: "salary",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.LessThan(tt.args.other, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.LessThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Record.LessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}
