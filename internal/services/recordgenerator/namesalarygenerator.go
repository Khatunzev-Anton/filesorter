package recordgenerator

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Khatunzev-Anton/filesorter/internal/models"
)

var BufferSizeDoesNotMatch error = errors.New("record size differs from the number of bytes written")

//go:generate mockery --name=NamesRepository
type NamesRepository interface {
	Names(page int, pagesize int) ([]string, error)
}

type namesalarygenerator struct {
	rd          models.RecordDescriptor
	namesbuffer []string
	rdsize      int
	rsource     rand.Source
	minsalary   uint32
	maxsalary   uint32
}

func NewNameSalaryGenerator(rd models.RecordDescriptor, namesrepository NamesRepository, minsalary uint32, maxsalary uint32) (RecordGenerator, error) {
	result := &namesalarygenerator{
		rd:          rd,
		namesbuffer: []string{},
		rdsize:      rd.Size(),
		rsource:     rand.NewSource(time.Now().UnixNano()),
		minsalary:   minsalary,
		maxsalary:   maxsalary,
	}

	var err error
	result.namesbuffer, err = namesrepository.Names(0, 5000)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare names buffer: %w", err)
	}
	if len(result.namesbuffer) == 0 {
		return nil, errors.New("names buffer is empty")
	}

	return result, nil
}

func (r *namesalarygenerator) GenerateRecord() (models.SerializableRecord, error) {
	result := &models.Record{}

	for _, fd := range r.rd {
		rand.Seed(time.Now().UnixNano())
		f := models.Field{
			Descriptor: fd,
		}
		switch fd.Name {
		case "name":
			name := r.namesbuffer[rand.New(r.rsource).Intn(len(r.namesbuffer))]
			if len(name) < 3 {
				return nil, errors.New("namesalarygenerator: name is too short")
			}
			f.Value = name
		case "salary":
			salary := r.minsalary
			if r.maxsalary-r.minsalary > 0 {
				salary += (rand.New(r.rsource).Uint32() % (r.maxsalary - r.minsalary + 1))
			}
			f.Value = salary
		default:
			return nil, fmt.Errorf("namesalarygenerator: unknown field %s", fd.Name)
		}

		*result = append(*result, f)
	}

	return result, nil
}
