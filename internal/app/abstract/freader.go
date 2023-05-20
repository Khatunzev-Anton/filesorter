package abstract

import (
	"io"

	"github.com/Khatunzev-Anton/filesorter/internal/models"
)

type FReader interface {
	ReadSerializableRecord(f io.ReadSeeker, offset int64) (models.SerializableRecord, error)
	ReadComparableByField(f io.ReadSeeker, offset int64) (models.ComparableByField, error)
}
