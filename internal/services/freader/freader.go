package freader

import (
	"errors"
	"fmt"
	"io"

	"github.com/Khatunzev-Anton/filesorter/internal/models"
)

var BytesReadNumberMismatch error = errors.New("the number of bytes read from the source does not match the size of the record")

type FReader interface {
	ReadSerializableRecord(f io.ReadSeeker, offset int64) (models.SerializableRecord, error)
	ReadComparableByField(f io.ReadSeeker, offset int64) (models.ComparableByField, error)
}

type freader struct {
	rd     models.RecordDescriptor
	rdsize int
}

func NewFReader(rd models.RecordDescriptor) (FReader, error) {
	return &freader{
			rd:     rd,
			rdsize: rd.Size(),
		},
		nil
}

func (r *freader) recordbytes(f io.ReadSeeker, offset int64) ([]byte, error) {
	_, err := f.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek the source: %w", err)
	}
	b := make([]byte, r.rdsize)
	n, err := f.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read record: %w", err)
	}
	if n != r.rdsize {
		return nil, BytesReadNumberMismatch
	}
	return b, nil
}

func (r *freader) ReadSerializableRecord(f io.ReadSeeker, offset int64) (models.SerializableRecord, error) {
	b, err := r.recordbytes(f, offset)
	if err != nil {
		return nil, err
	}
	rec, err := models.SerializableRecordFromBytes(b, r.rd)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize record: %w", err)
	}
	return rec, nil
}

func (r *freader) ReadComparableByField(f io.ReadSeeker, offset int64) (models.ComparableByField, error) {
	b, err := r.recordbytes(f, offset)
	if err != nil {
		return nil, err
	}
	rec, err := models.ComparableByFieldFromBytes(b, r.rd)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize record: %w", err)
	}
	return rec, nil
}
